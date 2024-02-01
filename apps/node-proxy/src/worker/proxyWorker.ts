import { logger } from 'logger';

import { ProxyWorker } from '../classes/proxyWorker.class.js';

const proxyWorker = new ProxyWorker();

logger.info('Proxy worker started');
