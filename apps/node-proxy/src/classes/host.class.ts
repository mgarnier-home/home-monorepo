import path from 'path';
import { setTimeout } from 'timers/promises';
import { Worker } from 'worker_threads';

import { logger } from '@libs/logger';

import {
    HostConfig, ManagerThreadMessage, ServiceConfig, ThreadMessageType, WorkerThreadMessage
} from '../utils/interfaces';
import { sendStoppingServer } from '../utils/ntfy.utils';
import { ServerControl } from './serverControl.class';

// import { TCPServiceProxy } from "./tcpServiceProxy.class.ts.old";

interface WorkerConfig {
  worker: Worker;
  service: ServiceConfig;
}

const getServiceId = (service: ServiceConfig) =>
  `${service.name}-${service.proxyPort}:${service.servicePort || service.proxyPort}`;

export class Host {
  private lastPacket: number = Date.now();

  private _hostStarting: boolean = false;
  private _hostStarted: boolean = false;
  private _hostStopping: boolean = false;

  private hostRefreshStatusInterval?: NodeJS.Timeout;
  private hostRefreshServicesInterval?: NodeJS.Timeout;

  private _configUpdated: () => void;

  public config: HostConfig;
  private servicesConfig: ServiceConfig[] = [];
  private workers: Map<string, WorkerConfig> = new Map();

  public updateOptionValue(key: keyof HostConfig['options'], value: any) {
    if (!this.config.options) {
      (this.config.options as any) = {};
    }

    (this.config.options as any)[key] = value;

    this._configUpdated();
  }

  constructor(config: HostConfig, configUpdated: () => void) {
    this.refreshHostStatus = this.refreshHostStatus.bind(this);
    this.spawnServiceWorker = this.spawnServiceWorker.bind(this);
    this.handleWorkerMessage = this.handleWorkerMessage.bind(this);
    this.refreshServices = this.refreshServices.bind(this);

    this.config = config;
    this._configUpdated = configUpdated;

    this.log('Host created');
    this.log(`IP: ${this.config.ip}`);
    this.log(`MAC: ${this.config.macAddress}`);
    this.log(`SSH Username: ${this.config.sshUsername}`);
    this.log(`SSH Password: ${this.config.sshPassword}`);

    this.setupHost();
  }

  private log(...args: any[]) {
    logger.info(`[${this.config.name}]`, ...args);
  }

  private async setupHost() {
    await this.startHost();

    this.hostRefreshStatusInterval = setInterval(this.refreshHostStatus, 1000);
    this.hostRefreshServicesInterval = setInterval(this.refreshServices, 30 * 1000);

    this.refreshServices();
  }

  private async refreshServices() {
    try {
      this.servicesConfig = [
        ...(this.config.enableDocker === true
          ? await ServerControl.getServicesFromDocker(this.config.ip, this.config.dockerPort)
          : []),
        ...(this.config.additionalServices || []),
      ];

      logger.debug('Services refreshed: ', this.servicesConfig);

      for (const service of this.servicesConfig) {
        if (!this.workers.has(getServiceId(service))) {
          this.spawnServiceWorker(service);
        }
      }

      for (const [id, { worker, service }] of this.workers) {
        if (!this.servicesConfig.find((s) => getServiceId(s) === id)) {
          this.disposeWorker(id, worker);
        }
      }
    } catch (error) {
      this.log('Error refreshing services ', error);
      return;
    }
  }

  public async dispose() {
    await Promise.all(
      Array.from(this.workers.values()).map(async ({ worker, service }) => {
        this.sendMessageToWorker(worker, {
          type: ThreadMessageType.DISPOSE_WORKER,
        });

        await setTimeout(1500);

        worker.removeAllListeners();

        worker.terminate();

        this.log(`Worker ${worker.threadId} for ${service.name} disposed`);
      })
    );

    this.workers.clear();

    clearInterval(this.hostRefreshServicesInterval);
    clearInterval(this.hostRefreshStatusInterval);
  }

  private spawnServiceWorker(service: ServiceConfig) {
    const workerPath = path.resolve(__dirname, './worker/proxyWorker');

    const worker = new Worker(workerPath);

    worker.on('online', async () => {
      this.log(`Worker ${worker.threadId} for ${service.name} started`);

      await setTimeout(2000);

      this.sendMessageToWorker(worker, {
        type: ThreadMessageType.START_PROXY,
        data: {
          protocol: service.protocol,
          hostIp: this.config.ip,
          hostName: this.config.name,
          serviceName: service.name,
          proxyPort: service.proxyPort,
          servicePort: service.servicePort || service.proxyPort,
        },
      });
    });

    worker.on('message', this.handleWorkerMessage);

    worker.on('error', (error) => {
      this.log(`Worker ${worker.threadId} for ${service.name} errored with error: ${error}`);

      this.workers.delete(getServiceId(service));

      worker.terminate();
    });

    worker.on('exit', (code) => {
      this.log(`Worker ${worker.threadId} for ${service.name} exited with code ${code}`);

      setImmediate(() => this.spawnServiceWorker(service));
    });

    this.workers.set(getServiceId(service), { worker, service });
  }

