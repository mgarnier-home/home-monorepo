import esbuild, { BuildOptions } from 'esbuild';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const copyFiles = (srcDir, destDir) => {
  if (!fs.existsSync(destDir)) fs.mkdirSync(destDir);

  fs.readdirSync(srcDir).forEach((file) => {
    fs.copyFileSync(path.join(srcDir, file), path.join(destDir, file));
  });

  console.log('Copied files');
};

const main = async () => {
  const __filename = fileURLToPath(import.meta.url);

  const __dirname = path.dirname(__filename);

  const distDir = path.join(__dirname, '../app-dist');
  const publicDir = path.join(__dirname, '../app/public');

  const args = process.argv.slice(2);

  const context /*: BuildOptions*/ = {
    entryPoints: ['app/src/main.tsx'],
    bundle: true,
    outdir: 'app-dist',
    format: 'cjs',
    logLevel: 'info',
  };

  console.log('Args are :', args);

  if (fs.existsSync(distDir)) fs.rmSync(distDir, { recursive: true });

  if (args[0] === 'serve') {
    const ctx = await esbuild.context({ ...context });

    await ctx.watch();
    await ctx.serve({ servedir: distDir, port: 3000 });

    copyFiles(publicDir, distDir);

    fs.watch(publicDir, async (eventType, filename) => {
      copyFiles(publicDir, distDir);
    });
  } else if (args[0] === 'dev') {
    const ctx = await esbuild.context({ ...context });

    await ctx.watch();

    copyFiles(publicDir, distDir);

    fs.watch(publicDir, async (eventType, filename) => {
      copyFiles(publicDir, distDir);
    });
  } else if (args[0] === 'build') {
    await esbuild.build({ ...context, minify: true });

    console.log('Build done');

    copyFiles(publicDir, distDir);
  }
};

main();
