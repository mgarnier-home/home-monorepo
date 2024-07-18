import Express, { NextFunction, Request, Response } from 'express';

import { logger } from '@libs/logger';

import { Host } from '../classes/host.class';
import { getHost } from './host.utils';

type ResponseWithHost = Response & { locals: { host: Host } };

const log = (...args: any[]) => {
  logger.info(`[API]`, ...args);
};

export const createExpressApp = (apiPort: number) => {
  const app = Express();

  app.use((err: Error, req: Request, res: Response, next: NextFunction) => {
    logger.error(err.stack);

    res.status(500).send(err.message);
  });

  app.use('/', (req, res, next) => {
    log(`${req.method} ${req.url}`);
    next();
  });

  app.use('/control/:host', (req, res: ResponseWithHost, next) => {
    try {
      const host = getHost(req.params.host);

      if (host) {
        res.locals.host = host;
        next();
      } else {
        throw new Error('Host not found');
      }
    } catch (err) {
      res.send('Error: ' + err);
    }
  });

  app.get('/control/:host/start', async (req: Request, res: ResponseWithHost) => {
    await res.locals.host.startHost();

    res.send('Server started');
  });

  app.get('/control/:host/stop', async (req: Request, res: ResponseWithHost) => {
    await res.locals.host.stopHost();

    res.send('Server stopped');
  });

  app.get('/control/:host/start-stop', async (req: Request, res: ResponseWithHost) => {
    if (await res.locals.host.getHostStatus()) {
      await res.locals.host.stopHost();
      res.send('Server stopped');
    } else {
      await res.locals.host.startHost();
      res.send('Server started');
    }
  });

  app.get('/control/:host/status', async (req, res: ResponseWithHost) => {
    const status = await res.locals.host.getHostStatus();

    res.status(status ? 200 : 500).send(status ? 'Started' : 'Unreachable');
  });

  app.get('/control/:host/autostop-disable', (req, res: ResponseWithHost) => {
    res.locals.host.updateOptionValue('autoStop', false);

    res.send('Auto stop disabled');
  });

  app.get('/control/:host/autostop-enable', (req, res: ResponseWithHost) => {
    res.locals.host.updateOptionValue('autoStop', true);

    res.send('Auto stop enabled');
  });

  app.get('/control/:host/autostop-toggle', (req, res: ResponseWithHost) => {
    res.locals.host.updateOptionValue('autoStop', !res.locals.host.config.options.autoStop);

    res.send(`Auto stop toggled ${res.locals.host.config.options.autoStop ? 'on' : 'off'}`);
  });

  app.get('/control/:host/autostop-status', (req, res: ResponseWithHost) => {
    res
      .status(res.locals.host.config.options.autoStop ? 200 : 569)
      .send(res.locals.host.config.options.autoStop ? 'Enabled' : 'Disabled');
  });

  app.get('/', (req, res) => {
    res.send('OK');
  });

  app.listen(apiPort, () => {
    log(`API listening on port ${apiPort}`);
  });

  log('Express API created');
};
