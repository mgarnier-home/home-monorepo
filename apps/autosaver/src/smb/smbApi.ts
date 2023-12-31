import { spawn } from 'child_process';
import fs from 'fs';

import { SmbConfig } from '../utils/interfaces';

export class SmbApi {
  static async mountSmbFolders(smbConfig: SmbConfig) {
    for (const smbHost of smbConfig?.hosts || []) {
      for (const smbMount of smbHost.mounts) {
        const root = smbHost.mountRoot.startsWith('/') ? smbHost.mountRoot : `/${smbHost.mountRoot}`;
        const src = smbMount.src.startsWith('/') ? smbMount.src : `/${smbMount.src}`;
        const dest = smbMount.dest.startsWith('/') ? smbMount.dest : `/${smbMount.dest}`;

        await SmbApi.mountSmbFolder(smbHost.ip, smbHost.user, smbHost.password, `${root}${src}`, dest);

        console.log(`Successfully mounted smb`);
      }
    }
  }

  static async unmountSmbFolders(smbConfig: SmbConfig) {
    for (const smbHost of smbConfig?.hosts || []) {
      for (const smbMount of smbHost.mounts) {
        await SmbApi.unmountSmbFolder(smbMount.dest);
      }
    }
  }

  static async mountSmbFolder(
    smbHost: string,
    smbUser: string,
    smbPassword: string,
    smbFolder: string,
    destFolder: string
  ) {
    const options = `rw,username=${smbUser},password=${smbPassword},iocharset=utf8,uid=1000,sec=ntlmv2,file_mode=0777,dir_mode=0777`;

    console.log(`Mounting smb folder from ${smbHost}${smbFolder} to ${destFolder}`);

    if (fs.existsSync(destFolder)) {
      try {
        await SmbApi.unmountSmbFolder(destFolder);
      } catch (error) {
        console.log(`Error unmounting folder ${destFolder}: ${error}`);
      }
    }
    await fs.promises.mkdir(destFolder, { recursive: true });

    return new Promise<void>((resolve, reject) => {
      const smbMountSpawn = spawn('mount', ['-t', 'cifs', '-o', options, `//${smbHost}${smbFolder}`, destFolder], {});

      smbMountSpawn.stdout.on('data', (data) => {
        console.log(`smbMountSpawn stdout: ${data}`);
      });

      smbMountSpawn.stderr.on('data', (data) => {
        console.log(`smbMountSpawn stderr: ${data}`);

        smbMountSpawn.kill(1);
      });

      smbMountSpawn.on('close', (code) => {
        console.log(`smbMountSpawn process exited with code ${code}`);

        if (code !== 0) {
          reject(new Error(`smbMountSpawn process exited with code ${code}`));
        } else {
          resolve();
        }
      });
    });
  }

  static async unmountSmbFolder(destFolder: string): Promise<void> {
    console.log(`Unmounting smb folder from ${destFolder}`);

    await new Promise<void>((resolve, reject) => {
      const smbUnmountSpawn = spawn('umount', [destFolder], {});

      smbUnmountSpawn.stdout.on('data', (data) => {
        console.log(`smbUnmountSpawn stdout: ${data}`);
      });

      smbUnmountSpawn.stderr.on('data', (data) => {
        console.log(`smbUnmountSpawn stderr: ${data}`);

        smbUnmountSpawn.kill(1);
      });

      smbUnmountSpawn.on('close', (code) => {
        console.log(`smbUnmountSpawn process exited with code ${code}`);

        if (code !== 0) {
          reject(new Error(`smbUnmountSpawn process exited with code ${code}`));
        } else {
          resolve();
        }
      });
    });

    await fs.promises.rm(destFolder, { recursive: true });
  }
}
