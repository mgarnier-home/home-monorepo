import { watch as watchFile } from 'chokidar';
import Express, { NextFunction, Request, Response } from 'express';
import { logger } from 'logger';

import { HostServer } from './classes/hostServer.class.js';
import { configFilePath, loadConfig } from './utils/config.js';
import { Config } from './utils/interfaces.js';

logger.setAppName('node-proxy');

let config: Config | undefined = undefined;

const apiPort = process.env.SERVER_PORT || 3000;

let hostServers: HostServer[] = [];

const log = (...args: any[]) => {
  logger.info(`[API]`, ...args);
};

const reloadConfig = async () => {
  config = await loadConfig();
};

const reloadHostsServers = async () => {
  if (hostServers.length > 0) {
    log('Disposing old hosts');

    for (const host of hostServers) {
      await host.dispose();
    }
  }

  if (!config) {
    throw new Error('Config not loaded');
  } else {
    hostServers = [];
    for (const host of config.hosts) {
      hostServers.push(new HostServer(host));
    }
  }
};

const getHost = (host: string): HostServer => {
  if (!config) {
    throw new Error('Config not loaded');
  }

  const hostIndex = Number.isNaN(parseInt(host, 10))
    ? config.hosts.findIndex((h) => h.name === host)
    : parseInt(host, 10);

  log(`Searching for host ${host} at index ${hostIndex}`);

  if (hostIndex === -1 || hostIndex >= hostServers.length) {
    throw new Error('Host not found');
  }

  return hostServers[hostIndex] as HostServer;
};

const startHost = async (req: Request, res: Response) => {
  const host = res.locals.host as HostServer;

  if (await host.startHost()) {
    res.send('Server started');
  } else {
    res.status(500).send('Error starting server');
  }
};

const stopHost = async (req: Request, res: Response) => {
  const host = res.locals.host as HostServer;

  if (await host.stopHost()) {
    res.send('Server stopped');
  } else {
    res.status(500).send('Error stopping server');
  }
};

const main = async () => {
  if (process.env.AUTO_RELOAD_CONFIG === 'true') {
    log('Watching config file for changes');
    watchFile(configFilePath, { ignoreInitial: false }).on('all', async (event, path) => {
      log('config file changed, reloading config');

      await reloadConfig();

      await reloadHostsServers();
    });
  } else {
    log('Not watching config file for changes');
    await reloadConfig();

    await reloadHostsServers();
  }

  const app = Express();

  app.use((err: Error, req: Request, res: Response, next: NextFunction) => {
    logger.error(err.stack);

    res.status(500).send(err.message);
  });

  app.use('/', (req, res, next) => {
    log(`${req.method} ${req.url}`);
    next();
  });

  app.use('/control/:host', (req, res, next) => {
    try {
      res.locals.host = getHost(req.params.host);

      if (res.locals.host) {
        next();
      } else {
        throw new Error('Host not found');
      }
    } catch (err) {
      res.send('Error: ' + err);
    }
  });

  app.get('/control/:host/startstop', async (req, res) => {
    const host = res.locals.host as HostServer;

    if (host.hostStarted) {
      stopHost(req, res);
    } else {
      startHost(req, res);
    }
  });

  app.get('/control/:host/start', startHost);

  app.get('/control/:host/stop', stopHost);

  app.get('/control/:host/disable', (req, res) => {
    const host = res.locals.host as HostServer;

    host.disableAutoStop();

    res.send('Auto stop disabled');
  });

  app.get('/control/:host/enable', (req, res) => {
    const host = res.locals.host as HostServer;

    host.enableAutoStop();

    res.send('Auto stop enabled');
  });

  app.get('/control/:host/autostop', (req, res) => {
    const host = res.locals.host as HostServer;

    log('Autostop : ', host.getAutoStop());

    if (host.getAutoStop()) {
      host.disableAutoStop();

      res.status(569).send('Auto stop disabled');
    } else {
      host.enableAutoStop();

      res.status(200).send('Auto stop enabled');
    }
  });

  app.get('/control/:host/autostop-status', (req, res) => {
    const host = res.locals.host as HostServer;

    res.status(host.getAutoStop() ? 200 : 569).send(host.getAutoStop() ? 'Enabled' : 'Disabled');
  });

  app.get('/control/:host/status', async (req, res) => {
    const host = res.locals.host as HostServer;

    const status = await host.getHostStatus();

    res.status(status ? 200 : 500).send(status ? 'Started' : 'Unreachable');
  });

  app.get('/status', (req, res) => {
    res.send('OK');
  });

  app.listen(apiPort, () => {
    log(`API listening on port ${apiPort}`);
  });

  log('Express API created');
};

main();
