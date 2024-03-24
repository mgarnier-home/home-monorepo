import { connect, Socket } from 'socket.io-client';
import { environment } from 'src/environments/environment';

import { Injectable } from '@angular/core';

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
