import { io } from 'socket.io-client';

import { ServerRoutes } from '../../../shared/routes';

class Socket {
  private _socket = io('', { autoConnect: false });
  private _connected = false;

  constructor() {
    this.onConnect = this.onConnect.bind(this);
    this.onDisconnect = this.onDisconnect.bind(this);

    this.socket.on('connect', this.onConnect);
    this.socket.on('disconnect', this.onDisconnect);
  }

  public get socket() {
    return this._socket;
  }

  public get connected() {
    return this._connected;
  }

  private onConnect() {
    console.log(`Socket connected: ${this.socket.id}`);

    this._connected = true;
  }

  private onDisconnect() {
    console.log(`Socket disconnected: ${this.socket.id}`);

    this._connected = false;
  }

  public apiRequest<T, U>(route: ServerRoutes, data: T): Promise<U> {
    return new Promise((resolve, reject) => {
      if (this.connected) {
        setTimeout(() => {
          reject('Request timeout');
        }, 10000);

        this.socket.emit('apiRequest', { route, data }, (response: U) => {
          resolve(response);
        });
      } else {
        reject('Socket not connected');
      }
    });
  }

  public connect() {
    this.socket.connect();
  }

  public dispose() {
    this.socket.disconnect();

    this.socket.off('connect', this.onConnect);
    this.socket.off('disconnect', this.onDisconnect);
  }
}

export const socket = new Socket();
