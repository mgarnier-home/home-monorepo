import colors from 'colors';
import { logger } from 'logger';
import { parentPort, threadId } from 'worker_threads';

import {
  HostStatusThreadMessageData,
  ManagerThreadMessage,
  Protocol,
  StartProxyThreadMessageData,
  ThreadMessageType,
  WorkerThreadMessage,
} from '../utils/interfaces.js';
import { TCPServiceProxy } from './tcpServiceProxy.class.js';

const colorList = ['red', 'green', 'yellow', 'blue', 'magenta', 'cyan', 'gray', 'grey'];

export class ProxyWorker {
  public hostStarted: boolean = false;

  public hostName: string = '';

  private proxy: TCPServiceProxy | null = null;

  constructor() {
    this.handleManagerMessage = this.handleManagerMessage.bind(this);
    this.handleStartProxy = this.handleStartProxy.bind(this);

    parentPort?.on('message', this.handleManagerMessage);
  }

  public log(...args: any[]) {
    logger.info(`[${this.hostName}:${threadId}]`, ...args);
  }

  private handleManagerMessage(message: ManagerThreadMessage) {
    switch (message.type) {
      case ThreadMessageType.START_PROXY:
        return this.handleStartProxy(message.data);
      case ThreadMessageType.HOST_STATUS:
        return this.handleHostStatus(message.data);
      case ThreadMessageType.DESTROY_SOCKETS:
        return this.proxy?.destroyServiceSockets();
      case ThreadMessageType.DISPOSE_WORKER:
        return this.dispose();
      default:
        break;
    }
  }

  private dispose() {
    this.proxy?.dispose();
  }

  private handleStartProxy(data: StartProxyThreadMessageData) {
    this.hostName = data.hostName;

    if (data.protocol === Protocol.TCP) {
      this.proxy = new TCPServiceProxy(this, data.serviceName, data.hostIp, data.proxyPort, data.servicePort);
    }
  }

  private handleHostStatus(data: HostStatusThreadMessageData) {
    this.hostStarted = data.isStarted;
  }

  private sendMessageToManager(message: WorkerThreadMessage) {
    parentPort?.postMessage(message);
  }

  public notifyPacketReceived() {
    this.sendMessageToManager({
      type: ThreadMessageType.PACKET_RECEIVED,
    });
  }

  public async startHost(): Promise<boolean> {
    return new Promise((resolve) => {
      const timeout = setTimeout(() => {
        parentPort?.off('message', messageReceived);

        this.log('Host start timeout');

        resolve(false);
      }, 10000);

      const messageReceived = (message: ManagerThreadMessage) => {
        if (message.type === ThreadMessageType.HOST_STATUS && message.data.isStarted) {
          this.hostStarted = message.data.isStarted;

          parentPort?.off('message', messageReceived);
          clearTimeout(timeout);

          resolve(message.data.isStarted);
        }
      };

      parentPort?.on('message', messageReceived);

      this.sendMessageToManager({
        type: ThreadMessageType.START_HOST,
      });
    });
  }
}
