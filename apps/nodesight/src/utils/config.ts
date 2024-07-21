import Fs from 'node:fs';
import Path from 'node:path';

import { getEnvVariable } from '@libs/env-config';

import { Config } from './interfaces';

const loadConfig = (): Config => {
  const enableStatsApi = getEnvVariable('ENABLE_STATS_API', false, true);

  const config: Config = {
    hostname: getEnvVariable('HOSTNAME', true),
    serverPort: getEnvVariable('SERVER_PORT', false, 3000),
    updateInterval: getEnvVariable('UPDATE_INTERVAL', false, 15000),
    enableStatsApi,
    statsApiUrl: getEnvVariable('STATS_API_URL', enableStatsApi),
    disableCpuTemps: getEnvVariable('DISABLE_CPU_TEMPS', false, false),
  };

  return config;
};

export const config = loadConfig();
