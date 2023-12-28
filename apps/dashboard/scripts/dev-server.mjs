import * as esbuild from 'esbuild';

const context = await esbuild.context({ entryPoints: ['server/src/main.ts'], bundle: true, outdir: 'server-dist' });

await context.watch();
