import dotenv from 'dotenv';

import { Config } from './interfaces';

dotenv.config();

const loadConfigFromEnv = (): Config => {
  const config: Config = {
    ntfyProtocol: process.env.NTFY_PROTOCOL || 'http',
    ntfyTopic: process.env.NTFY_TOPIC || '',
    ntfyServer: process.env.NTFY_SERVER || '',
  };

  return config;
};

export const config = loadConfigFromEnv();
