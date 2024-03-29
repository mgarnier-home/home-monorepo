import express from 'express';
import fs from 'fs';
import path from 'path';

import { commands } from './commands';
import { config, getStack } from './utils/config';
import { Stack } from './utils/schemas';
import { getStackHostPath } from './utils/utils';

console.log(config);

type Request = express.Request & { stackConfig?: Stack };

const checkStackForHost = (stack: string, hostName: string) => {
  //
  // if a file config.composeFolderPath/${stack}/${host}.${stack}.yml exists then return true

  return fs.existsSync(getStackHostPath(stack, hostName));
};

const main = async () => {
  const app = express();

  app.use('/', (req: Request, res, next) => {
    console.log('middleware');
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

  (Object.keys(commands) as (keyof typeof commands)[]).forEach((command) => {
    app.get(`/${command}/:stack/:host`, async (req: Request, res) => {
      const { stack, host: hostName } = req.params;
      console.log('stack', stack);
      console.log('hostName', hostName);
      const stackConfig: Stack = req.stackConfig!;

      const host = stackConfig.hosts.find((host) => host.name === hostName)!;

      if (!stack || !hostName) {
        return res.status(400).send('Invalid request');
      }

      if (stack !== 'all' && !stackConfig.stacks.includes(stack)) {
        return res.status(400).send('Invalid stack');
      }

      if (hostName !== 'all' && !stackConfig.hosts.includes(host)) {
        return res.status(400).send('Invalid host');
      }

      if (stack !== 'all' && hostName !== 'all' && !checkStackForHost(stack, hostName)) {
        return res.status(400).send(`Stack ${stack} not found for host ${hostName}`);
      }

      await commands[command](stack, host);

      res.send(`Running ${command} on ${stack} for host ${hostName}`);
    });
  });
};

main();
