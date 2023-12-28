import esbuild, { BuildOptions } from 'esbuild';

const main = async () => {
  const args = process.argv.slice(2);

  const context: BuildOptions = {
    entryPoints: ['server/src/main.ts'],
    bundle: true,
    outdir: 'server-dist',
    logLevel: 'info',
  };

  console.log('Args are :', args);

  if (args[0] === 'dev') {
    const ctx = await esbuild.context({ ...context });

    await ctx.watch();

    console.log('Watching for changes...');
  } else if (args[0] === 'build') {
    await esbuild.build({ ...context, minify: true });

    console.log('Build done');
  }
};

main();
