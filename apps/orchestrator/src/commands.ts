import { exec } from 'child_process';
import fs from 'fs';
import path from 'path';

import { config } from './utils/config';
import { Host } from './utils/schemas';
import { getStackHostPath } from './utils/utils';

const execCommand = (command: string): Promise<string> => {
  return new Promise<string>((resolve, reject) => {
    exec(command, (error, stdout, stderr) => {
      if (error) {
        console.error(`exec error ${error}`);
        return reject(error);
      }

      if (stderr) {
        console.error(`stderr ${stderr}`);
        return reject(new Error(stderr));
      }

      resolve(stdout);
    });
  });
};

const getAdditionalComposeFiles = () => {
  // list all .yml files directly under config.composeFolderPath
  const files = fs.readdirSync(config.composeFolderPath);

  return files.filter((file) => file.endsWith('.yml')).map((file) => path.join(config.composeFolderPath, file));
};

const hasDockerContext = async (host: Host) => {
  const command = `docker context ls | grep ${host.name}e`;

  try {
    const result = await execCommand(command);

    return result.includes(host.name) && result.includes(`ssh://${host.username}@${host.ip}`);
  } catch (error) {
    return false;
  }
};

export const commands = {
  up: async (stack: string, host: Host) => {
    console.log('up', stack, host, getAdditionalComposeFiles());

    const additionalComposeFiles = getAdditionalComposeFiles();

    const composeCommand = `docker compose -f ${
      //
      getStackHostPath(stack, host.name)
    } ${
      //
      additionalComposeFiles.map((file) => `-f ${file}`).join(' ')
    } -p ${
      //
      stack
    } up -d`;

    console.log(composeCommand);

    console.log(await hasDockerContext(host));
    // execute docker compose up
  },
  down: async (stack: string, host: Host) => {
    console.log('down', stack, host);
  },
  pull: async (stack: string, host: Host) => {
    console.log('pull', stack, host);
  },
};
