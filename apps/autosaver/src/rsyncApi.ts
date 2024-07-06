import { spawn } from 'child_process';
import fs from 'fs';

import { logger } from '@libs/logger';

export namespace RsyncApi {
  export const rsyncFolder = async (source: string, destination: string) => {
    let firstTime = false;

    if (!fs.existsSync(destination)) {
      await fs.promises.mkdir(destination, { recursive: true });

      firstTime = true;
    }

    logger.info(`Rsync ${source} to ${destination}`);
    logger.debug(`running command: rsync -a --info=progress2 --inplace --no-i-r --delete ${source} ${destination}`);

    return new Promise<void>((resolve, reject) => {
      const rsyncSpawn = spawn(
        'rsync',
        ['-a', '--info=progress2', '--inplace', '--no-i-r', '--delete', source, destination],
        {}
      );

      logger.info(`Rsync ${source} to ${destination}`);

      let lastPercent = 0;

      rsyncSpawn.stdout.on('data', (data: Buffer) => {
        const dataString = data.toString();

        const matches = dataString.match(/\s*(\d+(?:,\d{3})*)\s+(\d+)%\s+([\d.]+[KMGTPEZY]B\/s)/);

        if (matches) {
          const [, transferred, percent, speed] = matches as [string, string, string, string];

          if (parseInt(percent, 10) > lastPercent) {
            const transferredInt = parseInt(transferred.replaceAll(',', ''), 10);

            lastPercent = parseInt(percent, 10);

            logger.info(`Rsync  ${percent}% - ${Math.ceil(transferredInt / 1024 / 1024)}MB ${speed}`);
          }
        }
      });

      const errorBuffers: Buffer[] = [];

      rsyncSpawn.stderr.on('data', (data: Buffer) => {
        errorBuffers.push(data);
      });

      rsyncSpawn.on('close', (code, signal) => {
        logger.debug(`rsync process exited with code ${code}, signal ${signal}`);

        if (code !== 0) {
          const error = Buffer.concat(errorBuffers).toString();

          reject({ message: `rsync process exited with code ${code}`, code, error });
        } else {
          resolve();
        }
      });
    });
  };
}