  private async disposeWorker(workerId: string, worker: Worker) {
    this.sendMessageToWorker(worker, { type: ThreadMessageType.DISPOSE_WORKER });

    this.log(`Disposing worker ${worker.threadId} for ${workerId}`);

    await setTimeout(1500);

    worker.removeAllListeners();

    worker.terminate();

    this.workers.delete(workerId);
  }

  private sendMessageToWorker(worker: Worker, message: ManagerThreadMessage) {
    worker.postMessage(message);
  }

  private sendMessageToAllWorkers(message: ManagerThreadMessage) {
    this.workers.forEach(({ worker }) => this.sendMessageToWorker(worker, message));
  }

  private handleWorkerMessage(message: WorkerThreadMessage) {
    switch (message.type) {
      case ThreadMessageType.START_HOST:
        return this.startHost();
      case ThreadMessageType.PACKET_RECEIVED:
        return this.handlePacketReceived();
      default:
        break;
    }
  }

  private handlePacketReceived() {
    this.lastPacket = Date.now();
  }

  // public async dispose() {
  //   await Promise.all(
  //     Array.from(this.workers.values()).map(async ({ worker, service }) => {
  //       this.sendMessageToWorker(worker, {
  //         type: ThreadMessageType.DISPOSE_WORKER,
  //       });

  //       await setTimeout(1500);

  //       worker.removeAllListeners();

  //       worker.terminate();

  //       this.log(`Worker ${worker.threadId} for ${service.name} disposed`);
  //     })
  //   );

  //   this.workers.clear();

  //   clearInterval(this.hostRefreshInterval);
  // }

  private shouldStopHost() {
    if (this.config.options?.autoStop === true) {
      const maxAliveTime = this.config.options?.maxAliveTime || 15 * 60 * 1000;

      const timeRemaining = maxAliveTime - (Date.now() - this.lastPacket);
      const seconds = timeRemaining >= 0 ? Math.floor(timeRemaining / 1000) : Math.ceil(timeRemaining / 1000);
      const minutes = timeRemaining >= 0 ? Math.floor(seconds / 60) : Math.ceil(seconds / 60);

      if (Math.abs(seconds) % 15 === 0) {
        if (timeRemaining < 0) {
          this.log(`Server stopped since ${Math.abs(minutes)}m ${Math.abs(seconds % 60)}s`);
        } else {
          this.log(`Time remaining before stopping: ${minutes}m ${seconds % 60}s`);
        }
      }

      return timeRemaining <= 0;
    }
  }

  private async refreshHostStatus() {
    this._hostStarted = await ServerControl.getServerStatus(this.config.ip);

    this.sendMessageToAllWorkers({ type: ThreadMessageType.HOST_STATUS, data: { isStarted: this._hostStarted } });

    if (this.shouldStopHost()) {
      if (this._hostStarted && !this._hostStopping) {
        this.log('Host idle for too long, shutting down');

        await this.stopHost();
      }
      if (!this._hostStarted) {
        this.sendMessageToAllWorkers({ type: ThreadMessageType.DESTROY_SOCKETS });
      }
    }
  }

  public async startHost(): Promise<boolean> {
    this._hostStarting = true;
    this.log('Starting host');

    try {
      await ServerControl.startServer(this.config.macAddress);
    } catch (e) {
      this.log('Error starting host: ' + e);

      this._hostStarting = false;

      return false;
    }

    while (this._hostStarting) {
      this._hostStarting = !(await ServerControl.getServerStatus(this.config.ip));
      await setTimeout(1000);
    }

    this.log('Host started');

    this._hostStarted = true;
    this.sendMessageToAllWorkers({ type: ThreadMessageType.HOST_STATUS, data: { isStarted: true } });

    return true;
  }

  public async stopHost(): Promise<boolean> {
    this._hostStopping = true;
    this.log('Stopping host');

    try {
      await ServerControl.stopServer(this.config.ip, this.config.sshUsername, this.config.sshPassword);

      sendStoppingServer(this.config.name);
    } catch (e) {
      this.log('Error stopping host: ' + e);

      this._hostStopping = false;

      return false;
    }

    while (this._hostStopping) {
      this._hostStopping = await ServerControl.getServerStatus(this.config.ip);
      await setTimeout(1000);
    }

    this.log('Host stopped');

    this._hostStarted = false;
    this.sendMessageToAllWorkers({ type: ThreadMessageType.HOST_STATUS, data: { isStarted: false } });

    return true;
  }

  public async getHostStatus(): Promise<boolean> {
    return ServerControl.getServerStatus(this.config.ip);
  }
}
