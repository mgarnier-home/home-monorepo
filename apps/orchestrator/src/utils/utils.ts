import fs from 'fs';
import path from 'path';

import { config } from './config';

export const checkStackForHost = (stack: string, hostName: string) => fs.existsSync(getStackHostPath(stack, hostName));

export const getStackHostPath = (stack: string, host: string) =>
  path.join(config.composeFolderPath, stack, `${host}.${stack}.yml`);

export const getAdditionalComposeFiles = () => {
  // list all .yml files directly under config.composeFolderPath
  const files = fs.readdirSync(config.composeFolderPath);

  return files.filter((file) => file.endsWith('.yml')).map((file) => path.join(config.composeFolderPath, file));
};
