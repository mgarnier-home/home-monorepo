import path from 'path';
import webpack from 'webpack';

import { Args, getConfig } from '../../default.webpack.config';

const nodeProxyConfig = (env: any, args: Args): webpack.Configuration => {
  const config = getConfig(args, 'autoscaler');

  config.entry = {
    main: './src/main.ts',
  };

  return config;
};

export default nodeProxyConfig;
