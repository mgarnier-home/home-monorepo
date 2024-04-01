import express from 'express';
import fs from 'fs';
import path from 'path';

import { commands } from './commands';
import { config, getStack } from './utils/config';
import { Stack } from './utils/schemas';
import { checkStackForHost } from './utils/utils';

console.log(config);

type Request = express.Request & { stackConfig?: Stack };

const main = async () => {
  const app = express();

  app.use('/', (req: Request, res, next) => {
    try {
      const stackConfig = getStack();

      stackConfig.hosts = stackConfig.hosts.sort((a, b) => a.name.localeCompare(b.name));
      stackConfig.stacks = stackConfig.stacks.sort();

      req.stackConfig = stackConfig;

      // console.log('stack config', stackConfig);
      next();
    } catch (error) {
      console.error(`error getting stack config: ${error}`);
      res.status(500).send('Internal server error');
    }
  });

  app.get('/', (req: Request, res) => {
    res.send(req.stackConfig);
  });

  app.listen(config.serverPort, () => {
    console.log(`Server listening on port ${config.serverPort}`);
  });

  for (const command of Object.keys(commands) as (keyof typeof commands)[]) {
    app.get(`/${command}/:stackName/:hostName`, async (req: Request, res) => {
      const { stackName, hostName } = req.params;
      const stackConfig: Stack = req.stackConfig!;

      const host = stackConfig.hosts.find((host) => host.name === hostName)!;

      if (!stackName || !hostName) {
        return res.status(400).send('Invalid request');
      }

      if (stackName !== 'all' && !stackConfig.stacks.includes(stackName)) {
        return res.status(400).send('Invalid stack');
      }

      if (hostName !== 'all' && !stackConfig.hosts.includes(host)) {
        return res.status(400).send('Invalid host');
      }

      if (stackName !== 'all' && hostName !== 'all' && !checkStackForHost(stackName, hostName)) {
        return res.status(400).send(`Stack ${stackName} not found for host ${hostName}`);
      }

      res.send(`Running ${stackName} ${command} on host ${hostName}`);

      const stacksToRun = stackName === 'all' ? stackConfig.stacks : [stackName];

      const hostsToRun = hostName === 'all' ? stackConfig.hosts : [host];

      await commands[command](stacksToRun, hostsToRun);
    });
  }
};

main();
