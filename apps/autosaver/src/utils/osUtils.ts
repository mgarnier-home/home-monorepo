import { exec } from 'child_process';
import fs from 'fs';
import path from 'path';

import { config } from './config.js';

export class OsUtils {
  static async rmFiles(filePaths: string[]) {
    if (config.deleteFiles === false) {
      console.log('not deleting files : ', filePaths);
      return;
    }

    try {
      for (const path of filePaths) {
        await fs.promises.unlink(path);

        console.log(`File ${path} deleted`);
      }
    } catch (error) {
      console.error(`Error while deleting [${filePaths.join(', ')}] : `, error);
    }
  }

  static async getFilesInFolder(
    folderPath: string,
    foldersToIgnore: string[] = []
  ): Promise<{ path: string; stat: fs.Stats }[]> {
    const files: { path: string; stat: fs.Stats }[] = [];

    const filesInFolder = await fs.promises.readdir(folderPath);
    const folderName = path.basename(folderPath);

    if (foldersToIgnore.includes(folderName)) {
      return files;
    }

    for (const fileInFolder of filesInFolder) {
      const filePath = path.join(folderPath, fileInFolder);
      const stat = await fs.promises.lstat(filePath);

      if (stat.isDirectory()) {
        const filesInSubFolder = await OsUtils.getFilesInFolder(filePath, foldersToIgnore);

        files.push(...filesInSubFolder);
      }

      if (stat.isFile()) {
        files.push({ path: filePath, stat });
      }
    }

    return files;
  }

  static async osCountFiles(dir: string) {
    return new Promise<number>((resolve, reject) => {
      exec(`find "${dir}" -type f | wc -l`, (error, stdout, stderr) => {
        if (error) {
          reject(error);
          return;
        }

        const filesNb = parseInt(stdout.trim(), 10);

        resolve(filesNb);
      });
    });
  }

  static async countFiles(dir: string) {
    const files = await OsUtils.getFilesInFolder(dir);

    return files.length;
  }
}
