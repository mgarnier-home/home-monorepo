import { logger } from 'logger';
//the config is coming from a file name config.json
import Fs from 'node:fs';
import Path from 'node:path';
import { fileURLToPath } from 'node:url';

import { Config, Host } from './interfaces.js';

export const __filename = fileURLToPath(import.meta.url);

export const __dirname = Path.dirname(__filename);

const configFilePath = process.env.CONFIG_FILE || Path.resolve(__dirname, '../../config.json');

const loadConfig = (): Config | undefined => {
  try {
    const config = Fs.readFileSync(configFilePath, 'utf-8');

    return JSON.parse(config) as Config;
  } catch (error) {
    logger.error(error);

    return undefined;
  }
};

const saveConfig = (config: Config): boolean => {
  try {
    Fs.writeFileSync(configFilePath, JSON.stringify(config, null, 2), 'utf-8');

    return true;
  } catch (error) {
    logger.error(error);

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
