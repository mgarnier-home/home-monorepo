import { createSocket, Socket } from 'dgram';
import { EventEmitter } from 'stream';

import { config } from './config';

export class Syslog extends EventEmitter {
  private socket: Socket | null = null;

  constructor() {
    super();
  }

  public async start(options: { port: number }) {
    if (this.socket) {
      return this.socket;
    }

    return new Promise((resolve, reject) => {
      this.socket = createSocket('udp4');

      this.socket.on('listening', () => {
        this.emit('start');

        resolve(this.socket);
      });

      this.socket.on('error', (err) => {
        this.emit('error', err);
      });

      this.socket.on('message', (msg, rinfo) => {
        this.emit('message', {
          date: new Date(),
          host: config.hostsMap[rinfo.address] || rinfo.address,
          message: msg.toString('utf8'),
          protocol: rinfo.family,
        });
      });

      this.socket.on('close', () => {
        this.emit('close');
      });

      this.socket.bind(options.port);
    });
  }

  public async stop() {
    if (!this.socket) {
      return;
    }

    return new Promise<void>((resolve, reject) => {
      this.socket?.close(() => {
        this.socket = null;

        resolve();
      });
    });
  }
}
