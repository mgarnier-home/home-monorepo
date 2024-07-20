import { threadId } from 'worker_threads';

import { logger } from '@libs/logger';

import { ProxyWorker } from '../classes/proxyWorker.class';

logger.setAppName('node-proxy:worker-' + threadId);

const proxyWorker = new ProxyWorker();

logger.info('Proxy worker started');
