import express from 'express';
import fs from 'fs';
import http from 'http';
import { logger } from 'logger';
import { Server as SocketIOServer } from 'socket.io';

import { DEFAULT_ROOM, SOCKET_EVENTS } from '@shared/interfaces/socket';

import { setup } from './setup.class';
import { config } from './utils/config';

import type { Setup } from '@shared/interfaces/setup';
logger.setAppName('dashboard-server');

const log = (...args: any[]) => {
  logger.info(`[API]`, ...args);
};

const expressApp = express();
const httpServer = http.createServer(expressApp);

log(config);

expressApp.use('/', (req, res, next) => {
  log(`${req.method} ${req.url}`);
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

let appStatus;
let configInterval: NodeJS.Timeout | null = null;

const loadAppSetup = async () => {
  if (fs.existsSync(config.appSetupPath)) {
    const appSetupContentStr = await fs.promises.readFile(config.appSetupPath, 'utf-8');

    setup.reloadAppSetup(appSetupContentStr);

    log('App setup reloaded');

    logger.debug(setup);
  } else {
    logger.error('Config file not found');
  }
};

const getNumberOfClients = () => {
  return socketIOServer.sockets.adapter.rooms.get(DEFAULT_ROOM)?.size ?? 0;
};

const setupInterval = () => {
  if (configInterval) {
    return;
  }

  configInterval = setInterval(() => {}, setup.globalConfig.statusCheckInterval);
};

socketIOServer.on('connection', async (socket) => {
  log(`Socket connected: ${socket.id}`);
  socket.join(DEFAULT_ROOM);

  loadAppSetup();

  setupInterval();

  socket.on(SOCKET_EVENTS.reloadAppSetup, loadAppSetup);

  socket.on('disconnect', () => {
    log(`Socket disconnected: ${socket.id}`);

    if (getNumberOfClients() === 0 && configInterval) {
      clearInterval(configInterval);

      configInterval = null;
    }
  });
});
