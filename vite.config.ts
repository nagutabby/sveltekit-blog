import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";
import { dirname, join } from 'path';
import { fileURLToPath } from 'url';

const currentDir = dirname(fileURLToPath(import.meta.url))

export default defineConfig({
  plugins: [sveltekit()],

  build: {
    target: "ES2020",
  },

  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `@use "${join(currentDir, './src/variables.scss')}" as *;`,
      },
    },
  },
});
