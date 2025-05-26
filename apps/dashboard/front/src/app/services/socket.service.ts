import { connect, Socket } from 'socket.io-client';

import { inject, Injectable, signal } from '@angular/core';

import { environment } from '../../environments/environment';
import { socketEvents } from '@shared/socketEvents.enum';
import { DashboardConfig, dashboardConfigSchema } from '@shared/schemas/dashboard-config.schema';
import { z } from 'zod';
import { HostState, ServiceState } from '@shared/schemas/dashboard-state.schema';
import { StateService } from './state.service';

@Injectable({
  providedIn: 'root',
})
export class SocketService {
  private ws: WebSocket;

  private stateService = inject(StateService);

  public dashboardConfig = signal<DashboardConfig | null>(null);

  constructor() {
    this.ws = new WebSocket(environment.apiUrl.replace('http', 'ws') + '/ws');

    this.ws.onopen = () => {
      console.log('Connected to server');
      this.sendMessage('Hello from Angular!');
    };

    this.ws.onmessage = (event) => {
      console.log(event);
    };

    this._onDashboardConfig = this._onDashboardConfig.bind(this);
    // this._onHostStateUpdate = this._onHostStateUpdate.bind(this);
    // this._onServiceStateUpdate = this._onServiceStateUpdate.bind(this);
  }

  private sendMessage(message: string) {
    if (this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(message);
    }
  }

  // public connect() {
  //   this.socket = connect(environment.apiUrl);
  //   console.log('SocketService connected to', environment.apiUrl);

  //   this.socket.on(socketEvents.Enum.dashboardConfig, this._onDashboardConfig);
  //   this.socket.on(socketEvents.Enum.hostStateUpdate, this._onHostStateUpdate);
  //   this.socket.on(socketEvents.Enum.serviceStateUpdate, this._onServiceStateUpdate);
  // }

  private _onDashboardConfig(config: DashboardConfig) {
    console.log('SocketService received config', config);
    this.dashboardConfig.set(config);
  }

  // private _onHostStateUpdate(hostId: string, hostState: HostState) {
  //   console.log('SocketService received host state', hostId, hostState);

  //   this.stateService.updateHostState(hostId, hostState);
  // }

  // private _onServiceStateUpdate(serviceId: string, serviceState: ServiceState) {
  //   console.log('SocketService received service state', serviceId, serviceState);

  //   this.stateService.updateServiceState(serviceId, serviceState);
  // }
}
