import express from 'express';
import fs from 'fs';
import generator from 'generate-password';
import { Color, logger } from 'logger';
import cron from 'node-cron';
import path from 'path';

import { ArchiveApi } from './archiveApi';
import { CifsApi } from './cifsApi';
import { MailApi } from './mailApi';
import { RsyncApi } from './rsyncApi';
import { SaveApi } from './saveApi';
import { config } from './utils/config';
import { sendBackupRecap } from './utils/ntfy';
import { OsUtils } from './utils/osUtils';
import { CifsDirectory, DirectoryToBackup, DirectoryType } from './utils/types';

logger.setAppName('autosaver');

const { backupConfig } = config;

let lastExecutionSuccess = false;
let isExecuting = false;

const getMountPath = (cifsDirectory: CifsDirectory) =>
  cifsDirectory.mountPath.startsWith('/')
    ? cifsDirectory.mountPath
    : `${backupConfig.backupPath}/${cifsDirectory.mountPath}`;

const run = async () => {
  if (isExecuting) {
    logger.info('Autosaver already running');
    return;
  }

  logger.info('Backup script started');
  isExecuting = true;
  try {
    logger.info('Step 1');
    if (backupConfig.rsync && backupConfig.rsync.type === DirectoryType.cifs) {
      logger.info('Mounting rsync');
      await CifsApi.mountCifsFolder(
        backupConfig.rsync.ip,
        backupConfig.rsync.user,
        backupConfig.rsync.password,
        backupConfig.rsync.hostPath,
        backupConfig.rsync.mountPath,
        backupConfig.rsync.port
      );
    } else {
      logger.info('Rsync is not a cifs directory');
    }

    logger.info('Step 2');
    if (backupConfig.cifsDirectories && backupConfig.cifsDirectories.length > 0) {
      logger.info('Mounting cifsDirectories');
      for (const cifsDirectory of backupConfig.cifsDirectories) {
        await CifsApi.mountCifsFolder(
          cifsDirectory.ip,
          cifsDirectory.user,
          cifsDirectory.password,
          cifsDirectory.hostPath,
          getMountPath(cifsDirectory),
          cifsDirectory.port
        );
      }
    } else {
      logger.info('No cifsDirectories to mount');
    }

    logger.info('Step 3');

    let rsyncPath = null;

    if (backupConfig.rsync) {
      logger.info('Running rsync');

      rsyncPath =
        backupConfig.rsync.type === DirectoryType.local ? backupConfig.rsync.path : backupConfig.rsync.mountPath;

      await RsyncApi.rsyncFolder(backupConfig.backupPath, rsyncPath);
    } else {
      logger.info('Rsync is not enabled');
    }
    logger.info('RsycPath : ', rsyncPath);

    logger.info('Step 4');
    const backupPath = rsyncPath ? path.join(rsyncPath, 'backup') : backupConfig.backupPath;
    logger.info('BackupPath : ', backupPath);

    const directoriesToBackup: DirectoryToBackup[] = fs
      .readdirSync(backupPath)
      .filter((file) => fs.lstatSync(path.join(backupPath, file)).isDirectory())
      .map((directory) => ({ name: directory, path: path.join(backupPath, directory) }));

    logger.info('Folders to backup : ', directoriesToBackup);

    logger.info('Step 5');
    const backupDestPath =
      backupConfig.backupDest.type === DirectoryType.local
        ? backupConfig.backupDest.path
        : backupConfig.backupDest.mountPath;
    if (backupConfig.backupDest.type === DirectoryType.cifs) {
      logger.info('Mounting backupDest');
      await CifsApi.mountCifsFolder(
        backupConfig.backupDest.ip,
        backupConfig.backupDest.user,
        backupConfig.backupDest.password,
        backupConfig.backupDest.hostPath,
        backupConfig.backupDest.mountPath,
        backupConfig.backupDest.port
      );
    } else {
      logger.info('BackupDest is not a cifs directory');
    }
    logger.info('BackupDestPath : ', backupDestPath);

    for (const directory of directoriesToBackup) {
      logger.info('Step 6');

      const archivePassword = generator.generate({ length: 12, numbers: true });

      logger.info(`Archiving ${directory.path}`);
      const { nbFilesArchived, archiveSize, archivePath } = await ArchiveApi.archiveFolder(
        directory.path,
        archivePassword
      );

      logger.info(`Archived ${nbFilesArchived} files (${archiveSize} bytes) to ${archivePath}`);

      logger.info('Step 7');

      await SaveApi.cpFile(archivePath, backupDestPath);

      await OsUtils.rmFiles([archivePath]);

      await MailApi.sendFileInfos(archivePassword, archivePath, `${directory.name}`);

      logger.colored.info(Color.GREEN, `Backup of the directory ${directory.name} done`);

      directory.success = true;
      directory.filesNb = nbFilesArchived;
      directory.size = archiveSize;
    }

    logger.info('Step 8');
    await SaveApi.cleanOldDirectories(backupDestPath);

    logger.info('Step 9');
    sendBackupRecap(directoriesToBackup);
    lastExecutionSuccess = directoriesToBackup.filter((f) => !f.success).length === 0;

    logger.info('Step 10');
    if (backupConfig.backupDest.type === DirectoryType.cifs) {
      logger.info('Unmounting backupDest');
      await CifsApi.unmountSmbFolder(backupConfig.backupDest.mountPath);
    } else {
      logger.info('BackupDest is not a cifs directory');
    }

    logger.info('Step 11');
    if (backupConfig.cifsDirectories && backupConfig.cifsDirectories.length > 0) {
      logger.info('Unmounting cifsDirectories');
      for (const cifsDirectory of backupConfig.cifsDirectories) {
        await CifsApi.unmountSmbFolder(getMountPath(cifsDirectory));
      }
    } else {
      logger.info('No cifsDirectories to unmount');
    }

    logger.info('Step 12');
    if (backupConfig.rsync && backupConfig.rsync.type === DirectoryType.cifs) {
      logger.info('Unmounting rsync');
      await CifsApi.unmountSmbFolder(backupConfig.rsync.mountPath);
    } else {
      logger.info('Rsync is not a cifs directory');
    }

    logger.info('Backup script finished');
  } catch (error) {
    logger.error('Error running autosaver');
    logger.error(error);
  }
};

cron.schedule(config.backupConfig.cronSchedule, run);

logger.info('Script scheduled with the following cron schedule : ', config.backupConfig.cronSchedule);

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
