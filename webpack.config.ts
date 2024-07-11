import { CleanWebpackPlugin } from 'clean-webpack-plugin';
import ForkTsCheckerWebpackPlugin from 'fork-ts-checker-webpack-plugin';
import path from 'path';
import { RunScriptWebpackPlugin } from 'run-script-webpack-plugin';
import { TsconfigPathsPlugin } from 'tsconfig-paths-webpack-plugin';
import webpack from 'webpack';
import webpackCli from 'webpack-cli';
import nodeExternals from 'webpack-node-externals';

const appsList = ['autosaver', 'dashboard'];

const tsConfigFile = 'tsconfig.json';

const config: webpack.Configuration = {
  output: {
    filename: `[name].js`,
    path: path.join(__dirname, 'dist'),
  },
  cache: {
    type: 'filesystem',
    cacheDirectory: path.resolve(__dirname, '.build_cache'),
  },
  // devtool: false,
  target: 'node',
  node: {
    __filename: false,
    __dirname: false,
  },
  externals: [nodeExternals()],
  module: {
    rules: [
      {
        test: /.ts?$/,
        use: [
          {
            loader: 'swc-loader',
            options: {
              sourceMaps: true,
              jsc: {
                target: 'es2022',
              },
            },
          },
        ],
        exclude: /node_modules/,
      },
    ],
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js'],
    plugins: [
      new TsconfigPathsPlugin({
        configFile: tsConfigFile,
      }),
    ],
  },
  plugins: [new webpack.ProgressPlugin(), new CleanWebpackPlugin()],
};

type Env = webpackCli.WebpackRunOptions['env'] & {
  apps?: string; // 'app1' | 'app1,app2'
};

type Args = webpackCli.WebpackRunOptions & {
  env: Env;
};

const getAppPath = (app: string) => path.join(__dirname, 'apps', app);

const getConfig = (env: Env, args: Args) => {
  const apps = env.apps ? env.apps.split(',') : appsList;

  config.entry = apps.reduce((acc: any, app) => {
    acc[app] = path.join(getAppPath(app), 'src', 'main.ts');
    return acc;
  }, {});

  console.log('config.entry', config.entry);

  config.mode = args.mode || 'none';

  if (config.mode === 'development') {
    config.watch = args.watch || false;
    config.devtool = 'source-map';

    if (config.watch && apps.length > 1) {
      console.error('Cannot run multiple apps in watch mode');
      process.exit(1);
    }
    if (config.watch && apps.length === 1) {
      const app = apps[0];
      config.watchOptions = {
        ignored: /node_modules/,
      };

      // config.cache = { ...(config.cache as webpack.FileCacheOptions), name: 'development_hmr' };
      config.entry = {
        main: [(config.entry as any)[app]],
      };
      // config.externals = [
      //   nodeExternals({
      //     allowlist: ['webpack/hot/poll?100'],
      //   }),
      // ];
      console.log(path.join(getAppPath(app), '.env'));
      config.plugins = [
        ...(config.plugins as Array<any>),
        new webpack.WatchIgnorePlugin({ paths: [/\.js$/, /\.d\.ts$/] }),
        // new webpack.HotModuleReplacementPlugin({}),
        new RunScriptWebpackPlugin({
          name: 'main.js',

          nodeArgs: ['--inspect=0.0.0.0:9229'], // Allow debugging
          // nodeArgs: ['--env-file', path.join(getAppPath(app), '.env')],
          autoRestart: true, // auto
          // signal: true, // Signal to send for HMR (defaults to `false`, uses 'SIGUSR2' if `true`)
          // keyboard: true, // Allow typing 'rs' to restart the server. default: only if NODE_ENV is 'development'
          // args: ['scriptArgument1', 'scriptArgument2'], // pass args to script
          cwd: getAppPath(app),
        }),
      ];
    }
  }

  // if (argv.mode === 'production') {
  //   config.devtool = 'source-map';
  //   config.cache.name = 'production';
  //   config.optimization.minimize = true;
  //   config.optimization.minimizer = [
  //     new TerserPlugin({
  //       terserOptions: {
  //         keep_classnames: true,
  //         keep_fnames: true,
  //       },
  //     }),
  //   ];
  // }

  return config;
};

module.exports = getConfig;
