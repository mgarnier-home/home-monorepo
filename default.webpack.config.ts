import { CleanWebpackPlugin } from 'clean-webpack-plugin';
import path from 'path';
import { RunScriptWebpackPlugin } from 'run-script-webpack-plugin';
import { TsconfigPathsPlugin } from 'tsconfig-paths-webpack-plugin';
import webpack from 'webpack';
import webpackCli from 'webpack-cli';
import nodeExternals from 'webpack-node-externals';

export type Env = webpackCli.WebpackRunOptions['env'] & {};

export type Args = webpackCli.WebpackRunOptions & {
  env: Env;
};

const tsConfigFile = 'tsconfig.json';

const defaultConfig: webpack.Configuration = {
  cache: {
    type: 'filesystem',
    cacheDirectory: path.resolve(__dirname, '.build_cache'),
  },
  target: 'node',
  node: {
    __filename: false,
    __dirname: false,
  },
  externals: [
    nodeExternals(),
    // {
    // modulesDir: path.join(__dirname, 'node_modules'),
    // importType: (moduleName) => `import ${moduleName}`,
    // }
  ],
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
      {
        test: /\.node$/,
        loader: 'node-loader',
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

export const getConfig = (env: Env, args: Args): webpack.Configuration => {
  defaultConfig.mode = args.mode || 'development';
  if (defaultConfig.mode === 'development') {
    defaultConfig.devtool = 'source-map';

    defaultConfig.watch = args.watch || false;
    if (defaultConfig.watch) {
      defaultConfig.watchOptions = {
        ignored: /node_modules/,
      };
      defaultConfig.plugins = [
        ...(defaultConfig.plugins as Array<any>),
        new webpack.WatchIgnorePlugin({ paths: [/\.js$/, /\.d\.ts$/] }),
        new RunScriptWebpackPlugin({
          name: 'main.js',

          nodeArgs: ['--inspect=0.0.0.0:9229'], // Allow debugging
          // nodeArgs: ['--env-file', path.join(getAppPath(app), '.env')],
          autoRestart: true, // auto
          // signal: true, // Signal to send for HMR (defaults to `false`, uses 'SIGUSR2' if `true`)
          // keyboard: true, // Allow typing 'rs' to restart the server. default: only if NODE_ENV is 'development'
          // args: ['scriptArgument1', 'scriptArgument2'], // pass args to script
          // cwd: getAppPath(app),
        }),
      ];
    }
  }

  return defaultConfig;
};
