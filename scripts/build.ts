import esbuild, { BuildOptions } from 'esbuild';


const getProject = (args: string[]): string => {

  
  
}

const main = async () => {
  const args = process.argv.slice(2);

  const [script, project] = args;



  const context: BuildOptions = {
    entryPoints: ['src/main.ts'],
    bundle: true,
    outdir: 'dist',
    logLevel: 'info',
    platform: 'node',
    tsconfig: 'tsconfig.json',
    treeShaking: true,
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
