import esbuild from 'esbuild';
// import fs from 'node:fs';
// import path from 'node:path';
// import { fileURLToPath } from 'node:url';

// const __filename = fileURLToPath(import.meta.url);

// const __dirname = path.dirname(__filename);

const args = process.argv.slice(2);

const context = {
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
