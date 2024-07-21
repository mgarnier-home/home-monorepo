import { getEnvVariable } from '@libs/env-config';
import { logger } from '@libs/logger';

import { createExpressApp } from './utils/express.utils';
import { setupConfigListenner } from './utils/host.utils';

logger.setAppName('node-proxy');

const main = async () => {
  setupConfigListenner();

  createExpressApp(getEnvVariable('SERVER_PORT', false, 3000));
};

main();
