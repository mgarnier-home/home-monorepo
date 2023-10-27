import { resolve } from "path";
import { defineConfig } from "vite";

import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  root: "./app",
  build: {
    outDir: "../app-dist",
  },
  resolve: {
    alias: {
      "@shared": resolve(__dirname, "./shared"),
    },
  },
});
