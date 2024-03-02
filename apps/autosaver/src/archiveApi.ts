import { spawn } from 'child_process';
import fs from 'fs';
import { logger } from 'logger';
import path from 'path';

import { EncryptApi } from './encryptApi';
import { OsUtils } from './utils/osUtils';

export namespace ArchiveApi {
  export const archiveFolder = async (
    folderToBackupPath: string,
    archivePassword?: string
  ): Promise<{ archivePath: string; nbFilesArchived: number; archiveSize: number }> => {
    const nbFiles = await OsUtils.osCountFiles(folderToBackupPath);

    let archivePath = await archive(folderToBackupPath, nbFiles);

    logger.info(`Archive created : ${archivePath}`);

    if (archivePassword) {
      logger.info(`Encrypting archive : ${archivePath}...`);

      archivePath = await EncryptApi.encryptFile(archivePath, archivePassword);

      logger.info(`Archive encrypted : ${archivePath}`);

      await OsUtils.rmFiles([archivePath.replace('.gpg', '')]);
    }

    const archiveSize = (await fs.promises.stat(archivePath)).size;

    return { archivePath, nbFilesArchived: nbFiles, archiveSize };
  };

  const archive = (folderToBackupPath: string, totalNbFiles: number) => {
    const folderName = path.basename(folderToBackupPath);
    const backupFolder = path.join(folderToBackupPath, '../');
    const tarFileName = `${folderName}.tar.gz`;
    const tarPath = path.join(backupFolder, tarFileName);

    logger.info(`Tarring ${totalNbFiles} files ${folderName} to ${tarPath}`);

    let commandArgs = ['--ignore-failed-read', '-c', '-z', '-v', '-f', tarFileName, folderName];

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
}
