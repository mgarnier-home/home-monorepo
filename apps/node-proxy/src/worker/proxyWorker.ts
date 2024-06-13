import { logger } from 'logger';
import { threadId } from 'worker_threads';

import { ProxyWorker } from '../classes/proxyWorker.class.js';

logger.setAppName('node-proxy:worker-' + threadId);

const proxyWorker = new ProxyWorker();

logger.info('Proxy worker started');
