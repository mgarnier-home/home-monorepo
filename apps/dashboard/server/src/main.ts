import express from 'express';
import http from 'http';
import { Server as SocketIOServer } from 'socket.io';

import { setVersionEndpoint } from '@libs/api-version';
import { logger } from '@libs/logger';
// import { DEFAULT_ROOM } from '@shared/interfaces/socket';

// import { state } from './state.class';
import { dashboardConfigSchema } from '@shared/schemas/config.schema';
import { socketEvents } from '@shared/socketEvents.enum';
import { readFileSync } from 'fs';
import YAML from 'yaml';
import { z } from 'zod';
import { config } from './utils/config';

logger.setAppName('dashboard-server');

const expressApp = express();
const httpServer = http.createServer(expressApp);

logger.info(config);

setVersionEndpoint(expressApp);

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

logger.info(`Serving icon files from ${config.iconsPath} and app files from ${config.appDistPath}`);

expressApp.use(express.static(config.iconsPath));
expressApp.use(express.static(config.appDistPath));

httpServer.listen(config.serverPort, () => {
  logger.info(`Server listening on port ${config.serverPort}`);
});

const socketIOServer = new SocketIOServer(httpServer, { cors: { origin: '*' } });

const getNumberOfClients = () => {
  return socketIOServer.sockets.sockets.size;
};

const loadConfig = (): z.infer<typeof dashboardConfigSchema> => {
  const yamlConf = YAML.parse(readFileSync(config.configFilePath, 'utf8'));

  const result = dashboardConfigSchema.parse(yamlConf);

  return result;
};

socketIOServer.on('connection', async (socket) => {
  logger.info(`Socket connected: ${socket.id}`);

  socket.on('disconnect', () => {
    logger.info(`Socket disconnected: ${socket.id}`);

    if (getNumberOfClients() === 0) {
      logger.info('No clients connected, stopping tracking');
      // state.stopTracking();
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
