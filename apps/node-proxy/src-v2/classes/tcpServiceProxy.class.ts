import { Server, Socket } from 'net';

import { sendStartingServer } from '../utils/ntfy.js';
import { ProxyWorker } from './proxyWorker.class.js';

interface SocketContext {
  key: string;
  proxySocket: Socket;
  buffers: Buffer[];
  serviceConnected: boolean;
  serviceSocket?: Socket;
}

function uniqueKey(socket: Socket) {
  var key = socket.remoteAddress + ':' + socket.remotePort;
  return key;
}

export class TCPServiceProxy {
  private proxyWorker: ProxyWorker;

  private proxyPort: number;
  private hostIp: string;
  private servicePort: number;
  private serviceName: string;

  private proxySockets: Map<string, SocketContext> = new Map();

  private proxyServer: Server;

  constructor(worker: ProxyWorker, serviceName: string, hostIp: string, proxyPort: number, servicePort: number) {
    this.handleNewClientConnection = this.handleNewClientConnection.bind(this);

    this.proxyWorker = worker;
    this.hostIp = hostIp;
    this.proxyPort = proxyPort;
    this.servicePort = servicePort;
    this.serviceName = serviceName;

    this.proxyServer = new Server(this.handleNewClientConnection);

    this.proxyServer.listen(this.proxyPort);

    this.log(`Proxy server listening on tcp port`);
  }

  private log(...args: any[]) {
    this.proxyWorker.log(`[${this.serviceName}:${this.proxyPort}]`, ...args);
  }

  public dispose() {
    this.proxySockets.forEach((context) => {
      context.proxySocket.destroy();
    });
    this.proxySockets.clear();

    this.proxyServer.close();
    this.proxyServer.unref();
  }

  public destroyServiceSockets() {
    this.proxySockets.forEach((context) => {
      if (context.serviceSocket) context.serviceSocket.destroy();
      context.serviceConnected = false;
    });
  }

  private gracefullyEndSocket(statusCode: number, socket: Socket) {
    const httpResponse =
      `HTTP/1.1 ${statusCode} OK\r\n` +
      'Content-Type: text/plain\r\n' +
      'Connection: close\r\n' +
      '\r\n' + // Important: Blank line to separate headers from body
      'Success';

    socket.write(httpResponse, () => {
      // The 'end' function ensures that all the data is flushed to the underlying system
      // before the socket is fully closed
      socket.end();
    });
  }

  private handleNewClientConnection(proxySocket: Socket) {
    const key = uniqueKey(proxySocket);

    this.log(`Client connection established: ${key}`);

    const context: SocketContext = {
      key,
      proxySocket,
      buffers: [],
      serviceConnected: false,
    };

    this.proxySockets.set(key, context);

    proxySocket.on('data', async (data) => {
      const isStatusRequest = await this.handleIncomingDataFromClient(context, data);

      if (isStatusRequest) {
        this.gracefullyEndSocket(204, proxySocket);
      }
    });

    proxySocket.on('close', () => {
      this.proxySockets.delete(key);

      this.log(`Client connection closed: ${key}`);

      if (context.serviceSocket) context.serviceSocket.destroy();
    });

    proxySocket.on('error', (err) => {
      this.log(`Client connection error: ${key} : ${err.message}`);

      if (context.serviceSocket) context.serviceSocket.destroy();
    });
  }

  private async handleIncomingDataFromClient(context: SocketContext, data: Buffer): Promise<boolean | undefined> {
    // this.log("Incoming data from client", data.toString());

    if (context.serviceConnected && context.serviceSocket) {
      this.proxyWorker.notifyPacketReceived();

      context.serviceSocket.write(data);
    } else {
      const [headersSection] = data.toString().split('\r\n\r\n', 1) as [string];
      const headersArray = headersSection.split('\r\n').slice(1);
      const headers: { [key: string]: string } = {};

      headersArray.forEach((headerLine) => {
        const [key, value] = headerLine.split(': ', 2) as [string, string];
        headers[key.toLowerCase()] = value.toLowerCase();
      });

      const incomingIp = headers['x-real-ip'] || headers['x-forwarded-for'];

      if (incomingIp) {
        this.log(`Request coming from ${incomingIp}`);
      }

      if (headers['status'] === 'true' && !this.proxyWorker.hostStarted) {
        this.log('Status packet and host is not started => ignoring');
        return true;
      }

      this.proxyWorker.notifyPacketReceived();

      context.buffers.push(data);

      if (!this.proxyWorker.hostStarted) {
        if (!(await this.proxyWorker.startHost())) {
          return;
        } else {
          sendStartingServer(this.proxyWorker.hostName, incomingIp ?? '', this.serviceName);
        }
      }

      // this.log("Creating service socket");

      this.createServiceSocket(context);
    }
  }

  private createServiceSocket(context: SocketContext) {
    context.serviceSocket = new Socket();
    context.serviceSocket.connect({ port: this.servicePort, host: this.hostIp }, () => {
      this.writeBuffer(context);
    });

    context.serviceSocket.on('data', (data) => this.handleIncomingDataFromService(context, data));

    context.serviceSocket.on('close', () => {
      context.proxySocket.destroy();
      delete context.serviceSocket;

      context.serviceConnected = false;
    });

    context.serviceSocket.on('error', (err) => {
      this.log(`Service connection error: ${context.key} : ${err.message}`);
      context.proxySocket.destroy();
      delete context.serviceSocket;

      context.serviceConnected = false;
    });
  }

  private handleIncomingDataFromService(context: SocketContext, data: Buffer) {
    context.proxySocket.write(data);
  }

  private writeBuffer(context: SocketContext) {
    context.serviceConnected = true;

    if (context.buffers.length > 0) {
      context.buffers.forEach((buffer) => context.serviceSocket!.write(buffer));
      context.buffers = [];
    }
  }
}
