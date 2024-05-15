import { Docker } from 'docker-api';
import express from 'express';
import fs from 'fs';
import { logger } from 'logger';

import { config } from './utils/config.js';
import { AppData, TraefikService } from './utils/interfaces.js';
import { getTraefikDynamicConf } from './utils/traefikUtils.js';

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
    console.error('Error while loading data : ', error);
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

    const traefikServices: TraefikService[] = [];

    for (const host of appData.hosts) {
      const docker = new Docker(`${host.ip}`, host.apiPort);
      const containers = await docker.listContainers();
      traefikServices.push(
        ...containers
          .filter((container) => {
            return container.Labels['traefik-conf.port'] != undefined;
          })
          .map((container) => {
            const portVariable = container.Labels['traefik-conf.port'] || '';
            const serviceName =
              container.Labels['traefik-conf.name'] || container.Labels['com.docker.compose.service'] || '';

            return { host, serviceName, portVariable };
          })
      );
    }

    const dynamicYml = getTraefikDynamicConf(traefikServices, appData);

    res.status(200).send(dynamicYml);
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
