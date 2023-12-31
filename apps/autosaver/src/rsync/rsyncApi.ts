import { spawn } from 'child_process';
import fs from 'fs';

export class RsyncApi {
  static async rsyncFolder(source: string, destination: string) {
    let firstTime = false;

    if (!fs.existsSync(destination)) {
      await fs.promises.mkdir(destination, { recursive: true });

      firstTime = true;
    }

    return new Promise<void>((resolve, reject) => {
      const rsyncSpawn = spawn(
        'rsync',
        ['-a', '--info=progress2', '--inplace', '--no-i-r', '--delete', source, destination],
        {}
      );

      console.log(`Rsync ${source} to ${destination}`);

      const dataRegex = /\s*(\d+(?:,\d{3})*)\s+(\d+)%\s+([\d.]+[KMGTPEZY]B\/s)/;
      let lastPercent = 0;

      rsyncSpawn.stdout.on('data', (data: Buffer) => {
        const dataString = data.toString();

        const matches = dataString.match(dataRegex);

        if (matches) {
          const [, transferred, percent, speed] = matches as [string, string, string, string];

          if (parseInt(percent, 10) > lastPercent) {
            const transferredInt = parseInt(transferred.replaceAll(',', ''), 10);

            lastPercent = parseInt(percent, 10);

            console.log(`Rsync  ${percent}% - ${Math.ceil(transferredInt / 1024 / 1024)}MB ${speed}`);
          }
        }
      });

      rsyncSpawn.stderr.on('data', (data) => {
        console.log(`rsync stderr: ${data}`);
      });

      rsyncSpawn.on('close', (code) => {
        console.log(`rsync process exited with code ${code}`);

        if (code !== 0) {
          reject(new Error(`rsync process exited with code ${code}`));
        } else {
          resolve();
        }
      });
    });
  }
}
