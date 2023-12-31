import path from 'path';
import { setTimeout } from 'timers/promises';
import { Worker } from 'worker_threads';

import { __dirname, updateHostConfig } from '../utils/config.js';
import { Host, ManagerThreadMessage, Service, ThreadMessageType, WorkerThreadMessage } from '../utils/interfaces.js';
import { sendStoppingServer } from '../utils/ntfy.js';
import { ServerControl } from '../utils/serverCtrl.js';

// import { TCPServiceProxy } from "./tcpServiceProxy.class.ts.old";

interface WorkerConfig {
  worker: Worker;
  service: Service;
}

export class HostServer {
  private lastPacket: number = Date.now();

  private _hostStarting: boolean = false;
  private _hostStarted: boolean = false;
  private _hostStopping: boolean = false;

  private hostRefreshInterval?: NodeJS.Timeout;

  private config: Host;
  private servicesConfig: Service[] = [];
  private workers: Map<number, WorkerConfig> = new Map();

  constructor(host: Host) {
    this.refreshHost = this.refreshHost.bind(this);
    this.spawnServiceWorker = this.spawnServiceWorker.bind(this);
    this.handleWorkerMessage = this.handleWorkerMessage.bind(this);

    this.config = host;

    this.log('Host created');
    this.log(`IP: ${this.config.ip}`);
    this.log(`MAC: ${this.config.macAddress}`);
    this.log(`SSH Username: ${this.config.sshUsername}`);
    this.log(`SSH Password: ${this.config.sshPassword}`);

    this.setupHost();
  }

  private log(...args: any[]) {
    console.log(`[${this.config.name}]`, ...args);
  }

  private async setupHost() {
    await this.startHost();

    this.servicesConfig = [
      // ...(await ServerControl.getServices(this.config.ip, this.config.sshUsername, this.config.sshPassword)),
      ...ServerControl.getServicesFromEnv(this.config.name),
      ...(this.config.additionalServices || []),
    ];

    this.servicesConfig.forEach(this.spawnServiceWorker);

    this.hostRefreshInterval = setInterval(this.refreshHost, 1000);

    this.log('Services', this.servicesConfig);
  }

  private spawnServiceWorker(service: Service) {
    const workerPath = path.resolve(__dirname, '../worker/proxyWorker');

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

      this.workers.delete(worker.threadId);

      worker.terminate();
    });

    worker.on('exit', (code) => {
      this.log(`Worker ${worker.threadId} for ${service.name} exited with code ${code}`);

      setImmediate(() => this.spawnServiceWorker(service));
    });

    this.workers.set(worker.threadId, { worker, service });
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

    clearInterval(this.hostRefreshInterval);
  }

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

  private async refreshHost() {
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

  public get hostStarted() {
    return this._hostStarted;
  }

  public updateOptionValue(key: string, value: any) {
    if (!this.config.options) {
      (this.config.options as any) = {};
    }

    (this.config.options as any)[key] = value;
  }

  public disableAutoStop() {
    this.updateOptionValue('autoStop', false);

    updateHostConfig(this.config);
  }

  public async enableAutoStop() {
    this.updateOptionValue('autoStop', true);

    updateHostConfig(this.config);
  }

  public getAutoStop() {
    return this.config.options?.autoStop;
  }
}
