import path from 'path';
import webpack from 'webpack';

import { Args, Env, getConfig } from '../../default.webpack.config';

const autosaverConfig = (env: Env, args: Args): webpack.Configuration => {
  const config = getConfig(env, args);

  config.entry = {
    main: './src/main.ts',
  };
  config.output = {
    filename: 'main.js',
    path: path.join(__dirname, 'dist'),
  };

  console.log(config);

  return config;
};

export default autosaverConfig;
