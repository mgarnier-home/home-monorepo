import dotenv from 'dotenv';
import Fs from 'node:fs';
import Path from 'node:path';
import { parse as ymlParse } from 'yaml';

import { BackupConfig, Config } from './types';

dotenv.config();

const configFilePath = process.env.CONFIG_FILE || Path.resolve(__dirname, '../../config.json');

const loadConfigFromFile = (): Config => {
  const config = Fs.readFileSync(configFilePath, 'utf-8');

  return JSON.parse(config) as Config;
};

const loadConfigFromEnv = (): Config => {
  const backupConfigPath = process.env.BACKUP_CONFIG_PATH || './config.yaml';
  const fullBackupConfigPath = backupConfigPath.startsWith('/')
    ? backupConfigPath
    : Path.join(__dirname, '../../', backupConfigPath);
  const backupConfigYml = Fs.readFileSync(fullBackupConfigPath, 'utf-8');
  const backupConfigJson: BackupConfig = ymlParse(backupConfigYml, { merge: true }).config as BackupConfig;

  const config: Config = {
    serverPort: Number(process.env.SERVER_PORT) || 3000,
    backupConfig: backupConfigJson,
  };

  return config;
};

export const config = (process.env.CONFIG_FILE ? loadConfigFromFile() : loadConfigFromEnv()) as Config;
