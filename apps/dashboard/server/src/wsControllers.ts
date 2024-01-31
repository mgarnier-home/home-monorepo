import http from 'http';
import { logger } from 'logger';
import { Server } from 'socket.io';

import { ServerRoutes } from '@shared/routes';

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
        case ServerRoutes.CONF:
          response = await Handlers.getAppConf();
          break;
        case ServerRoutes.PING_HOST:
          response = await Handlers.pingHost(data.ip);
          break;
        case ServerRoutes.MAKE_REQUEST:
          response = await Handlers.makeRequest(data.url, data.method, data.body);
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
