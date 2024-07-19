import { connect, Socket } from 'socket.io-client';

import { Injectable } from '@angular/core';

import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class SocketService {
  private socket!: Socket;

  constructor() {}

  public connect() {
    this.socket = connect(environment.apiUrl);
    console.log('SocketService connected to', environment.apiUrl);
  }
}
