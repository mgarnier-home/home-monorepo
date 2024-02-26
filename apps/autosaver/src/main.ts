import { logger } from 'logger';

import { CifsApi } from './cifs/cifsApi';
import { config } from './utils/config';
import { DirectoryType } from './utils/types';

logger.setAppName('autosaver');

const run = async () => {
  console.log('test44');

  console.log(config.backupConfig);

  if (config.backupConfig.rsync && config.backupConfig.rsync.type === DirectoryType.cifs) {
    await CifsApi.mountCifsFolder(
      config.backupConfig.rsync.ip,
      config.backupConfig.rsync.user,
      config.backupConfig.rsync.password,
      config.backupConfig.rsync.hostPath,
      config.backupConfig.rsync.mountPath,
      config.backupConfig.rsync.port
    );
  }

  for (const directory of config.backupConfig.directories) {
    if (directory.type === DirectoryType.cifs) {
      // await CifsApi.unmountSmbFolder(directory.mountPath);

      await CifsApi.mountCifsFolder(
        directory.ip,
        directory.user,
        directory.password,
        directory.hostPath,
        directory.mountPath,
        directory.port
      );
    } else {
      console.log(directory);
    }
  }
};
run();
