import express from 'express';
import fs from 'fs';
import http from 'http';
import { logger } from 'logger';
import { Server as SocketIOServer } from 'socket.io';

import { DEFAULT_ROOM, SOCKET_EVENTS } from '@shared/interfaces/socket';

import { state } from './state.class';
import { config } from './utils/config';

logger.setAppName('dashboard-server');

const expressApp = express();
const httpServer = http.createServer(expressApp);

logger.info(config);

expressApp.use('/', (req, res, next) => {
  logger.info(`${req.method} ${req.url}`);
  next();
});

expressApp.use(express.json());
expressApp.use((req, res, next) => {
  res.header('Access-Control-Allow-Origin', '*');
  next();
});
expressApp.use(express.urlencoded({ extended: true }));

expressApp.use((err: Error, req: express.Request, res: express.Response, next: express.NextFunction) => {
  logger.error(err.stack);

  res.status(500).send(err.message);
});

expressApp.use(express.static(config.iconsPath));
expressApp.use(express.static(config.appDistPath));

httpServer.listen(config.serverPort, () => {
  logger.info(`Server listening on port ${config.serverPort}`);
});

const socketIOServer = new SocketIOServer(httpServer, { cors: { origin: '*' } });

const getNumberOfClients = () => {
  return socketIOServer.sockets.adapter.rooms.get(DEFAULT_ROOM)?.size ?? 0;
};

socketIOServer.on('connection', async (socket) => {
  logger.info(`Socket connected: ${socket.id}`);
  socket.join(DEFAULT_ROOM);

  await state.reloadSetup();
  state.startTracking();

  socket.on(SOCKET_EVENTS.reloadAppSetup, state.reloadSetup);

  socket.on('disconnect', () => {
    logger.info(`Socket disconnected: ${socket.id}`);

    if (getNumberOfClients() === 0) {
      logger.info('No clients connected, stopping tracking');
      state.stopTracking();
    }
  });
});
