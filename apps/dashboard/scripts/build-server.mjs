import * as esbuild from 'esbuild'


await esbuild.build({
  entryPoints: ['server/src/main.ts'],
  bundle: true,
  outdir: 'server-dist',
  
  
})