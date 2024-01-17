import type { DockerMessage, SyslogMessage } from '@shared/interfaces';

import fs from 'fs';
import path from 'path';

import { config } from './config';
import { SyslogServer } from './syslog';
import { log } from './utils';

const main = async () => {
  log('Starting syslog server');
  log(config);

  const syslogServer = new SyslogServer();

  syslogServer.start(config.syslogPort);
};

main();
