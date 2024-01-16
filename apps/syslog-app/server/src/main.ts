import type { DockerMessage, SyslogMessage } from '@shared/interfaces';

import fs from 'fs';
import path from 'path';

import { config } from './config';
import { SyslogServer } from './syslog';

const main = async () => {
  console.log('Starting syslog server');
  console.log(config);
  console.log(new Date().toLocaleString());

  const syslogServer = new SyslogServer();

  syslogServer.start(config.syslogPort);
};

main();
