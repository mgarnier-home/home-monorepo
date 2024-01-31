import archiver from 'archiver';
import fs from 'fs';
import gracefulFs from 'graceful-fs';
import { logger } from 'logger';
import path from 'path';

import { OsUtils } from '../utils/osUtils.js';
import { ArchiveApi } from './archiveApi.js';

export class ArchiveApiArchiver extends ArchiveApi {
  protected async archive(folderToBackupPath: string, totalNbFiles: number) {
    const folderName = path.basename(folderToBackupPath);
    const backupFolder = path.join(folderToBackupPath, '../');
    const zipFileName = `${folderName}.zip`;
    const zipPath = path.join(backupFolder, zipFileName);

    logger.info(`Zipping ${folderName} to ${zipPath}`);

    const output = fs.createWriteStream(zipPath);

    output.on('close', function () {
      logger.info(`${folderName} in progress : ${(archive.pointer() / 1024 / 1024).toFixed(0)} MB`);
    });

    logger.info('Getting files in folder...');

    const files = await OsUtils.getFilesInFolder(folderToBackupPath);

    logger.info('Files in folder : ', files.length);

    const archive = archiver.create('zip', {
      zlib: { level: 9 }, // Sets the compression level.
    });

    archive.pipe(output);

    archive.on('progress', (progress) => {
      logger.info(
        `file added to archive : ${progress.entries.processed}/${files.length} : ${(
          (progress.entries.processed / files.length) *
          100
        ).toFixed(2)}%`
      );
    });

    const concurrentFiles = 50; // Adjust this based on your system's capability
    let processingFiles = 0;
    const filesQueue = [...files]; // Create a copy of your files array

    return new Promise<string>((resolve, reject) => {
      const processNextFile = async () => {
        if (!filesQueue.length && processingFiles === 0) {
          // If there are no more files to process and no files currently processing, finalize the archive
          await archive.finalize();

          resolve(zipPath);
          return;
        }

        while (processingFiles < concurrentFiles && filesQueue.length) {
          const file = filesQueue.shift()!;
          processingFiles++;

          const readStream = gracefulFs.createReadStream(file.path);

          readStream.on('end', () => {
            processingFiles--;
            processNextFile(); // Process the next file once the current one is done
          });

          archive.append(readStream, { name: path.relative(folderToBackupPath, file.path), stats: file.stat });
        }
      };

      processNextFile(); // Start processing the files
    });
  }
}
