import { spawn } from 'child_process';
import fs from 'fs';
import path from 'path';

import { logger } from '@libs/logger';

import { EncryptApi } from './encryptApi';
import { config } from './utils/config';
import { OsUtils } from './utils/osUtils';
import { ArchiveApiType } from './utils/types';

export type ArchiveFn = (folderToBackupPath: string, totalNbFiles: number) => Promise<string>;

const archiveTar: ArchiveFn = (folderToBackupPath: string, totalNbFiles: number) => {
  const folderName = path.basename(folderToBackupPath);
  const backupFolder = path.join(folderToBackupPath, '../');
  const tarFileName = `${folderName}.tar.gz`;
  const tarPath = path.join(backupFolder, tarFileName);

  logger.info(`Tarring ${totalNbFiles} files ${folderName} to ${tarPath}`);

  if (fs.existsSync(tarPath)) {
    logger.info(`Removing existing tar file ${tarPath}`);
    fs.rmSync(tarPath);
  }

  let commandArgs = [
    '--ignore-failed-read',
    '--warning=no-file-changed',
    '-c',
    '-z',
    '-v',
    '-f',
    tarFileName,
    folderName,
  ];

  return new Promise<string>((resolve, reject) => {
    let filesNb = 0;
    let lastPercent = 0;

    logger.debug(`running command: tar ${commandArgs.join(' ')}`);

    const tarSpawn = spawn(`tar`, commandArgs, { cwd: backupFolder });
    tarSpawn.stdout.on('data', (data: Buffer) => {
      const dataStr = data.toString();

      logger.debug(dataStr);

      if (!dataStr.endsWith('/\n')) {
        filesNb++;

        const percent = Math.floor((filesNb / totalNbFiles) * 100);

        if (percent > lastPercent) {
          logger.info(`Tarring ${filesNb}/${totalNbFiles} files: ${percent.toFixed(2)}%`);
          lastPercent = percent;
        }
      }
    });

    const errorBuffers: Buffer[] = [];

    tarSpawn.stderr.on('data', (data: Buffer) => {
      errorBuffers.push(data);

      tarSpawn.kill();
    });

    tarSpawn.on('close', (code) => {
      if (code !== 0) {
        const error = Buffer.concat(errorBuffers).toString();

        reject({ message: `tar process exited with code ${code}`, code, error });
      } else {
        resolve(tarPath);
      }
    });
  });
};

const archiveZip: ArchiveFn = async (folderToBackupPath: string, totalNbFiles: number) => {
  const folderName = path.basename(folderToBackupPath);
  const backupFolder = path.join(folderToBackupPath, '../');
  const zipFileName = `${folderName}.zip`;
  const zipPath = path.join(backupFolder, zipFileName);

  logger.info(`Zipping ${totalNbFiles} files from ${folderName} to ${zipPath}`);

  if (fs.existsSync(zipPath)) {
    logger.info(`Removing existing zip file ${zipPath}`);
    fs.rmSync(zipPath);
  }

  let commandArgs = ['a', '-bt', '-bb3', '-tzip', '-mmt=on', zipPath, folderToBackupPath];

  return new Promise<string>((resolve, reject) => {
    logger.debug(`running command: 7z ${commandArgs.join(' ')}`);

    const zipSpawn = spawn('7z', commandArgs);
    let filesNb = 0;

    const newFileRegex = new RegExp(`(U|\\+) ${folderName}`, 'gm');

    let lastPercent = 0;

    let lastLine = '';

    // Live stdout logging
    zipSpawn.stdout.on('data', (data: Buffer) => {
      let dataStr = lastLine + data.toString();

      if (!dataStr.endsWith('\n')) {
        lastLine = dataStr.substring(dataStr.lastIndexOf('\n') + 1);
        dataStr = dataStr.substring(0, dataStr.length - lastLine.length);
      } else {
        lastLine = '';
      }

      const newFiles = dataStr.match(newFileRegex);
      if (newFiles) {
        filesNb += newFiles.length;
        const percent = Math.floor((filesNb / totalNbFiles) * 100);
        if (percent > lastPercent) {
          lastPercent = percent;
          logger.info(`Zipping ${filesNb}/${totalNbFiles} files: ${percent}%`);
        }
      } else {
        logger.info('stdout: ', dataStr);
      }
    });

    zipSpawn.stderr.on('data', (data) => {
      logger.error(`zip stderr: ${data}`);

      zipSpawn.kill();

      reject(new Error(`7z process exited with error`));
    });

    zipSpawn.on('close', (code) => {
      if (code !== 0) {
        reject(new Error(`7z process exited with code ${code}`));
      } else {
        resolve(zipPath);
      }
    });
  });
};

const archive = {
  withFn: (archiveFn: ArchiveFn) => {
    return {
      archiveFolder: async (
        folderToBackupPath: string,
        archivePassword?: string
      ): Promise<{ archivePath: string; nbFilesArchived: number; archiveSize: number }> => {
        const nbFiles = await OsUtils.osCountFiles(folderToBackupPath);

        let archivePath = await archiveFn(folderToBackupPath, nbFiles);

        logger.info(`Archive created : ${archivePath}`);

        if (archivePassword) {
          logger.info(`Encrypting archive : ${archivePath}...`);

          archivePath = await EncryptApi.encryptFile(archivePath, archivePassword);

          logger.info(`Archive encrypted : ${archivePath}`);

          await OsUtils.rmFiles([archivePath.replace('.gpg', '')]);
        }

        const archiveSize = (await fs.promises.stat(archivePath)).size;

        return { archivePath, nbFilesArchived: nbFiles, archiveSize };
      },
    };
  },
};

export const archiveApi = archive.withFn(config.archiveApiType === ArchiveApiType.TAR ? archiveTar : archiveZip);
