export enum Protocol {
  TCP = 'tcp',
  UDP = 'udp',
}

export interface ServiceConfig {
  name: string;
  proxyPort: number;
  servicePort?: number;
  protocol: Protocol;
}

export interface HostConfig {
  name: string;
  ip: string;
  enableDocker: boolean;
  sshPort: number;
  macAddress: string;
  sshUsername: string;
  sshPassword: string;
  additionalServices?: ServiceConfig[];
  options: {
    maxAliveTime: number;
    autoStop: boolean;
  };
}

export interface AppData {
  hostsConfigs: HostConfig[];
}

export enum ThreadMessageType {
  // Manager messages
  START_PROXY = 'startProxy',
  HOST_STATUS = 'hostStatus',
  DESTROY_SOCKETS = 'destroySockets',
  DISPOSE_WORKER = 'disposeWorker',

  // Worker messages
  START_HOST = 'startHost',
  PACKET_RECEIVED = 'packetReceived',
}

export interface StartProxyThreadMessageData {
  protocol: Protocol;
  hostIp: string;
  hostName: string;
  serviceName: string;
  proxyPort: number;
  servicePort: number;
}

export interface HostStatusThreadMessageData {
  isStarted: boolean;
}

export type ManagerThreadMessage =
  | {
      type: ThreadMessageType.START_PROXY;
      data: StartProxyThreadMessageData;
    }
  | {
      type: ThreadMessageType.HOST_STATUS;
      data: HostStatusThreadMessageData;
    }
  | {
      type: ThreadMessageType.DESTROY_SOCKETS;
    }
  | {
      type: ThreadMessageType.DISPOSE_WORKER;
    };

export type WorkerThreadMessage =
  | {
      type: ThreadMessageType.START_HOST;
    }
  | {
      type: ThreadMessageType.PACKET_RECEIVED;
    };
