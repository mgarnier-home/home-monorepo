import { logger } from 'logger';

import { createExpressApp } from './utils/express.utils.js';
import { setupConfigListenner } from './utils/host.utils.js';

logger.setAppName('node-proxy');

const main = async () => {
  setupConfigListenner();

  createExpressApp((process.env.SERVER_PORT || 3000) as number);
};

main();
