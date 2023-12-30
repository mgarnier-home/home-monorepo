import express from 'express';
import fs from 'fs';

import { changeRedirection } from './freeboxApi.js';
import { config } from './utils/config.js';
import { AppData, TraefikService } from './utils/interfaces.js';
import { listFiles, readFiles } from './utils/osUtils.js';
import { getComposeStacksPaths, getTraefikDynamicConf, getTraefikServices } from './utils/traefikUtils.js';
import { mergeYamls } from './utils/utils.js';

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

let appData = await loadData();

console.log('appData : ', appData);

const app = express();

app.use((err: any, req: any, res: any, next: any) => {
  console.error(err.stack);
  res.status(500).send('Something broke!');
});

app.get('/', (req, res) => {
  console.log('root');
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
  console.log('proxy-enable');

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
  console.log('proxy-status');

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
  console.log('Server started on port ' + config.serverPort);
});

// const success1 = await changeRedirection("ovh", "192.168.0.100");

// console.log("redirection to athena : ", success1);

// const success2 = await changeRedirection("ovh", "192.168.0.98");

// console.log("redirection to artemis : ", success2);

// const freeboxRegister = new freebox.FreeboxRegister({
//   app_id: "mgarnier11.traefik_conf",
//   app_name: "Traefik conf",
//   app_version: "1.0.0",
//   device_name: "Server",
// });

// const access = await freeboxRegister.register();

// console.log(access);
