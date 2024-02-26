import { spawnSync } from 'child_process';
import fs from 'fs';
import { logger } from 'logger';

export class CifsApi {
  static async mountCifsFolder(
    host: string,
    user: string,
    password: string,
    hostPath: string,
    mountPath: string,
    port: number = 445
  ): Promise<void> {
    const options = `rw,username=${user},password=${password},iocharset=utf8,uid=1000,sec=ntlmv2,file_mode=0777,dir_mode=0777,port=${port}`;

    logger.info(`Mounting smb folder from ${host}${hostPath} to ${mountPath}`);

    if (fs.existsSync(mountPath)) {
      try {
        await CifsApi.unmountSmbFolder(mountPath);
      } catch (error) {
        logger.error(`Error unmounting folder ${mountPath}: ${error}`);
      }
    }
    await fs.promises.mkdir(mountPath, { recursive: true });

    logger.debug(`running command: mount -t cifs -o ${options} //${host}${hostPath} ${mountPath}`);

    const smbMountSpawn = spawnSync('mount', ['-t', 'cifs', '-o', options, `//${host}${hostPath}`, mountPath], {});

    if (smbMountSpawn.status !== 0) {
      throw new Error(`smbMountSpawn process exited with code ${smbMountSpawn.status}`, {
        cause: smbMountSpawn.error || smbMountSpawn.stderr.toString(),
      });
    }

    logger.info(`Mounted smb folder`);
  }

  static async unmountSmbFolder(mountPath: string): Promise<void> {
    logger.info(`Unmounting smb folder from ${mountPath}`);

    const smbUnmountSpawn = spawnSync('umount', [mountPath], {});

    if (smbUnmountSpawn.status !== 0) {
      throw new Error(`smbMountSpawn process exited with code ${smbUnmountSpawn.status}`, {
        cause: smbUnmountSpawn.error || smbUnmountSpawn.stderr.toString(),
      });
    }

    await fs.promises.rm(mountPath, { recursive: true });

    logger.info(`Unmounted smb folder`);
  }
}
