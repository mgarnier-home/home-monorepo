import dotenv from 'dotenv';
//the config is coming from a file name config.json
import Fs from 'node:fs';
import Path from 'node:path';

import { Config } from './interfaces.js';

dotenv.config();

const configFilePath = process.env.CONFIG_FILE || Path.resolve(__dirname, '../../config.json');

const loadConfigFromFile = (): Config => {
  const config = Fs.readFileSync(configFilePath, 'utf-8');

  return JSON.parse(config) as Config;
};

const loadConfigFromEnv = (): Config => {
  const config: Config = {
    dataFilePath: Path.resolve('./', process.env.DATA_FILE_PATH || '/data/data.json'),
    serverPort: Number(process.env.SERVER_PORT) || 3000,
    stacksToIgnore: process.env.STACKS_TO_IGNORE?.split(',').map((stack) => stack.toLowerCase()) || [],
  };

  return config;
};

export const config = (process.env.CONFIG_FILE ? loadConfigFromFile() : loadConfigFromEnv()) as Config;
