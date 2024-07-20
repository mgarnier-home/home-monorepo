import { logger } from '@libs/logger';

import { createExpressApp } from './utils/express.utils';
import { setupConfigListenner } from './utils/host.utils';

logger.setAppName('node-proxy');

const main = async () => {
  setupConfigListenner();

  createExpressApp((process.env.SERVER_PORT || 3000) as number);
};

main();
