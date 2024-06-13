import fs from 'fs';
import jsYaml from 'js-yaml';
import path from 'path';
import { fileURLToPath } from 'url';

import { Host } from '../classes/host.class';
import { HostConfig } from './interfaces';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const hosts: Host[] = [];
const configFilePath = process.env.CONFIG_FILE || path.resolve(__dirname, '../../config.json');

const loadConfig = async (): Promise<HostConfig[]> => {
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

export const setupHosts = async () => {
  const hostConfigs = await loadConfig();

  hosts.push(...hostConfigs.map((config) => new Host(config, hostConfigUpdated)));
};
