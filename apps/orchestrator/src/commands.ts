import { execSync } from 'child_process';
import fs from 'fs';
import path from 'path';

import { config } from './utils/config';
import { Host } from './utils/schemas';
import { getStackHostPath } from './utils/utils';

type ExecError = {
  status: number;
  signal: string;
  output: [Buffer | null, Buffer | null, Buffer | null];
};

type ExecResponse = {
  success: boolean;
  output: string;
};

const execCommand = (command: string): ExecResponse => {
  try {
    console.log(`executing command: ${command}`);

    const result = execSync(command, { stdio: 'pipe' });
    console.log(`command output: ${result.toString()}`);

    return { success: true, output: result.toString() };
  } catch (_error) {
    const error = _error as ExecError;

    console.error(`error executing command: ${command}`);
    console.error(`status: ${error.status}`);
    console.error(`signal: ${error.signal}`);
    console.error(`stdout: ${error.output[1]}`);
    console.error(`stderr: ${error.output[2]}`);

    return { success: false, output: error.output[2]?.toString() ?? '' };
  }
};

const getAdditionalComposeFiles = () => {
  // list all .yml files directly under config.composeFolderPath
  const files = fs.readdirSync(config.composeFolderPath);

  return files.filter((file) => file.endsWith('.yml')).map((file) => path.join(config.composeFolderPath, file));
};

const hasDockerContext = (host: Host) => {
  const command = `docker context ls | grep ${host.name}`;

  const result = execCommand(command);

  if (result.success) {
    return result.output.includes(host.name) && result.output.includes(`ssh://${host.username}@${host.ip}`);
  }

  return false;
};

const createContext = (host: Host) => {
  const command = `docker context create ssh ${host.name} --docker "host=ssh://${host.username}@${host.ip}"`;

  const result = execCommand(command);

  if (!result.success) {
    console.error(result.output);
  }

  return result.success;
};

const useContext = (host: Host) => {
  const command = `docker context use ${host.name}`;

  const result = execCommand(command);

  if (!result.success) {
    console.error(result.output);
  }

  return result.success;
};

const execDockerCommandOnHost = (host: Host, command: string) => {
  const contextPresent = hasDockerContext(host);

  if (!contextPresent) {
    createContext(host);
  }

  useContext(host);

  return execCommand(command);
};

const getDockerComposeCommand = (stack: string, host: Host, command: string) => {
  const additionalComposeFiles = getAdditionalComposeFiles();

  return `docker compose ${
    //
    config.composeEnvFilesPaths.map((file) => `--env-file ${file}`).join(' ')
  } -f ${
    //
    getStackHostPath(stack, host.name)
  } ${
    //
    additionalComposeFiles.map((file) => `-f ${file}`).join(' ')
  } -p ${
    //
    stack
  } ${
    //
    command
  }`;
};

const execDockerCommandListOnHost = (host: Host, commands: string[]) => {
  const contextPresent = hasDockerContext(host);

  if (!contextPresent) {
    createContext(host);
  }

  useContext(host);

  return commands.map((command) => execCommand(command));
};

export const commands = {
  up: async (stack: string, host: Host) => {
    // console.log(composeCommand);

    console.log(execDockerCommandOnHost(host, getDockerComposeCommand(stack, host, 'up -d')));
  },
  down: async (stack: string, host: Host) => {
    console.log(execDockerCommandOnHost(host, getDockerComposeCommand(stack, host, 'up -d')));
  },
  pull: async (stack: string, host: Host) => {
    console.log(
      execDockerCommandListOnHost(host, [
        getDockerComposeCommand(stack, host, 'down'),
        getDockerComposeCommand(stack, host, 'pull'),
        getDockerComposeCommand(stack, host, 'up -d'),
      ])
    );
  },
};
