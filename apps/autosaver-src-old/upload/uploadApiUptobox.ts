import { logger } from 'logger';
import path from 'path';
import { Utils } from 'utils';

import { UptoboxApi } from '../uptobox/uptoboxApi.js';
import { UploadApi } from './uploadApi.js';

export class UploadApiUptobox implements UploadApi {
  async cleanOldFolders(backupFolder: string): Promise<void> {
    const uptoboxBackupsFolder = await UptoboxApi.getFolder(backupFolder);

    const allFolders = (await UptoboxApi.listContent(uptoboxBackupsFolder)).folders;

    const foldersToDelete = allFolders.filter((f) => {
      const folderDate = new Date(f.name);
      const now = new Date();
      const diff = now.getTime() - folderDate.getTime();
      const diffInDays = diff / (1000 * 3600 * 24);
      return diffInDays > 14;
    });

    logger.info(
      'Folders to delete : ',
      foldersToDelete.map((f) => f.fld_name)
    );

    for (const folderToDelete of foldersToDelete) {
      await UptoboxApi.deleteFolder(folderToDelete);

      logger.info(`Folder ${folderToDelete.fld_name} deleted`);
    }
  }

  async uploadFile(filePath: string, destFolder: string): Promise<string> {
    const archiveFileName = path.basename(filePath);

    const date = new Date();
    const uptoboxFolderPath = `${destFolder}/${
      date.getFullYear() //
    }-${
      Utils.padStart(date.getMonth() + 1, 2) //
    }-${
      Utils.padStart(date.getDate(), 2) //
    }`;

    //upload the file
    const fileUrl = await UptoboxApi.uploadFileWithRetry(filePath, 5);

    //retrieve the folder
    const uptoboxDestFolder = await UptoboxApi.getFolder(uptoboxFolderPath);
    //list the files in the folder
    const filesInFolder = (await UptoboxApi.listContent(uptoboxDestFolder)).files;

    //delete the old file if it exists
    for (const file of filesInFolder) {
      if (file.file_name === archiveFileName) {
        await UptoboxApi.deleteFiles([file]);
      }
    }

    //move the file to the folder
    await UptoboxApi.moveFile(fileUrl, uptoboxDestFolder);

    return fileUrl;
  }

  protectFile(fileLocation: string): Promise<string> {
    return UptoboxApi.protectFile(fileLocation, true);
  }
}
