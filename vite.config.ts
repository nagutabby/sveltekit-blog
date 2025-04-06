import { defineConfig } from "vitest/config";
import { sveltekit } from "@sveltejs/kit/vite";
import { svelteTesting } from '@testing-library/svelte/vite';

export default defineConfig({
  // @ts-expect-error
  plugins: [sveltekit(), svelteTesting()],

  build: {
    target: "ES2020",
  },
  test: {
    include: ['src/**/*.{test,spec}.{js,ts}'],
    environment: 'jsdom'
  }
});
