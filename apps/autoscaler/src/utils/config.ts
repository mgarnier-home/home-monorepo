import { getEnvVariable } from '@libs/env-config';

import { Config } from './interfaces';

const loadConfig = (): Config => {
  const config: Config = {
    serverPort: getEnvVariable('SERVER_PORT', false, 3000),
    nodeEnv: getEnvVariable('NODE_ENV', false, 'development'),
    webhookSecret: getEnvVariable('WEBHOOK_SECRET', true),
    arm64Host: getEnvVariable('ARM64_HOST', true),
    amd64Host: getEnvVariable('AMD64_HOST', true),
    smeeUrl: getEnvVariable('SMEE_URL', false, ''),
  };

  return config;
};

export const config = loadConfig();
