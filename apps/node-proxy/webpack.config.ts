import path from 'path';
import webpack from 'webpack';

import { Args, Env, getConfig } from '../../default.webpack.config';

const nodeProxyConfig = (env: Env, args: Args): webpack.Configuration => {
  const config = getConfig(env, args, 'node-proxy');

  config.entry = {
    main: './src/main.ts',
    worker: { import: './src/worker/proxyWorker.ts', filename: 'worker/proxyWorker.js' },
  };

  return config;
};

export default nodeProxyConfig;
