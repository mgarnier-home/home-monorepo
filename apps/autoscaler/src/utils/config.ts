import { readFileSync } from 'fs';
import { resolve } from 'path';
import { cwd } from 'process';
import * as YAML from 'yaml';

import { getEnvVariable } from '@libs/env-config';

import { AutoscalerConfig, Config } from './interfaces';

const loadConfig = (): Config => {
  const autoscalerConfigPath = getEnvVariable<string>('AUTOSCALER_CONFIG_PATH', false, './config.yml');
  const fullAutoscalerConfigPath = autoscalerConfigPath.startsWith('/')
    ? autoscalerConfigPath
    : resolve(cwd(), autoscalerConfigPath);

  const config: Config = {
    serverPort: getEnvVariable('SERVER_PORT', false, 3000),
    nodeEnv: getEnvVariable('NODE_ENV', false, 'development'),
    webhookSecret: getEnvVariable('WEBHOOK_SECRET', true),
    smeeUrl: getEnvVariable('SMEE_URL', false, ''),
    runnerImage: getEnvVariable('RUNNER_IMAGE', false, 'myoung34/github-runner:latest'),
    runnerRepoUrl: getEnvVariable('RUNNER_REPO_URL', true),
    runnerAccessToken: getEnvVariable('RUNNER_ACCESS_TOKEN', true),
    autoscalerConfigPath: fullAutoscalerConfigPath,
  };

  return config;
};

export const config = loadConfig();

export const getAutoscalerConfig = (): AutoscalerConfig => {
  const autoscalerConfigYml = readFileSync(config.autoscalerConfigPath, 'utf-8');
  const autoscalerConfig: AutoscalerConfig = YAML.parse(autoscalerConfigYml, { merge: true }) as AutoscalerConfig;

  return autoscalerConfig;
};
