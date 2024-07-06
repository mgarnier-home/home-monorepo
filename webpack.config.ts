import { CleanWebpackPlugin } from 'clean-webpack-plugin';
import ForkTsCheckerWebpackPlugin from 'fork-ts-checker-webpack-plugin';
import path from 'path';
import { RunScriptWebpackPlugin } from 'run-script-webpack-plugin';
import { TsconfigPathsPlugin } from 'tsconfig-paths-webpack-plugin';
import webpack from 'webpack';
import { ConfigOptions } from 'webpack-cli';
import nodeExternals from 'webpack-node-externals';

const tsConfigFile = 'tsconfig.json';

const config: webpack.Configuration = {
  output: {
    filename: `[name].js`,
    path: path.join(__dirname, 'dist'),
    devtoolModuleFilenameTemplate: `${path.sep}[absolute-resource-path][loaders]`,
  },
  cache: {
    type: 'filesystem',
    cacheDirectory: path.resolve(__dirname, '.build_cache'),
  },
  // devtool: false,
  target: 'node',
  mode: 'development',
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
              sourceMaps: false,
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
  plugins: [
    new webpack.ProgressPlugin(),
    new CleanWebpackPlugin(),
    // only in lint mode
    // new ForkTsCheckerWebpackPlugin({
    //   typescript: {
    //     configFile: tsConfigFile,
    //     configOverwrite: {
    //       include: ['apps/autosaver'],
    //       exclude: ['node_modules'],
    //     },
    //     diagnosticOptions: {
    //       semantic: true,
    //       syntactic: true,
    //     },
    //     mode: 'write-references',
    //   },
    //   issue: {},
    // }),
  ],
};

module.exports = (options: ConfigOptions) => {
  console.log(options);
  // config.mode = argv.mode;

  // if (argv.mode === 'development') {
  // config.watch = argv.watch;
  config.watch = true;
  config.watchOptions = {
    // for some systems, watching many files can result in a lot of CPU or memory usage
    // https://webpack.js.org/configuration/watch/#watchoptionsignored
    // don't use this pattern, if you have a monorepo with linked packages
    ignored: /node_modules/,
  };
  // config.cache.name = env.debug ? 'development_debug' : 'development';
  // config.devtool = env.debug ? 'eval-source-map' : undefined;

  // If running in watch mode
  if (config.watch) {
    // config.cache.name = env.debug ? 'development_debug_hmr' : 'development_hmr';
    // config.entry = {
    //   ...(config as any).entry,
    //   main: ['webpack/hot/poll?100',
    //     //'webpack/hot/signal',
    //     (config as any).entry.main],
    // };
    config.entry = ['webpack/hot/poll?100', `./apps/autosaver/src/main.ts`];
    config.externals = [
      nodeExternals({
        allowlist: [
          'webpack/hot/poll?100',
          //'webpack/hot/signal'
        ],
      }),
    ];
    config.plugins = [
      ...(config.plugins as Array<any>),
      new webpack.WatchIgnorePlugin({ paths: [/\.js$/, /\.d\.ts$/] }),
      new webpack.HotModuleReplacementPlugin({}),
      new RunScriptWebpackPlugin({
        name: 'main.js',
        // nodeArgs: env.debug ? ['--inspect'] : undefined, // Allow debugging
        nodeArgs: ['--env-file', '../.env'],
        autoRestart: false, // auto
        // signal: true, // Signal to send for HMR (defaults to `false`, uses 'SIGUSR2' if `true`)
        keyboard: true, // Allow typing 'rs' to restart the server. default: only if NODE_ENV is 'development'
        // args: ['scriptArgument1', 'scriptArgument2'], // pass args to script
        cwd: path.join(__dirname, 'dist'),
      }),
    ];
  }
  // }

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
