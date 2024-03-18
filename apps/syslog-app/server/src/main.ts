import express from 'express';
import { logger } from 'logger';

import { config } from './config';
import { SyslogServer } from './syslog';

logger.setAppName('syslog-server');

const main = async () => {
  logger.info('Starting syslog server');

  logger.info(config);

  const syslogServer = new SyslogServer();
  const app = express();

  syslogServer.start(config.syslogPort);

  app.get('/', (req, res) => {
    res.send('OK');
  });

  app.listen(config.serverPort, () => {
    logger.info(`Server listening on port ${config.serverPort}`);
  });
};

main();
