import path from 'path';

import { config } from './config';

export const getStackHostPath = (stack: string, host: string) =>
  path.join(config.composeFolderPath, stack, `${host}.${stack}.yml`);
