import { logger } from '@libs/logger';
import { NtfyUtils } from '@libs/ntfy-utils';
import { Webhooks } from '@octokit/webhooks';

import { startRunner, stopRunner } from './docker';
import { config, getAutoscalerConfig } from './utils/config';

const webhooks = new Webhooks({
  secret: config.webhookSecret,
});

webhooks.on('workflow_job', async (event) => {
  logger.debug('===============================');
  logger.debug('workflow job event');
  logger.debug('workflow_job.name : ', event.payload.workflow_job.name);
  logger.debug('workflow_job.status : ', event.payload.workflow_job.status);
  logger.debug('workflow_job.conclusion : ', event.payload.workflow_job.conclusion);
  logger.debug('workflow_job.labels : ', event.payload.workflow_job.labels);
  logger.debug('workflow_job.id : ', event.payload.workflow_job.id);
  logger.debug('workflow_job.run_id : ', event.payload.workflow_job.run_id);
  logger.debug('===============================');

  const autoscalerConfig = getAutoscalerConfig();

  const targetHost = autoscalerConfig.autoscalerHosts.find((host) => {
    return event.payload.workflow_job.labels.includes(host.label);
  });

  if (!targetHost) {
    logger.error('No host found for the given label', event.payload.workflow_job.labels);
    return;
  }

  if (event.payload.workflow_job.status === 'queued') {
    try {
      logger.info(`Should add a runner to ${targetHost.label}`);
      await startRunner(targetHost, event.payload.workflow_job.id);
    } catch (error) {
      logger.error('Error while starting runner', error);
      NtfyUtils.sendNotification('Autoscaler error', `Error while starting runner\n${error}`, '');
    }
  } else if (event.payload.workflow_job.status === 'completed') {
    try {
      logger.info(`Should remove a runner from ${targetHost.label}`);
      await stopRunner(targetHost, event.payload.workflow_job.id);
    } catch (error) {
      logger.error('Error while stopping runner', error);
      NtfyUtils.sendNotification('Autoscaler error', `Error while stopping runner\n${error}`, '');
    }
  }
});

export { webhooks };
