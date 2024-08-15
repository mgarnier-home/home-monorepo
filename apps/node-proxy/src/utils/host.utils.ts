import fs from 'fs';
import path, { resolve } from 'path';
import { cwd } from 'process';
import * as YAML from 'yaml';

import { getEnvVariable } from '@libs/env-config';
import { logger } from '@libs/logger';

import { Host } from '../classes/host.class';
import { HostConfig } from './interfaces';

const hosts: Host[] = [];
const configFilePath = getEnvVariable('CONFIG_FILE', false, '../config.yml');
const fullConfigPath = configFilePath.startsWith('/') ? configFilePath : resolve(cwd(), configFilePath);

const lastConfig: HostConfig[] = [];

const loadConfig = async (): Promise<HostConfig[]> => {
  try {
    if (fs.existsSync(fullConfigPath)) {
      const dataStr = await fs.promises.readFile(fullConfigPath, 'utf-8');

      if (dataStr !== '') {
        if (fullConfigPath.endsWith('on')) {
          return JSON.parse(dataStr) as HostConfig[];
        }
        return YAML.parse(dataStr, { merge: true }) as HostConfig[];
      }
    } else {
      logger.warn('Config file not found : ', fullConfigPath);
    }
  } catch (err) {
    logger.error('Error loading config file : ', err);
  }

  return [];
};

const saveConfig = async (data: HostConfig[]) => {
  const stringData = fullConfigPath.endsWith('json') ? JSON.stringify(data, null, 4) : YAML.stringify(data);

  await fs.promises.writeFile(fullConfigPath, stringData, 'utf-8');
};

const hostConfigUpdated = () => {
  const data = hosts.map((host) => host.config);

  saveConfig(data);
};

export const getHost = (host: string): Host | undefined => {
  logger.debug('Getting host : ', host);
  return hosts.find((h) => h.config.name.toLowerCase() === host.toLowerCase());
};

export const setupConfigListenner = () => {
  logger.debug('Loading config from : ', fullConfigPath);

  const configFileChanged = async (configs: HostConfig[]) => {
    await disposeHosts();

    logger.info('Hosts loaded : ', configs);

    hosts.push(...configs.map((config) => new Host(config, hostConfigUpdated)));
  };

  const checkConfig = async () => {
    const config = await loadConfig();

    if (JSON.stringify(config) !== JSON.stringify(lastConfig)) {
      lastConfig.splice(0, lastConfig.length);
      lastConfig.push(...config);

      logger.info('Config file changed, triggering callback');

      configFileChanged(config);
    }
  };

  setInterval(checkConfig, 30 * 1000);

  checkConfig();
};

export const disposeHosts = async () => {
  for (const host of hosts) {
    await host.dispose();

    logger.info('Host disposed : ', host.config.name);
  }

  hosts.splice(0, hosts.length);

  logger.info('All hosts disposed');
};
