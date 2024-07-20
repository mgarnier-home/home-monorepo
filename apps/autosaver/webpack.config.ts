import path from 'path';
import webpack from 'webpack';

import { Args, Env, getConfig } from '../../default.webpack.config';

const autosaverConfig = (env: Env, args: Args): webpack.Configuration => {
  const config = getConfig(env, args, 'autosaver');

  config.entry = {
    main: './src/main.ts',
  };

  console.log(config);

  return config;
};

export default autosaverConfig;
