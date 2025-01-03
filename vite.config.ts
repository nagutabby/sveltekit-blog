import { defineConfig } from "vitest/config";
import { sveltekit } from "@sveltejs/kit/vite";
import { dirname, join } from 'path';
import { fileURLToPath } from 'url';

const currentDir = dirname(fileURLToPath(import.meta.url))

export default defineConfig({
  // @ts-expect-error
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

  test: {
    include: ['src/**/*.{test,spec}.{js,ts}']
  }
});
