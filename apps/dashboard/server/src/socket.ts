import http from 'http';
import YAML from 'yaml';

import { DashboardConfig, dashboardConfigSchema } from '@shared/schemas/dashboard-config.schema';
import { readFileSync } from 'fs';

import { Server as SocketIOServer } from 'socket.io';
import { config } from './utils/config';
import { logger } from '@libs/logger';
import { socketEvents } from '@shared/socketEvents.enum';
import { DashboardState } from './state.class';

let dashboardState: DashboardState | null = null;

export function startSocketIOServer(httpServer: http.Server) {
  const socketIOServer = new SocketIOServer(httpServer, { cors: { origin: '*' } });

  const getNumberOfClients = () => {
    return socketIOServer.sockets.sockets.size;
  };

  const loadConfig = (): DashboardConfig => {
    const yamlConf = YAML.parse(readFileSync(config.configFilePath, 'utf8'));

    const result = dashboardConfigSchema.parse(yamlConf);

    return result;
  };

  socketIOServer.on('connection', async (socket) => {
    logger.info(`Socket connected: ${socket.id}`);

    if (!dashboardState) {
      dashboardState = new DashboardState(loadConfig(), socketIOServer.sockets.sockets);
    }

    socket.on('disconnect', () => {
      logger.info(`Socket disconnected: ${socket.id}`);

      if (getNumberOfClients() === 0 && dashboardState) {
        logger.info('No clients connected, stopping tracking');
        dashboardState.dispose();
        dashboardState = null;
      }
    });

    try {
      const config = loadConfig();

      socket.emit(socketEvents.Enum.dashboardConfig, config);
    } catch (error) {
      socket.emit('error', 'Error loading config file');
      logger.error('Error loading config file', error);

      socket.disconnect();
    }
  });

  logger.info('Socket.IO server started on port', config.serverPort);
}
