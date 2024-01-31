import { spawn } from 'child_process';
import { logger } from 'logger';
import path from 'path';

import { ArchiveApi } from './archiveApi.js';

export class ArchiveApiZip extends ArchiveApi {
  protected async archive(folderToBackupPath: string, totalNbFiles: number) {
    const folderName = path.basename(folderToBackupPath);
    const backupFolder = path.join(folderToBackupPath, '../');
    const zipFileName = `${folderName}.zip`;
    const zipPath = path.join(backupFolder, zipFileName);

    logger.info(`Zipping ${totalNbFiles} files from ${folderName} to ${zipPath}`);

    let commandArgs = ['a', '-bt', '-bb3', '-tzip', '-mmt=on', zipPath, folderToBackupPath];

    return new Promise<string>((resolve, reject) => {
      const zipSpawn = spawn('7z', commandArgs);
      let filesNb = 0;

      const newFileRegex = new RegExp(`(U|\\+) ${folderName}`, 'gm');

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

          logger.info(`Zipping ${filesNb}/${totalNbFiles} files: ${((filesNb / totalNbFiles) * 100).toFixed(2)}%`);
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
  }
}
