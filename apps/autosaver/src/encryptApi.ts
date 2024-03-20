import { spawn } from 'child_process';
import fs from 'fs';
import { logger } from 'logger';
import path from 'path';

export namespace EncryptApi {
  export const encryptFile = (filePath: string, password: string) => {
    return new Promise<string>((resolve, reject) => {
      const encryptedArchivePath = `${filePath}.gpg`;

      if (fs.existsSync(encryptedArchivePath)) {
        logger.info(`Encrypted archive ${encryptedArchivePath} already exists, skipping`);

        resolve(encryptedArchivePath);

        return;
      }

      logger.debug(`running command: gpg --batch --passphrase ######## --symmetric ${filePath}`);

      const gpg = spawn('gpg', ['--batch', '--passphrase', password, '--symmetric', filePath], {
        cwd: path.dirname(filePath),
      });

      gpg.stdout.on('data', (data) => {
        logger.debug(`gpg stdout: ${data}`);
      });

      const errorBuffers: Buffer[] = [];

      gpg.stderr.on('data', (data: Buffer) => {
        const stringError = data.toString();

        if (
          stringError.includes(`gpg: directory '/root/.gnupg' created`) ||
          stringError.includes(`gpg: keybox '/root/.gnupg/pubring.kbx' created`)
        ) {
          logger.info('ignored error : ', stringError);
        } else {
          errorBuffers.push(data);

          gpg.kill();
        }
      });

      gpg.on('close', (code) => {
        if (code !== 0) {
          const error = Buffer.concat(errorBuffers).toString();

          reject({ message: `gpg process exited with code ${code}`, code, error });
        } else {
          resolve(encryptedArchivePath);
        }
      });
    });
  };
}
