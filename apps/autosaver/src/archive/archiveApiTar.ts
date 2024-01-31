import { spawn } from 'child_process';
import { logger } from 'logger';
import path from 'path';

import { ArchiveApi } from './archiveApi.js';

export class ArchiveApiTar extends ArchiveApi {
  protected async archive(folderToBackupPath: string, totalNbFiles: number) {
    const folderName = path.basename(folderToBackupPath);
    const backupFolder = path.join(folderToBackupPath, '../');
    const tarFileName = `${folderName}.tar.gz`;
    const tarPath = path.join(backupFolder, tarFileName);

    logger.info(`Tarring ${totalNbFiles} files ${folderName} to ${tarPath}`);

    logger.debug('Tar command args : ', ['-c', '-z', '-v', '-f', tarFileName, folderName]);

    let commandArgs = ['-c', '-z', '-v', '-f', tarFileName, folderName];

    return new Promise<string>((resolve, reject) => {
      let filesNb = 0;
      let lastPercent = 0;

      const tarSpawn = spawn(`tar`, commandArgs, { cwd: backupFolder });
      tarSpawn.stdout.on('data', (data: Buffer) => {
        const dataStr = data.toString();

        if (!dataStr.endsWith('/\n')) {
          filesNb++;

          const percent = Math.floor((filesNb / totalNbFiles) * 100);

          if (percent > lastPercent) {
            logger.info(`Tarring ${filesNb}/${totalNbFiles} files: ${percent.toFixed(2)}%`);
            lastPercent = percent;
          }
        }
      });

      tarSpawn.stderr.on('data', (data) => {
        logger.error(`tar stderr: ${data}`);

        tarSpawn.kill();

        reject(new Error(`tar process exited with error`));
      });

      tarSpawn.on('close', (code) => {
        if (code !== 0) {
          reject(new Error(`tar process exited with code ${code}`));
        } else {
          resolve(tarPath);
        }
      });
    });
  }
}
