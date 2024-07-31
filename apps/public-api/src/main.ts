import EventSource from 'eventsource';
import express from 'express';

import { logger } from '@libs/logger';
import { Webhooks } from '@octokit/webhooks';

import { config } from './utils/config';

logger.setAppName('public-api');

logger.info('Webhooks is running');

logger.info('Config : ', config);

// const main = async () => {
//   const app = express();

//   setVersionEndpoint(app);

//   app.use(express.json());

//   app.get('/', (req, res) => {
//     res.send('OK');
//   });

//   app.post('/github', createNodeMiddleware(webhooks));

//   app.listen(config.serverPort, () => {
//     logger.info(`Server is running on port ${config.serverPort}`);
//   });
// };

// main();
const webhooks = new Webhooks({
  secret: config.secret,
});
const webhookProxyUrl = 'https://smee.io/VH3uwhHetkYuSq9p'; // replace with your own Webhook Proxy URL
const source = new EventSource(webhookProxyUrl);
console.log(config);
source.onmessage = async (event) => {
  const webhookEvent = JSON.parse(event.data);

  console.log(webhookEvent);
  console.log(webhookEvent.body.payload);

  const res = await webhooks.verify(webhookEvent.body.payload, webhookEvent['x-hub-signature-256']);

  console.log(res);
  // webhooks
  //   .verifyAndReceive({
  //     id: webhookEvent['x-request-id'],
  //     name: webhookEvent['x-github-event'],
  //     signature: webhookEvent['x-hub-signature-256'],
  //     payload: JSON.stringify(webhookEvent.body),
  //   })
  //   .catch(console.error);
};
