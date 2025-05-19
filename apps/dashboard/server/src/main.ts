import express from 'express';
import http from 'http';
import { setVersionEndpoint } from '@libs/api-version';
import { logger } from '@libs/logger';
import { config } from './utils/config';
import { startSocketIOServer } from './socket';

logger.setAppName('dashboard-server');

const expressApp = express();
const httpServer = http.createServer(expressApp);

logger.info('Config from env :', config);

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

startSocketIOServer(httpServer);
