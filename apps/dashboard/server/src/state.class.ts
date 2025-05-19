import { logger } from '@libs/logger';
import { DashboardConfig } from '@shared/schemas/dashboard-config.schema';
import { HostState, ServiceState } from '@shared/schemas/dashboard-state.schema';
import { Socket } from 'socket.io';

export class DashboardState {
  private servicesStates: Map<string, ServiceState> = new Map();
  private hostsStates: Map<string, HostState> = new Map();

  constructor(private dashboardConfig: DashboardConfig, private socketsMap: Map<string, Socket>) {
    logger.info('Starting tracking');
  }

  public dispose() {
    logger.info('Disposing of dashboard state');
  }
}
