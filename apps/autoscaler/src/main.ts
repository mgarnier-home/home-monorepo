// app.get('/', (req, res) => {
//   res.send('Hello World!');
// });
// app.listen(3000, () => {
//   console.log('App listening on port 3000');
// });
import EventSource from 'eventsource';
import express from 'express';

import { Webhooks } from '@octokit/webhooks';

// TODO : add a secret via env variable
const webhooks = new Webhooks({
  secret: 'SecureSecret',
});

webhooks.on('workflow_run', (event) => {
  console.log('===============================');
  console.log('workflow run event');
  console.log('workflow.name : ', event.payload.workflow?.name);
  console.log('workflow_run.name : ', event.payload.workflow_run.name);
  console.log('workflow_run.status : ', event.payload.workflow_run.status);
  console.log('workflow_run.conclusion : ', event.payload.workflow_run.conclusion);
  console.log('===============================');
});

webhooks.on('workflow_job', (event) => {
  console.log('===============================');
  console.log('workflow job event');
  console.log('workflow_job.name : ', event.payload.workflow_job.name);
  console.log('workflow_job.status : ', event.payload.workflow_job.status);
  console.log('workflow_job.conclusion : ', event.payload.workflow_job.conclusion);
  console.log('===============================');
});
//TODO : pourquoi node_env ne marche pas ?
console.log(process.env, process.env['MODE'], process.env.MODE === 'production');

if (process.env.MODE !== 'production') {
  console.log('Running in development mode');
  const webhookProxyUrl = 'https://smee.io/ZThAp1B0Bb5aTYu'; // replace with your own Webhook Proxy URL
  const source = new EventSource(webhookProxyUrl);
  source.onmessage = (event: any) => {
    const webhookEvent = JSON.parse(event.data);

    webhooks
      .verifyAndReceive({
        id: webhookEvent['x-request-id'],
        name: webhookEvent['x-github-event'],
        signature: webhookEvent['x-hub-signature-256'],
        payload: JSON.stringify(webhookEvent.body),
      })
      .catch(console.error);
  };
} else {
  console.log('Running in production mode');
  const app = express();

  app.use(express.json());

  app.post('/webhooks', async (req, res) => {
    try {
      await webhooks.verifyAndReceive({
        id: req.headers['x-request-id'] as string,
        name: req.headers['x-github-event'] as any,
        signature: req.headers['x-hub-signature-256'] as string,
        payload: JSON.stringify(req.body),
      });

      console.log('Webhook verified and received successfully');

      res.status(200).send('OK');
    } catch (error) {
      console.error('Error while verifying and receiving webhook : ', error);
      res.status(401).send('Unauthorized');
    }
  });

  app.listen(3000, () => {
    console.log('App listening on port 3000');
  });
}
