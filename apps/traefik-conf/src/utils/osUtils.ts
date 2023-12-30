import fs from 'fs';
import path from 'path';

export const listFiles = async (dir: string, extension: string = '', recursive: boolean = false): Promise<string[]> => {
  const filesInFolder = await fs.promises.readdir(dir);
  const files: string[] = [];

  for (const file of filesInFolder) {
    const filePath = path.join(dir, file);
    const stat = await fs.promises.lstat(filePath);

    if (stat.isDirectory()) {
      const subFiles = await listFiles(filePath, extension, recursive);

      files.push(...subFiles);
    }

    if (stat.isFile() && file.endsWith(extension)) {
      files.push(filePath);
    }
  }

  return files;
};

export const readFiles = async (filePaths: string[]): Promise<string[]> => {
  const fileDatas = await Promise.all(filePaths.map((path) => fs.promises.readFile(path, 'utf-8')));

  return fileDatas;
};
