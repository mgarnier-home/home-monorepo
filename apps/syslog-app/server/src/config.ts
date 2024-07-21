import * as path from 'path';

import { getEnvVariable } from '@libs/env-config';

import type { ServerConfig } from './interfaces';

const loadConfig = (): ServerConfig => {
  const config: ServerConfig = {
    devMode: process.env.NODE_ENV !== 'production',
    storagePath: path.resolve(__dirname, '../', getEnvVariable('STORAGE_PATH', false, '../storage')),
    hostsMap: JSON.parse(process.env.HOSTS_MAP || '{}'),
    maxLogFileSize: getEnvVariable('MAX_LOG_FILE_SIZE', false, 1024 * 1024 * 256),
    serverPort: getEnvVariable('SERVER_PORT', false, 3000),
    syslogPort: getEnvVariable('SYSLOG_PORT', false, 514),
  };

  return config;
};

export const config = loadConfig();
