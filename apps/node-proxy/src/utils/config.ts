//the config is coming from a file name config.json
import Fs from 'node:fs';
import Path from 'node:path';

import { Config, Host } from './interfaces.js';

const configFilePath = process.env.CONFIG_FILE || Path.resolve(__dirname, '../../config.json');

const loadConfig = (): Config | undefined => {
  try {
    const config = Fs.readFileSync(configFilePath, 'utf-8');

    return JSON.parse(config) as Config;
  } catch (error) {
    console.error(error);

    return undefined;
  }
};

const saveConfig = (config: Config): boolean => {
  try {
    Fs.writeFileSync(configFilePath, JSON.stringify(config, null, 2), 'utf-8');

    return true;
  } catch (error) {
    console.error(error);

    return false;
  }
};

const updateHostConfig = (host: Host) => {
  const config = loadConfig();

  if (!config) {
    return false;
  }

  const hostIndex = config.hosts.findIndex((h) => h.name === host.name);

  if (hostIndex === -1) {
    return false;
  }

  config.hosts[hostIndex] = host;

  return saveConfig(config);
};

export { loadConfig, updateHostConfig, configFilePath };
