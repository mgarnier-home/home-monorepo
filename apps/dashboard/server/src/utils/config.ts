import dotenv from 'dotenv';
import fs from 'node:fs';
import path from 'node:path';

import type { ServerConfig } from './serverConfig';

dotenv.config();

const configFilePath = process.env.CONFIG_FILE || path.resolve(__dirname, '../../config.json');

const loadConfigFromFile = (): ServerConfig => {
  const config = fs.readFileSync(configFilePath, 'utf-8');

  return JSON.parse(config) as ServerConfig;
};

const resolvePath = (pathToResolve: string): string => {
  return pathToResolve.startsWith('/') ? pathToResolve : path.resolve(__dirname, '../', pathToResolve);
};

const loadConfigFromEnv = (): ServerConfig => {
  const config: ServerConfig = {
    appDistPath: resolvePath(process.env.APP_DIST_PATH || './app-dist'),
    appConfPath: resolvePath(process.env.APP_CONF_PATH || 'conf.yml'),
    iconsPath: resolvePath(process.env.ICONS_PATH || './icons'),
    serverPort: Number(process.env.SERVER_PORT) || 3000,
  };

  return config;
};

export const config = (process.env.CONFIG_FILE ? loadConfigFromFile() : loadConfigFromEnv()) as ServerConfig;
