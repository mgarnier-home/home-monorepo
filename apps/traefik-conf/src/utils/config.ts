import dotenv from 'dotenv';
//the config is coming from a file name config.json
import Fs from 'node:fs';
import Path from 'node:path';
import { fileURLToPath } from 'node:url';

import { Config } from './interfaces.js';

dotenv.config();

console.log(process.env);

export const __filename = fileURLToPath(import.meta.url);

export const __dirname = Path.dirname(__filename);

const configFilePath = process.env.CONFIG_FILE || Path.resolve(__dirname, '../../config.json');

const loadConfigFromFile = (): Config => {
  const config = Fs.readFileSync(configFilePath, 'utf-8');

  return JSON.parse(config) as Config;
};

const loadConfigFromEnv = (): Config => {
  const traefikConfDirectory = process.env.TRAEFIK_CONF_DIR || '/traefik/dynamic';
  const composeDirectory = process.env.COMPOSE_DIR || '/compose';
  const saveDataFile = process.env.SAVE_DATA_FILE || '/data/save.json';

  const config: Config = {
    traefikConfDirectory: traefikConfDirectory.startsWith('/')
      ? traefikConfDirectory
      : Path.resolve(__dirname, '../../', traefikConfDirectory),
    composeDirectory: composeDirectory.startsWith('/')
      ? composeDirectory
      : Path.resolve(__dirname, '../../', composeDirectory),
    saveDataFile: saveDataFile.startsWith('/') ? saveDataFile : Path.resolve(__dirname, '../../', saveDataFile),
    redirectionName: process.env.REDIRECTION_NAME || '',
    serverPort: Number(process.env.SERVER_PORT) || 3000,
    stacksToIgnore: process.env.STACKS_TO_IGNORE?.split(',').map((stack) => stack.toLowerCase()) || [],
    fbxAppToken: process.env.FBX_APP_TOKEN || '',
    fbxAppId: process.env.FBX_APP_ID || '',
    fbxApiDomain: process.env.FBX_API_DOMAIN || '',
    fbxHttpsPort: Number(process.env.FBX_HTTPS_PORT) || 0,
    fbxApiBaseUrl: process.env.FBX_API_BASE_URL || '',
    fbxApiVersion: process.env.FBX_API_VERSION || '',
  };

  return config;
};

export const config = (process.env.CONFIG_FILE ? loadConfigFromFile() : loadConfigFromEnv()) as Config;
