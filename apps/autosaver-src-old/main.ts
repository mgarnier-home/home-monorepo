import archiver from 'archiver';
import express from 'express';
import fs from 'fs';
import { Color, logger } from 'logger';
import cron from 'node-cron';
import path from 'path';

import { getArchiveApi } from './archive/index.js';
import { MailApi } from './mail/mailApi.js';
import { RsyncApi } from './rsync/rsyncApi.js';
import { SmbApi } from './smb/smbApi.js';
import { getUploadApi } from './upload/index.js';
import { config } from './utils/config.js';
import { FolderToBackup } from './utils/interfaces.js';
import { sendBackupRecap, sendError } from './utils/ntfy.js';
import { OsUtils } from './utils/osUtils.js';
import { Utils } from './utils/utils.js';

if (!process.env.TZ) {
  process.env.TZ = 'Europe/Paris';
}

logger.setAppName('autosaver');

logger.debug(process.env);
logger.debug(config);

let lastExecutionSuccess = false;
let isExecuting = false;

const { uploadApi, uploadDestFolder } = getUploadApi(config.uploadStrategy, config.uploadDestFolder);
const archiveApi = getArchiveApi(config.archiveStrategy);
const backupsFolder = config.folderToBackup;
const rsyncDestFolderRoot = config.rsyncDestFolder;

let foldersToBackup: FolderToBackup[] = [];

const run = async () => {
  if (!isExecuting) {
    logger.info('Starting backup script');
    isExecuting = true;

    try {
      await SmbApi.mountSmbFolders(config.smbConfig);

      foldersToBackup = fs
        .readdirSync(backupsFolder)
        .filter((f) => fs.lstatSync(path.join(backupsFolder, f)).isDirectory())
        .map((f) => ({ name: f, path: path.join(backupsFolder, f) }));

      logger.info('Folders to backup : ', foldersToBackup);

      if (config.enableRsync) {
        for (const folderToBackup of foldersToBackup) {
          const folderToBackupPath = path.join(backupsFolder, folderToBackup.name);

          logger.info(`Rsyncing folder ${folderToBackupPath} to ${rsyncDestFolderRoot}`);

          await RsyncApi.rsyncFolder(folderToBackupPath, rsyncDestFolderRoot);

          folderToBackup.path = path.join(rsyncDestFolderRoot, folderToBackup.name);

          logger.info(`Rsync of the folder ${folderToBackup.name} done`);
        }
      }

      for (const folderToBackup of foldersToBackup) {
        try {
          const archivePassword = Utils.generatePassword();

          const { nbFilesArchived, archiveSize, archivePath } = await archiveApi.archiveFolder(
            folderToBackup.path,
            archivePassword
          );

          const uploadedFile = await uploadApi.uploadFile(archivePath, uploadDestFolder);

          const filePassword = await uploadApi.protectFile(uploadedFile);

          await MailApi.sendFileInfos(filePassword, archivePassword, uploadedFile, `${folderToBackup.name}.zip`);

          await OsUtils.rmFiles([archivePath]);

          logger.colored.info(Color.GREEN, `Backup of the folder ${folderToBackup.name} done`);

          folderToBackup.success = true;
          folderToBackup.filesNb = nbFilesArchived;
          folderToBackup.size = archiveSize;
        } catch (err) {
          logger.error(err);

          await MailApi.sendError(`An error happened during the backup of the file ${folderToBackup}.zip : ${err}`);

          folderToBackup.success = false;
        }
      }

      await uploadApi.cleanOldFolders(config.uploadDestFolder);

      sendBackupRecap(foldersToBackup);

      lastExecutionSuccess = foldersToBackup.filter((f) => !f.success).length === 0;
    } catch (err) {
      logger.error(err);

      await MailApi.sendError(`An error happened reading the folders to backup : ${err}`);

      sendError(err);

      lastExecutionSuccess = false;
    } finally {
      logger.info('Backup script finished');

      Utils.printRecapTable(foldersToBackup);

      await SmbApi.unmountSmbFolders(config.smbConfig);

      isExecuting = false;

      foldersToBackup = [];
    }
  } else {
    logger.warn('Backup script already running');
  }
};

cron.schedule(config.cronSchedule, run);

logger.info('Script scheduled with the following cron schedule : ', config.cronSchedule);

//start an express server
const app = express();

app.get('/', (req, res) => {
  if (isExecuting) {
    res.status(204).send('OK');
  } else {
    res.status(200).send('OK');
  }
});

app.get('/last', (req, res) => {
  res.status(lastExecutionSuccess ? 200 : 500).send(lastExecutionSuccess ? 'OK' : 'KO');
});

app.get('/run', (req, res) => {
  if (isExecuting) {
    res.status(204).send('Autosave already running');
  } else {
    run();
    res.status(200).send('Autosave started');
  }
});

app.listen(config.serverPort, () => {
  logger.info('Server started on port ' + config.serverPort);
});
