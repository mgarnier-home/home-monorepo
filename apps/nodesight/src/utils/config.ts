import dotenv from 'dotenv';
import Fs from 'node:fs';
import Path from 'node:path';

import { Config } from './interfaces';

dotenv.config();

const configFilePath = process.env.CONFIG_FILE || Path.resolve(__dirname, '../../config.json');

const loadConfigFromFile = (): Config => {
  const config = Fs.readFileSync(configFilePath, 'utf-8');

  return JSON.parse(config) as Config;
};

const loadConfigFromEnv = (): Config => {
  const config: Config = {
    hostname: process.env.HOSTNAME || 'localhost',
    serverPort: Number(process.env.SERVER_PORT) || 3000,
    updateInterval: Number(process.env.UPDATE_INTERVAL) || 15000,
    enableStatsApi: process.env.ENABLE_STATS_API === 'true',
    statsApiUrl: process.env.STATS_API_URL || 'http://localhost:3000',
    disableCpuTemps: process.env.DISABLE_CPU_TEMPS === 'true',
  };

  return config;
};

export const config = (process.env.CONFIG_FILE ? loadConfigFromFile() : loadConfigFromEnv()) as Config;
