import dotenv from 'dotenv';
import path from 'node:path';

import type { ServerConfig } from '../interfaces/serverConfig';

dotenv.config();

const resolvePath = (pathToResolve: string): string => {
  return pathToResolve.startsWith('/') ? pathToResolve : path.resolve(__dirname, '../', pathToResolve);
};

const loadConfigFromEnv = (): ServerConfig => {
  const config: ServerConfig = {
    appDistPath: resolvePath(process.env.APP_DIST_PATH || './app/browser'),
    configFilePath: resolvePath(process.env.CONFIG_FILE_PATH || 'setup.yml'),
    iconsPath: resolvePath(process.env.ICONS_PATH || './icons'),
    serverPort: Number(process.env.SERVER_PORT) || 3000,
  };

  return config;
};

export const config = loadConfigFromEnv();
