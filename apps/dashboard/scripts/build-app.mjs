import * as esbuild from 'esbuild'

await esbuild.build({
  entryPoints: ['app/src/main.tsx'],
  bundle: true,
  format: 'cjs',
  sourcemap: true,  
  outdir: 'app-dist',
})