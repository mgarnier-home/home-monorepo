import cpy from 'cpy';
import fs from 'fs';
import path from 'path';

import { logger } from '@libs/logger';
import { Utils } from '@libs/utils';

export namespace SaveApi {
  export const cleanOldDirectories = async (backupFolder: string): Promise<void> => {
    const folderList = (await fs.promises.readdir(backupFolder)).filter((f) =>
      fs.lstatSync(path.join(backupFolder, f)).isDirectory()
    );

    const foldersToDelete = folderList.filter((f) => {
      const folderDate = new Date(f);
      const now = new Date();
      const diff = now.getTime() - folderDate.getTime();
      const diffInDays = diff / (1000 * 3600 * 24);
      return diffInDays > 14;
    });

    logger.info('Folders to delete : ', foldersToDelete);

    for (const folderToDelete of foldersToDelete) {
      await fs.promises.rmdir(path.join(backupFolder, folderToDelete), { recursive: true });

      logger.info(`Folder ${folderToDelete} deleted`);
    }
  };

  export const cpFile = async (filePath: string, destFolder: string): Promise<string> => {
    const archiveFileName = path.basename(filePath);

    const date = new Date();
    const folderPath = path.join(
      destFolder,
      `/${
        date.getFullYear() //
      }-${
        Utils.padStart(date.getMonth() + 1, 2) //
      }-${
        Utils.padStart(date.getDate(), 2) //
      }`
    );

    if (!fs.existsSync(folderPath)) {
      fs.mkdirSync(folderPath, { recursive: true });
    }

    const fileSize = (await fs.promises.stat(filePath)).size;
    const totalMb = Math.floor(fileSize / 1024 / 1024);

    const destFilePath = path.join(folderPath, archiveFileName);

    let intPercentage = 0;

    logger.info(`Copying ${filePath} to ${destFilePath}`);

    await cpy(filePath, folderPath).on('progress', (progress) => {
      const actualPercentage = Math.floor((progress.completedSize / fileSize) * 100);
      const actualMb = Math.floor(progress.completedSize / 1024 / 1024);

      if (actualPercentage > intPercentage) {
        intPercentage = actualPercentage;
        logger.info(`Copying file : ${actualMb}Mb / ${totalMb}Mb (${actualPercentage}%)`);
      }
    });

    logger.info(`File copied to ${destFilePath}`);

    return destFilePath;
  };
}
