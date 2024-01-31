import { logger } from 'logger';

import { config } from './config';
import { SyslogServer } from './syslog';

logger.setAppName('syslog-server');

const main = async () => {
  logger.info('Starting syslog server');

  logger.info(config);

  const syslogServer = new SyslogServer();

  syslogServer.start(config.syslogPort);
};

main();
