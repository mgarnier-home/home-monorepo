import http from 'http';
import { Server } from 'socket.io';

import { ServerRoutes } from '@shared/routes';

import { Handlers } from './handlers';

const log = (...args: any[]) => {
  console.log(`[WS]`, ...args);
};

export const bindSocketIOServer = (httpServer: http.Server) => {
  const socketIoServer = new Server(httpServer, { cors: { origin: '*' } });

  socketIoServer.on('connection', (socket) => {
    console.log(`Socket connected: ${socket.id}`);

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
      console.log(`Socket disconnected: ${socket.id}`);
    });
  });
};
