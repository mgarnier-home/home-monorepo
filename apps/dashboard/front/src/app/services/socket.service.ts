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
  private socket!: Socket;

  private stateService = inject(StateService);

  public dashboardConfig = signal<DashboardConfig | null>(null);

  constructor() {
    this._onDashboardConfig = this._onDashboardConfig.bind(this);
    this._onHostStateUpdate = this._onHostStateUpdate.bind(this);
    this._onServiceStateUpdate = this._onServiceStateUpdate.bind(this);
  }

  public connect() {
    this.socket = connect(environment.apiUrl, {
      // transports: ['websocket'],
    });
    console.log('SocketService connected to', environment.apiUrl);

    this.socket.on(socketEvents.Enum.dashboardConfig, this._onDashboardConfig);
    this.socket.on(socketEvents.Enum.hostStateUpdate, this._onHostStateUpdate);
    this.socket.on(socketEvents.Enum.serviceStateUpdate, this._onServiceStateUpdate);
  }

  private _onDashboardConfig(config: DashboardConfig) {
    console.log('DashboardConfig received', config);
    this.dashboardConfig.set(config);
  }

  private _onHostStateUpdate(hostId: string, hostState: HostState) {
    console.log('SocketService received host state', hostId, hostState);

    this.stateService.updateHostState(hostId, hostState);
  }

  private _onServiceStateUpdate(serviceId: string, serviceState: ServiceState) {
    console.log('SocketService received service state', serviceId, serviceState);

    this.stateService.updateServiceState(serviceId, serviceState);
  }
}
