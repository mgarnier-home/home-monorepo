import express from 'express';
import http from 'http';
import { logger } from 'logger';

import { SERVER_ROUTES } from '@shared/routes';

import { RestControllers } from './restControllers';
import { config } from './utils/config';
import { bindSocketIOServer } from './wsControllers';

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

expressApp.get(SERVER_ROUTES.CONF, RestControllers.getConf);
expressApp.post(SERVER_ROUTES.PING_HOST, RestControllers.postPingHost);
expressApp.post(SERVER_ROUTES.MAKE_REQUEST, RestControllers.postMakeRequest);
expressApp.post(SERVER_ROUTES.STATUS_CHECKS, RestControllers.postStatusChecks);

bindSocketIOServer(httpServer);

httpServer.listen(config.serverPort, () => {
  logger.info(`Server listening on port ${config.serverPort}`);
});
