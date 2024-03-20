import { spawnSync } from 'child_process';
import fs from 'fs';
import { logger } from 'logger';

export namespace CifsApi {
  export const mountCifsFolder = async (
    host: string,
    user: string,
    password: string,
    hostPath: string,
    mountPath: string,
    port: number = 445
  ): Promise<void> => {
    logger.info(`Mounting smb folder from //${host}${hostPath} to ${mountPath}`);

    if (fs.existsSync(mountPath)) {
      try {
        await unmountSmbFolder(mountPath);
      } catch (error) {
        logger.error(`Error unmounting folder ${mountPath}`);
        logger.error(error);
      }
    }

    await fs.promises.mkdir(mountPath, { recursive: true });

    const options = `rw,vers=3.0,username=${user},password=${password},iocharset=utf8,uid=1000,sec=ntlmv2,file_mode=0777,dir_mode=0777,port=${port}`;

    logger.debug(`running command: mount -t cifs -o ${options} //${host}${hostPath} ${mountPath}`);

    const smbMountSpawn = spawnSync('mount', ['-t', 'cifs', '-o', options, `//${host}${hostPath}`, mountPath], {});

    if (smbMountSpawn.status !== 0) {
      throw {
        message: `smbMountSpawn process exited with code ${smbMountSpawn.status}`,
        code: smbMountSpawn.status,
        error: smbMountSpawn.error || smbMountSpawn.stderr.toString(),
      };
    }

    logger.info(`Mounted smb folder`);
  };

  export const unmountSmbFolder = async (mountPath: string): Promise<void> => {
    logger.info(`Unmounting smb folder from ${mountPath}`);

    logger.debug(`running command: umount ${mountPath}`);

    const smbUnmountSpawn = spawnSync('umount', [mountPath], {});

    if (smbUnmountSpawn.status !== 0) {
      throw {
        message: `smbMountSpawn process exited with code ${smbUnmountSpawn.status}`,
        code: smbUnmountSpawn.status,
        error: smbUnmountSpawn.error || smbUnmountSpawn.stderr.toString(),
      };
    }

    await fs.promises.rm(mountPath, { recursive: true });

    logger.info(`Unmounted smb folder`);
  };
}
