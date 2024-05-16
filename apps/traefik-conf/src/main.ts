import { Docker } from 'docker-api';
import express from 'express';
import fs from 'fs';
import jsYaml from 'js-yaml';
import { logger } from 'logger';

import { config } from './utils/config.js';
import { AppData } from './utils/interfaces.js';
import { parseTraefikLabels } from './utils/traefikUtils.js';

logger.setAppName('traefik-conf');

const saveData = async (data: AppData) => {
  await fs.promises.writeFile(config.dataFilePath, JSON.stringify(data, null, 4));
};

const loadData = async (): Promise<AppData> => {
  try {
    if (fs.existsSync(config.dataFilePath)) {
      const dataStr = await fs.promises.readFile(config.dataFilePath, 'utf-8');

      if (dataStr !== '') {
        return JSON.parse(dataStr) as AppData;
      }
    }
  } catch (error) {
    logger.error('Error while loading data : ', error);
  }

  return {
    proxies: [],
    hosts: [],
  };
};

const main = async () => {
  let appData = await loadData();

  logger.info('appData : ', appData);

  const app = express();

  app.use((err: any, req: any, res: any, next: any) => {
    logger.error(err.stack);
    res.status(500).send('Something broke!');
  });

  app.get('/', (req, res) => {
    res.status(200).send('OK');
  });

  app.get('/dynamic-config', async (req, res) => {
    appData = await loadData();

    const traefikConf: any = { http: { services: {}, routers: {} } };

    for (const host of appData.hosts) {
      const docker = new Docker(`${host.ip}`, host.apiPort);
      const containers = await docker.listContainers();

      for (const container of containers) {
        const result = parseTraefikLabels(host, container.Labels, appData);

        if (result.services) {
          traefikConf.http.services = { ...traefikConf.http.services, ...result.services };
        }

        if (result.routers) {
          traefikConf.http.routers = { ...traefikConf.http.routers, ...result.routers };
        }
      }
    }

    res.status(200).send(jsYaml.dump(traefikConf, { indent: 2 }));
  });

  app.get('/proxy-enable/:proxyIndex', async (req, res) => {
    logger.info('proxy-enable');

    const proxyIndex = Number(req.params.proxyIndex);
    const proxy = appData.proxies[proxyIndex];

    if (proxy) {
      const ip = proxy.activated ? proxy.sourceIP : proxy.destIP;

      proxy.activated = !proxy.activated;

      await saveData(appData);

      res.status(200).send('Ip set to ' + ip + ' enabled : ' + proxy.activated);
    } else {
      res.status(404).send('Proxy not found');
    }
  });

  app.get('/proxy-status/:proxyIndex', async (req, res) => {
    logger.info('proxy-status');

    const proxyIndex = Number(req.params.proxyIndex);
    const proxy = appData.proxies[proxyIndex];

    if (proxy) {
      if (proxy.activated) {
        res.status(200).send('Proxy activated');
      } else {
        res.status(500).send('Proxy not activated');
      }
    } else {
      res.status(404).send('Proxy not found');
    }
  });

  app.listen(config.serverPort, () => {
    logger.info('Server started on port ' + config.serverPort);
  });
};

main();
