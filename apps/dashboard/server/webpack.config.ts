import path from 'path';
import TsconfigPathsPlugin from 'tsconfig-paths-webpack-plugin';
import webpack from 'webpack';

import { Args, Env, getConfig } from '../../../default.webpack.config';

const nodeProxyConfig = (env: Env, args: Args): webpack.Configuration => {
  const config = getConfig(env, args, 'dashboard');

  config.entry = {
    main: { import: './server/src/main.ts', filename: 'main.js' },
  };

  config.output = {
    ...config.output,
    path: path.resolve(__dirname, '../dist/server'),
  };

  return config;
};

export default nodeProxyConfig;
