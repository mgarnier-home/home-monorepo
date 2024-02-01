import express from 'express';
import fs from 'fs';
import { logger } from 'logger';

import { changeRedirection } from './freeboxApi.js';
import { config } from './utils/config.js';
import { AppData, TraefikService } from './utils/interfaces.js';
import { listFiles, readFiles } from './utils/osUtils.js';
import { getComposeStacksPaths, getTraefikDynamicConf, getTraefikServices } from './utils/traefikUtils.js';
import { mergeYamls } from './utils/utils.js';

logger.setAppName('traefik-conf');

const saveData = async (data: AppData) => {
  await fs.promises.writeFile(config.saveDataFile, JSON.stringify(data, null, 4));
};

const loadData = async (): Promise<AppData> => {
  if (fs.existsSync(config.saveDataFile)) {
    const dataStr = await fs.promises.readFile(config.saveDataFile, 'utf-8');

    if (dataStr !== '') {
      return JSON.parse(dataStr) as AppData;
    }
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

  app.get('/config', async (req, res) => {
    const ymlFilesPaths = await listFiles(config.traefikConfDirectory, '.yml');

    const ymlFilesContent = await readFiles(ymlFilesPaths);
    const mergedYaml = mergeYamls(ymlFilesContent);

    res.status(200).send(mergedYaml);
  });

  app.get('/compose-config', async (req, res) => {
    appData = await loadData();

    const stacksInfos = await getComposeStacksPaths(appData.hosts, config.composeDirectory, config.stacksToIgnore);

    const traefikServices: TraefikService[] = [];

    for (const stackInfos of stacksInfos) {
      const fileContent = await fs.promises.readFile(stackInfos.path, 'utf-8');

      traefikServices.push(...getTraefikServices(stackInfos, fileContent));
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

      await changeRedirection(config.redirectionName, ip);

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
