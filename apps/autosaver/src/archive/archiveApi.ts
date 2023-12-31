import { spawn } from 'child_process';
import fs from 'fs';
import path from 'path';

import { OsUtils } from '../utils/osUtils.js';

export abstract class ArchiveApi {
  public async archiveFolder(
    folderToBackupPath: string,
    archivePassword?: string
  ): Promise<{ archivePath: string; nbFilesArchived: number; archiveSize: number }> {
    const nbFiles = await OsUtils.osCountFiles(folderToBackupPath);

    let archivePath = await this.archive(folderToBackupPath, nbFiles);

    console.log(`Archive created : ${archivePath}`);

    if (archivePassword) {
      console.log(`Encrypting archive : ${archivePath}...`);

      archivePath = await this.encryptArchive(archivePath, archivePassword);

      console.log(`Archive encrypted : ${archivePath}`);

      await OsUtils.rmFiles([archivePath.replace('.gpg', '')]);
    }

    const archiveSize = (await fs.promises.stat(archivePath)).size;

    return { archivePath, nbFilesArchived: nbFiles, archiveSize };
  }

  protected abstract archive(folderToBackupPath: string, totalNbFiles: number): Promise<string>;

  protected async encryptArchive(archivePath: string, archivePassword: string) {
    return new Promise<string>((resolve, reject) => {
      const encryptedArchivePath = `${archivePath}.gpg`;

      if (fs.existsSync(encryptedArchivePath)) {
        console.log(`Encrypted archive ${encryptedArchivePath} already exists, skipping`);

        resolve(encryptedArchivePath);

        return;
      }

      const gpg = spawn('gpg', ['--batch', '--passphrase', archivePassword, '--symmetric', archivePath], {
        cwd: path.dirname(archivePath),
      });

      gpg.stdout.on('data', (data) => {
        console.log(`gpg stdout: ${data}`);
      });

      gpg.stderr.on('data', (data) => {
        console.error(`gpg stderr: ${data}`);

        gpg.kill();

        reject(new Error(`gpg process exited with error`));
      });

      gpg.on('close', (code) => {
        if (code !== 0) {
          reject(new Error(`gpg process exited with code ${code}`));

          return;
        }

        resolve(encryptedArchivePath);
      });
    });
  }
}
