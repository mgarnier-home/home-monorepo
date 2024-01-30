import type { DockerMessage, SyslogMessage } from '@shared/interfaces';

import fs from 'fs';
import path from 'path';

import { config } from './config';
import { SyslogServer } from './syslog';
import { log } from './utils';

const main = async () => {
  console.log('\x1b[36m%s\x1b[0m', 'I am cyan'); //cyan
  console.log('\x1b[33m%s\x1b[0m', 'I am yellow'); //yellow

  log('Starting syslog server');

  log(config);

  const syslogServer = new SyslogServer();

  syslogServer.start(config.syslogPort);
};

main();
