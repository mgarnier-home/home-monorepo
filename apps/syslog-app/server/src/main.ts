import type { DockerMessage, SyslogMessage } from '@shared/interfaces';

import fs from 'fs';
import path from 'path';

import { config } from './config.js';
import { SyslogServer } from './syslog.js';

const main = async () => {
  console.log('Starting syslog server');
  console.log(config);

  const syslogServer = new SyslogServer();

  syslogServer.start(config.syslogPort);
};

main();
