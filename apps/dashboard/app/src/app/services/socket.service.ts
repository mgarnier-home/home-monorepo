import { connect, Socket } from 'socket.io-client';

import { Injectable, signal } from '@angular/core';

import { environment } from '../../environments/environment';
import { socketEvents } from '@shared/socketEvents.enum';
import { dashboardConfigSchema } from '@shared/schemas/config.schema';
import { z } from 'zod';

@Injectable({
  providedIn: 'root',
})
export class SocketService {
  private socket!: Socket;

  public dashboardConfig = signal<z.infer<typeof dashboardConfigSchema> | null>(null);

  constructor() {
    this._onDashboardConfig = this._onDashboardConfig.bind(this);
  }

  public connect() {
    this.socket = connect(environment.apiUrl);
    console.log('SocketService connected to', environment.apiUrl);

    this.socket.on(socketEvents.Enum.dashboardConfig, this._onDashboardConfig);
  }

  private _onDashboardConfig(config: z.infer<typeof dashboardConfigSchema>) {
    console.log('SocketService received config', config);
    this.dashboardConfig.set(config);
  }
}
