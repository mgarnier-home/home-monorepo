import { logger } from 'logger';

import { createExpressApp } from './utils/express.utils';
import { setupHosts } from './utils/host.utils';

logger.setAppName('node-proxy');

const main = async () => {
  await setupHosts();

  createExpressApp((process.env.SERVER_PORT || 3000) as number);
};

main();
