import dotenv from 'dotenv';
//the config is coming from a file name config.json
import Fs from 'node:fs';
import Path from 'node:path';
import { fileURLToPath } from 'node:url';

import { ArchiveStrategy, Config, SmbConfig, UploadStrategy } from './interfaces.js';

dotenv.config();

export const __filename = fileURLToPath(import.meta.url);

export const __dirname = Path.dirname(__filename);

const configFilePath = process.env.CONFIG_FILE || Path.resolve(__dirname, '../../config.json');

const loadConfigFromFile = (): Config => {
  const config = Fs.readFileSync(configFilePath, 'utf-8');

  return JSON.parse(config) as Config;
};

const loadConfigFromEnv = (): Config => {
  const backupDir = process.env.FOLDER_TO_BACKUP || 'backups';
  const folderToBackup = backupDir.startsWith('/') ? backupDir : Path.join(__dirname, '../../', backupDir);

  const rsyncDir = process.env.RSYNC_DEST_FOLDER || 'rsynced';
  const rsyncDestFolder = rsyncDir.startsWith('/') ? rsyncDir : Path.join(__dirname, '../../', rsyncDir);

  let smbConfig = process.env.SMB_CONFIG || '';
  if (smbConfig.endsWith('.json')) {
    const smbConfigPath = smbConfig.startsWith('/') ? smbConfig : Path.join(__dirname, '../../', smbConfig);
    smbConfig = Fs.readFileSync(smbConfigPath, 'utf-8');
  }
  const smbConfigJson: SmbConfig = JSON.parse(smbConfig || '{}');

  const config: Config = {
    enableMail: process.env.ENABLE_MAIL === 'true' || false,
    enableRsync: process.env.ENABLE_RSYNC === 'true' || false,
    uploadStrategy: (process.env.UPLOAD_STRATEGY as any) || UploadStrategy.CP,
    archiveStrategy: (process.env.ARCHIVE_STRATEGY as any) || ArchiveStrategy.ZIP,
    mailHost: process.env.MAIL_HOST || '',
    mailPort: Number(process.env.MAIL_PORT) || 465,
    mailSecure: process.env.MAIL_SECURE === 'true' || true,
    infoMailTo: process.env.INFO_MAIL_TO || '',
    errorMailTo: process.env.ERROR_MAIL_TO || '',
    mailLogin: process.env.MAIL_LOGIN || '',
    mailPassword: process.env.MAIL_PASSWORD || '',
    uptoboxToken: process.env.UPTOBOX_TOKEN || '',
    cronSchedule: process.env.CRON_SCHEDULE || '0 4 * * *',
    serverPort: Number(process.env.SERVER_PORT) || 3000,
    folderToBackup,
    rsyncDestFolder,
    uploadDestFolder: process.env.UPLOAD_DEST_FOLDER || 'backups',
    deleteFiles: process.env.DELETE_FILES === 'true' || true,
    smbConfig: smbConfigJson,
  };

  return config;
};

export const config = (process.env.CONFIG_FILE ? loadConfigFromFile() : loadConfigFromEnv()) as Config;
