import EventSource from 'eventsource';
import express from 'express';

import { setVersionEndpoint } from '@libs/api-version';
import { logger } from '@libs/logger';
import { Webhooks } from '@octokit/webhooks';

import { config } from './utils/config';

logger.setAppName('autoscaler');

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

if (config.nodeEnv === 'development') {
  logger.info('Running in development mode');

  const source = new EventSource(config.smeeUrl);
  source.onmessage = async (event: any) => {
    logger.debug('Received event from smee');
    const webhookEvent = JSON.parse(event.data);
    try {
      await webhooks.verifyAndReceive({
        id: webhookEvent['x-request-id'],
        name: webhookEvent['x-github-event'],
        signature: webhookEvent['x-hub-signature-256'],
        payload: JSON.stringify(webhookEvent.body),
      });

      logger.debug('Webhook verified and received successfully');
    } catch (error) {
      console.error('Error while verifying and receiving webhook : ', error);
    }
  };
}

const app = express();

setVersionEndpoint(app);

app.use(express.json());

app.get('/', (req, res) => {
  res.status(200).send('OK');
});

app.post('/webhooks', async (req, res) => {
  try {
    logger.debug('Received webhook');

    await webhooks.verifyAndReceive({
      id: req.headers['x-request-id'] as string,
      name: req.headers['x-github-event'] as any,
      signature: req.headers['x-hub-signature-256'] as string,
      payload: JSON.stringify(req.body),
    });

    logger.debug('Webhook verified and received successfully');

    res.status(200).send('OK');
  } catch (error) {
    console.error('Error while verifying and receiving webhook : ', error);
    res.status(401).send('Unauthorized');
  }
});

app.listen(config.serverPort, () => {
  logger.info(`App listening on port ${config.serverPort}`);
});
