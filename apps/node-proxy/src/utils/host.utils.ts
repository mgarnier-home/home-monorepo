import { watch as watchFile } from 'chokidar';
import fs from 'fs';
import jsYaml from 'js-yaml';
import { logger } from 'logger';
import path from 'path';
import { fileURLToPath } from 'url';

import { Host } from '../classes/host.class.js';
import { HostConfig } from './interfaces.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const hosts: Host[] = [];
const configFilePath = path.resolve(__dirname, process.env.CONFIG_FILE ?? '../../config.json');

const loadConfig = async (): Promise<HostConfig[]> => {
  logger.debug('Loading config from : ', configFilePath);
  if (fs.existsSync(configFilePath)) {
    const dataStr = await fs.promises.readFile(configFilePath, 'utf-8');

    if (dataStr !== '') {
      if (configFilePath.endsWith('.json')) {
        return JSON.parse(dataStr) as HostConfig[];
      }
      return jsYaml.load(dataStr) as HostConfig[];
    }
  }
  return [];
};

const getConfig = (cb: (config: HostConfig[]) => void) => {
  watchFile(configFilePath, { ignoreInitial: false }).on('all', async (event, path) => {
    logger.info('Config file changed, triggering callback');

    const config = await loadConfig();

    cb(config);
  });
};

const saveConfig = async (data: HostConfig[]) => {
  const stringData = configFilePath.endsWith('.json') ? JSON.stringify(data, null, 4) : jsYaml.dump(data);

  await fs.promises.writeFile(configFilePath, stringData, 'utf-8');
};

const hostConfigUpdated = () => {
  const data = hosts.map((host) => host.config);

  saveConfig(data);
};

export const getHost = (host: string): Host | undefined => {
  return hosts.find((h) => h.config.name.toLowerCase() === host.toLowerCase());
};

export const setupConfigListenner = () => {
  getConfig(async (configs: HostConfig[]) => {
    await disposeHosts();

    logger.info('Hosts loaded : ', configs);

    hosts.push(...configs.map((config) => new Host(config, hostConfigUpdated)));
  });
};

export const disposeHosts = async () => {
  for (const host of hosts) {
    await host.dispose();

    logger.info('Host disposed : ', host.config.name);
  }

  hosts.splice(0, hosts.length);

  logger.info('All hosts disposed');
};
