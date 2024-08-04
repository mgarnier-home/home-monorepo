import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { cwd } from 'node:process';
import * as YAML from 'yaml';

import { getEnvVariable } from '@libs/env-config';

import { Config, CronConfig } from './interfaces';

const loadConfig = (): Config => {
  const cronConfigPath = getEnvVariable<string>('CRON_CONFIG_PATH', false, './config.yml');
  const fullCronConfigPath = cronConfigPath.startsWith('/') ? cronConfigPath : resolve(cwd(), cronConfigPath);

  const config: Config = {
    serverPort: getEnvVariable('SERVER_PORT', false, 3000),
    dataFilePath: fullCronConfigPath,
  };

  return config;
};

export const config = loadConfig();

export const getCronConfig = (): CronConfig => {
  const cronConfigYml = readFileSync(config.dataFilePath, 'utf-8');
  const cronConfig: CronConfig = YAML.parse(cronConfigYml, { merge: true }) as CronConfig;

  return cronConfig;
};
