import { logger } from '@libs/logger';
import { Webhooks } from '@octokit/webhooks';

import { config } from './utils/config';

const webhooks = new Webhooks({
  secret: config.webhookSecret,
});

webhooks.on('workflow_job', (event) => {
  logger.debug('===============================');
  logger.debug('workflow job event');
  logger.debug('workflow_job.name : ', event.payload.workflow_job.name);
  logger.debug('workflow_job.status : ', event.payload.workflow_job.status);
  logger.debug('workflow_job.conclusion : ', event.payload.workflow_job.conclusion);
  logger.debug('workflow_job.labels : ', event.payload.workflow_job.labels);
  logger.debug('===============================');

  if (
    event.payload.workflow_job.labels.includes('linux/arm64') ||
    event.payload.workflow_job.labels.includes('linux/amd64')
  ) {
    if (event.payload.workflow_job.status === 'queued') {
      logger.info('Should deploy a new runner');
    } else if (event.payload.workflow_job.status === 'completed') {
      logger.info('Should remove the runner');
    }
  }
});

export { webhooks };
