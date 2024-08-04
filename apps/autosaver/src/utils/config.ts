import fs from 'node:fs';
import path from 'node:path';
import { cwd } from 'node:process';
import * as YAML from 'yaml';

import { getEnvVariable } from '@libs/env-config';

import { ArchiveApiType, BackupConfig, Config } from './types';

const loadConfig = (): Config => {
  const backupConfigPath = getEnvVariable<string>('BACKUP_CONFIG_PATH', false, './config.yml');
  const fullBackupConfigPath = backupConfigPath.startsWith('/')
    ? backupConfigPath
    : path.resolve(cwd(), backupConfigPath);

  const config: Config = {
    serverPort: getEnvVariable('SERVER_PORT', false, 3000),
    backupConfigPath: fullBackupConfigPath,
    cronSchedule: getEnvVariable('CRON_SCHEDULE', false, '0 0 * * *'),
    archiveApiType: getEnvVariable<ArchiveApiType>('ARCHIVE_API_TYPE', false, ArchiveApiType.TAR),
    keepAliveUrl: getEnvVariable<string>('KEEP_ALIVE_URL', false, ''),
  };

  return config;
};

export const config = loadConfig();

export const getBackupConfig = (): BackupConfig => {
  const backupConfigYml = fs.readFileSync(config.backupConfigPath, 'utf-8');
  const backupConfig: BackupConfig = YAML.parse(backupConfigYml, { merge: true })?.config as BackupConfig;

  return backupConfig;
};
