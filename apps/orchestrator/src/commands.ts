import { execSync } from 'child_process';
import fs from 'fs';
import path from 'path';

import { config } from './utils/config';
import { Host } from './utils/schemas';
import { checkStackForHost, getAdditionalComposeFiles, getStackHostPath } from './utils/utils';

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
    // console.log(`executing command: ${command}`);

    const result = execSync(command, { stdio: 'pipe' });
    // console.log(`command output: ${result.toString()}`);

    return { success: true, output: result.toString() };
  } catch (_error) {
    const error = _error as ExecError;

    // console.error(`error executing command: ${command}`);
    // console.error(`status: ${error.status}`);
    // console.error(`signal: ${error.signal}`);
    // console.error(`stdout: ${error.output[1]}`);
    // console.error(`stderr: ${error.output[2]}`);

    return { success: false, output: error.output[2]?.toString() ?? '' };
  }
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
  const command = `docker context create --docker "host=ssh://${host.username}@${host.ip}" ${host.name}`;

  const result = execCommand(command);

  if (!result.success) {
    console.error(result.output);
  }

  return result.success;
};

const useContext = (hostName: string) => {
  const command = `docker context use ${hostName}`;

  const result = execCommand(command);

  if (!result.success) {
    console.error(result.output);
  }

  return result.success;
};

const getDockerComposeCommand = (stack: string, host: Host, command: string) => {
  const additionalComposeFiles = getAdditionalComposeFiles();

  return `docker compose ${
    //
    config.composeEnvFilesPaths.map((file) => `--env-file ${file.trim()}`).join(' ')
  } -f ${
    //
    getStackHostPath(stack, host.name)
  } ${
    //
    additionalComposeFiles.map((file) => `-f ${file.trim()}`).join(' ')
  } -p ${
    //
    stack
  } ${
    //
    command
  }`;
};

const runCommands = (
  stacks: string[],
  hosts: Host[],
  commands: string[]
): {
  name: string;
  stacks: { name: string; commands: { full: string; base: string; success: boolean; output: string }[] }[];
}[] => {
  return hosts.map((host) => {
    const contextPresent = hasDockerContext(host);

    if (!contextPresent) {
      createContext(host);
    }

    useContext(host.name);

    const stacksResults = stacks
      .filter((stack) => checkStackForHost(stack, host.name))
      .map((stack) => ({
        name: stack,
        commands: commands.map((command) => {
          console.log(`Host : ${host.name} Stack : ${stack} Command : ${command}`);
          const fullCommand = getDockerComposeCommand(stack, host, command);
          const commandResult = execCommand(fullCommand);

          return {
            full: fullCommand,
            base: command,
            success: commandResult.success,
            output: commandResult.output,
          };
        }),
      }));

    useContext('default');

    return {
      name: host.name,
      stacks: stacksResults,
    };
  });
};

export const commands = {
  up: async (stacks: string[], hosts: Host[]) => {
    console.log('=====================================================================');
    console.log('UP');

    console.log(JSON.stringify(runCommands(stacks, hosts, ['up -d']), null, 2));

    console.log('UP COMPLETE');
    console.log('=====================================================================');
  },
  down: async (stacks: string[], hosts: Host[]) => {
    console.log('=====================================================================');
    console.log('DOWN');

    console.log(JSON.stringify(runCommands(stacks, hosts, ['down']), null, 2));

    console.log('DOWN COMPLETE');
    console.log('=====================================================================');
  },
  pull: async (stacks: string[], hosts: Host[]) => {
    console.log('=====================================================================');
    console.log('PULL');

    console.log(JSON.stringify(runCommands(stacks, hosts, ['down', 'pull', 'up -d']), null, 2));

    console.log('PULL COMPLETE');
    console.log('=====================================================================');
  },
};
