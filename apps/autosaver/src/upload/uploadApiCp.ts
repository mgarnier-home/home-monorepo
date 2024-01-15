import cpy from 'cpy';
import fs from 'fs';
import path from 'path';
import { Utils } from 'utils';

import { UploadApi } from './uploadApi.js';

export class UploadApiCp implements UploadApi {
  async cleanOldFolders(backupFolder: string): Promise<void> {
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

    console.log('Folders to delete : ', foldersToDelete);

    for (const folderToDelete of foldersToDelete) {
      await fs.promises.rmdir(path.join(backupFolder, folderToDelete), { recursive: true });

      console.log(`Folder ${folderToDelete} deleted`);
    }
  }

  async protectFile(fileLocation: string): Promise<string> {
    return 'can not protect file with cp';
  }

  async uploadFile(filePath: string, destFolder: string): Promise<string> {
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

    await cpy(filePath, folderPath).on('progress', (progress) => {
      const actualPercentage = Math.floor((progress.completedSize / fileSize) * 100);
      const actualMb = Math.floor(progress.completedSize / 1024 / 1024);

      if (actualPercentage > intPercentage) {
        intPercentage = actualPercentage;
        console.log(`copying file : ${actualMb}Mb / ${totalMb}Mb (${actualPercentage}%)`);
      }
    });

    return destFilePath;
  }
}
