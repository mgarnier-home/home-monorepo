import { CronJob } from 'cron';
import express from 'express';
import { schedule } from 'node-cron';
import { execSync } from 'node:child_process';

import { setVersionEndpoint } from '@libs/api-version';
import { logger } from '@libs/logger';

import { config, getCronConfig } from './utils/config';
import { CronConfig, CronExecution } from './utils/interfaces';

logger.setAppName('stats-api');

logger.debug('Config loaded : ', config);

let lastConfig: CronConfig | undefined = undefined;
const cronsMap: Map<string, CronJob> = new Map();
const cronsExecutions: Map<string, CronExecution[]> = new Map();
const app = express();

setVersionEndpoint(app);

logger.debug('Loading cron config from : ', config.dataFilePath);

// API Setup
app.get('/', (req, res) => {
  res.status(200).send('OK');
});

app.get('/executions', (req, res) => {
  const obj: { [key: string]: CronExecution[] } = {};

  for (const [name, executions] of cronsExecutions.entries()) {
    obj[name] = executions;
  }

  res.status(200).send(obj);
});

app.listen(config.serverPort, () => {
  logger.info('Server started on port ' + config.serverPort);
});

// Cron execution
const addExecution = (cronName: string, execution: CronExecution) => {
  const executions = cronsExecutions.get(cronName) || [];

  executions.push(execution);

  if (executions.length > 10) {
    executions.shift();
  }

  cronsExecutions.set(cronName, executions);
};

const setupCronConfig = (config: CronConfig) => {
  config.crons.forEach((cron) => {
    logger.info(`Scheduling cron : ${cron.name} with schedule : ${cron.schedule}`);

    const job = new CronJob(
      cron.schedule,
      () => {
        logger.info(`Running cron : ${cron.name}`);
        logger.debug(`Executing command : ${cron.command}`);

        const cronExecution: CronExecution = {
          cronName: cron.name,
          date: new Date(),
          result: '',
          success: false,
        };

        try {
          const commandResult = execSync(cron.command, { stdio: 'pipe' });

          logger.debug(`Command result : ${commandResult}`);

          cronExecution.result = commandResult.toString();
          cronExecution.success = true;
        } catch (error) {
          logger.error(`Error executing command : ${error}`);

          cronExecution.result = String(error);
          cronExecution.success = false;
        }

        logger.info(`Cron : ${cron.name} executed`);
        addExecution(cron.name, cronExecution);
      },
      null,
      true,
      null,
      null,
      true
    );

    cronsMap.set(cron.name, job);
  });
};

// Config observer
const disposeOldConfig = () => {
  for (const [name, job] of cronsMap.entries()) {
    logger.info(`Disposing cron : ${name}`);
    job.stop();
  }

  cronsMap.clear();
};

const configFileChanged = async (config: CronConfig) => {
  disposeOldConfig();

  logger.info('Cronconfig loaded : ', config);

  setupCronConfig(config);
};

const checkConfig = async () => {
  const config = await getCronConfig();

  if (JSON.stringify(config) !== JSON.stringify(lastConfig)) {
    lastConfig = config;

    logger.info('Config file changed, triggering callback');

    configFileChanged(config);
  }
};

setInterval(checkConfig, 30 * 1000);

checkConfig();
