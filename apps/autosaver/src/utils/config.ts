import dotenv from 'dotenv';
import fs from 'node:fs';
import path from 'node:path';
import { cwd } from 'node:process';
import { parse as ymlParse } from 'yaml';

import { ArchiveApiType, BackupConfig, Config } from './types';

dotenv.config();

const configFilePath = process.env.CONFIG_FILE || path.resolve(__dirname, '../../config.json');

const loadConfigFromFile = (): Config => {
  const config = fs.readFileSync(configFilePath, 'utf-8');

  return JSON.parse(config) as Config;
};

const loadConfigFromEnv = (): Config => {
  const backupConfigPath = process.env.BACKUP_CONFIG_PATH || './config.yml';
  const fullBackupConfigPath = backupConfigPath.startsWith('/')
    ? backupConfigPath
    : path.resolve(cwd(), backupConfigPath);

  const config: Config = {
    serverPort: Number(process.env.SERVER_PORT) || 3000,
    backupConfigPath: fullBackupConfigPath,
    cronSchedule: process.env.CRON_SCHEDULE || '0 0 * * *',
    archiveApiType: process.env.ARCHIVE_API_TYPE === ArchiveApiType.ZIP ? ArchiveApiType.ZIP : ArchiveApiType.TAR,
  };

  return config;
};

export const config = (process.env.CONFIG_FILE ? loadConfigFromFile() : loadConfigFromEnv()) as Config;

export const getBackupConfig = (): BackupConfig => {
  const backupConfigYml = fs.readFileSync(config.backupConfigPath, 'utf-8');
  const backupConfig: BackupConfig = ymlParse(backupConfigYml, { merge: true }).config as BackupConfig;

  return backupConfig;
};
