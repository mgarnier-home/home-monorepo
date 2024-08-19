import path from 'path';
import webpack from 'webpack';

import { Args, Env, getConfig } from '../../default.webpack.config';

const nodeProxyConfig = (env: Env, args: Args): webpack.Configuration => {
  const config = getConfig(env, args, 'autoscaler');

  config.entry = {
    main: './src/main.ts',
  };

  return config;
};

export default nodeProxyConfig;
