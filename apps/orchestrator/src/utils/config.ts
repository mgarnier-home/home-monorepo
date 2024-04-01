import dotenv from 'dotenv';
import fs from 'node:fs';
import path from 'node:path';
import { Utils } from 'utils';
import { parse as ymlParse } from 'yaml';

import { Config, Stack, stackSchema } from './schemas';

dotenv.config();

const configFilePath = process.env.CONFIG_FILE || path.resolve(__dirname, '../../config.json');

const loadConfigFromFile = (): Config => {
  const config = fs.readFileSync(configFilePath, 'utf-8');

  return JSON.parse(config) as Config;
};

const loadConfigFromEnv = (): Config => {
  const dirname = __dirname;
  const config: Config = {
    serverPort: parseInt(process.env.SERVER_PORT || '3000'),
    composeEnvFilesPaths: process.env.COMPOSE_ENV_FILES?.split(',').map((file) =>
      Utils.resolveConfigPath(file.trim(), path.resolve(dirname, '../')).trim()
    ) || ['.env'],
    composeFolderPath: Utils.resolveConfigPath(process.env.COMPOSE_FOLDER || '../compose', dirname),
    stackFilePath: Utils.resolveConfigPath(process.env.STACK_FILE || '../stack.yml', dirname),
  };

  return config;
};

export const config = (process.env.CONFIG_FILE ? loadConfigFromFile() : loadConfigFromEnv()) as Config;

export const getStack = (): Stack => {
  const backupConfigYml = fs.readFileSync(config.stackFilePath, 'utf-8');
  const backupConfig: Stack = stackSchema.parse(ymlParse(backupConfigYml, { merge: true }));

  return backupConfig;
};
