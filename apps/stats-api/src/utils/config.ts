import Fs from 'node:fs';
import Path from 'node:path';

import { getEnvVariable } from '@libs/env-config';

import { Config } from './interfaces';

const loadConfig = (): Config => {
  const config: Config = {
    serverPort: getEnvVariable('SERVER_PORT', false, 3000),
    dbHost: getEnvVariable('DB_HOST', true),
    dbOrg: getEnvVariable('DB_ORG', true),
    dbBucket: getEnvVariable('DB_BUCKET', true),
    dbToken: getEnvVariable('DB_TOKEN', true),
  };

  return config;
};

export const config = loadConfig();
