import Fs from 'node:fs';
import Path from 'node:path';

import { getEnvVariable } from '@libs/env-config';

import { Config } from './interfaces';

const loadConfig = (): Config => {
  const config: Config = {
    dataFilePath: Path.resolve(__dirname, getEnvVariable('DATA_FILE_PATH', false, '../data.json')),
    serverPort: getEnvVariable('SERVER_PORT', false, 3000),
    secret: getEnvVariable('SECRET', true),
  };

  return config;
};

export const config = loadConfig();
