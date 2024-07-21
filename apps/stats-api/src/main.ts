import express, { Request, Response } from 'express';
import { body } from 'express-validator';

import { setVersionEndpoint } from '@libs/api-version';
import { logger } from '@libs/logger';

import { setupApiRoutes } from './api/setupApi';
import { ApiUtils } from './api/utils';
import { databaseCpu } from './db/cpu';
import { databaseGpu } from './db/gpu';
import { databaseNetwork } from './db/network';
import { databaseRam } from './db/ram';
import { config } from './utils/config';

logger.setAppName('stats-api');
const app = express();

setVersionEndpoint(app);

app.use(express.json());
app.use((req, res, next) => {
  res.header('Access-Control-Allow-Origin', '*');
  next();
});
app.use((err: any, req: Request, res: Response, next: any) => {
  logger.error(err.stack);
  res.status(500).send('Something broke!');
});

app.use('/api/:hostname', (req: Request<ApiUtils.Hostname, any, any>, res, next) => {
  res.locals.hostname = req.params.hostname.toLowerCase();

  next();
});

setupApiRoutes(
  app,
  'cpu',
  [
    body('cores').isArray().withMessage('Cores must be an array'),
    body('cores.*.load').isFloat().withMessage('Load must be a number'),
    body('cores.*.temp').isFloat().optional({ nullable: true }).withMessage('Temp must be a number'),
    body('cores.*.core').isInt().withMessage('Core must be an integer'),
  ],
  databaseCpu
);
setupApiRoutes(
  app,
  'gpu',
  [
    body('gpus').isArray().withMessage('Gpus must be an array'),
    body('gpus.*.index').isInt().withMessage('Index must be an integer'),
    body('gpus.*.model').isString().withMessage('Model must be a string'),
    body('gpus.*.memoryUsed').isInt().withMessage('Memory used must be an integer'),
    body('gpus.*.load').isFloat().withMessage('Load must be a number'),
    body('gpus.*.temp').isFloat().optional({ nullable: true }).withMessage('Temp must be a number'),
    body('gpus.*.powerDraw').isFloat().optional({ nullable: true }).withMessage('Power draw must be a number'),
  ],
  databaseGpu
);
setupApiRoutes(
  app,
  'ram',
  [
    body('load').isInt().withMessage('Load must be an integer'), //
  ],
  databaseRam
);
setupApiRoutes(
  app,
  'network',
  [
    body('down').isFloat().withMessage('Down must be a number'),
    body('up').isFloat().withMessage('Up must be a number'),
  ],
  databaseNetwork
);

app.get('/', (req, res) => {
  res.send('OK');
});

app.listen(config.serverPort, () => {
  logger.info(`Server listening on port ${config.serverPort}`);
});
