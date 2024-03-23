import http from 'http';
import { logger } from 'logger';
import { Server } from 'socket.io';

import { SERVER_ROUTES } from '@shared/routes';

import { Handlers } from './handlers';

const log = (...args: any[]) => {
  logger.info(`[WS]`, ...args);
};

export const bindSocketIOServer = (httpServer: http.Server) => {
  const socketIoServer = new Server(httpServer, { cors: { origin: '*' } });

  socketIoServer.on('connection', (socket) => {
    log(`Socket connected: ${socket.id}`);

    socket.on('apiRequest', async (params, callback) => {
      const { data, route } = params;

      log(route);

      let response = null;

      switch (route) {
        case SERVER_ROUTES.CONF:
          response = await Handlers.getAppConf();
          break;
        case SERVER_ROUTES.PING_HOST:
          response = await Handlers.pingHost(data.ip);
          break;
        case SERVER_ROUTES.MAKE_REQUEST:
          response = await Handlers.makeRequest(data.url, data.method, data.body);
          break;
        case SERVER_ROUTES.STATUS_CHECKS:
          response = await Handlers.getStatusChecks(data.statusChecks);
          break;
        default:
          break;
      }

      callback(response);
    });

    socket.on('disconnect', () => {
      log(`Socket disconnected: ${socket.id}`);
    });
  });
};
