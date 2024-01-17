import { configDotenv } from 'dotenv';
import { readFileSync } from 'fs';
import * as path from 'path';

import type { ServerConfig } from './interfaces.js';

configDotenv();

const configFilePath = process.env.CONFIG_FILE || path.resolve(__dirname, '../../config.json');

const loadConfigFromFile = (): ServerConfig => {
  const config = readFileSync(configFilePath, 'utf-8');

  return JSON.parse(config) as ServerConfig;
};

const resolvePath = (pathToResolve: string): string => {
  return pathToResolve.startsWith('/') ? pathToResolve : path.resolve(__dirname, '../', pathToResolve);
};

const loadConfigFromEnv = (): ServerConfig => {
  const config: ServerConfig = {
    storagePath: resolvePath(process.env.STORAGE_PATH || '../storage'),
    hostsMap: JSON.parse(process.env.HOSTS_MAP || '{}'),
    maxLogFileSize: Number(process.env.MAX_LOG_FILE_SIZE) || 1024 * 1024 * 256,
    serverPort: Number(process.env.SERVER_PORT) || 3000,
    syslogPort: Number(process.env.SYSLOG_PORT) || 514,
  };

  return config;
};

export const config = (process.env.CONFIG_FILE ? loadConfigFromFile() : loadConfigFromEnv()) as ServerConfig;
