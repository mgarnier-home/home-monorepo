import { logger } from '@libs/logger';
import { Utils } from '@libs/utils';
import { actionTypeEnum, DashboardConfig, HealthCheck, Host, Service } from '@shared/schemas/dashboard-config.schema';
import { HealthCheckState, HostState, ServiceState } from '@shared/schemas/dashboard-state.schema';
import { Socket } from 'socket.io';
import { pingHost } from './utils/utils';

export class DashboardState {
  private servicesStates: Map<string, ServiceState> = new Map();
  private hostsStates: Map<string, HostState> = new Map();

  constructor(private dashboardConfig: DashboardConfig, private socketsMap: Map<string, Socket>) {
    logger.info('Starting tracking');
  }

  public dispose() {
    logger.info('Disposing of dashboard state');
  }

  private async getHostState(hostId: string, host: Host): Promise<HostState> {
    const { nodesight, ip } = host;

    const ping = await pingHost(ip);

    return {
      id: hostId,
      ping: ping.duration,
      status: ping.ping ? 'ok' : 'error',
      // dockerStatus: 'ok',
    };
  }

  private async getServiceState(serviceId: string, service: Service): Promise<ServiceState> {
    const { dockerName, healthCheck, healthChecks } = service;

    const healthCheckStates = await Promise.all(
      (healthChecks || []).map((healthCheck) => this.getHealthCheck(serviceId, healthCheck))
    );
    const healthCheckState = healthCheck ? await this.getHealthCheck(serviceId, healthCheck) : undefined;

    return {
      id: serviceId,
      healthCheck: healthCheckState,
      dockerStatus: 'ok',
      healthChecks: healthCheckStates,
    };
  }

  private async getHealthCheck(healthCheckId: string, healthCheck: HealthCheck): Promise<HealthCheckState> {
    if (healthCheck.action.type !== actionTypeEnum.enum.request) {
      throw new Error('Invalid action type for health check');
    }
    const startTime = Date.now();

    const data = await Utils.fetchWithTimeout(healthCheck.url, 5000, {
      method: healthCheck.action.method,
      headers: {
        Status: 'true',
      },
    });

    const text = await data.text();

    const endTime = Date.now();

    return {
      id: healthCheckId,
      name: healthCheck.name,
      code: data.status,
      responseTime: endTime - startTime,
      response: text,
    };
  }
}
