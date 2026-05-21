import express from 'express';
import { readFileSync } from 'fs';
import { resolve } from 'path';

import { getEnvVariable } from '@libs/env-config';

export const setVersionEndpoint = (app: express.Express) => {
  app.get('/version', (req, res) => {
    try {
      const versionFilePath = getEnvVariable('VERSION_FILE_PATH', false, './version.txt');
      const version = readFileSync(resolve(__dirname, versionFilePath), 'utf8');

      res.status(200).send(version);
    } catch (error) {
      console.error(error);
      res.status(500).send('Something broke!');
    }
  });
};
