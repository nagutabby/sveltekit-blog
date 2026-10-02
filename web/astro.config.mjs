import { defineConfig } from 'astro/config';
import svelte from '@astrojs/svelte';
import tailwindcss from '@tailwindcss/vite';
import { loadEnv } from 'vite';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const root = path.dirname(fileURLToPath(import.meta.url));
const env = loadEnv(process.env.NODE_ENV === 'production' ? 'production' : 'development', root, '');
const contactServiceProxy = {
  target: env.BACKEND_URL ?? 'http://localhost:8080',
  changeOrigin: true
};

export default defineConfig({
  site: 'https://blog.nagutabby.uk',
  output: 'static',
  publicDir: './static',
  outDir: './dist',
  trailingSlash: 'never',
  build: {
    format: 'file'
  },
  integrations: [svelte()],
  vite: {
    plugins: [tailwindcss()],
    resolve: {
      alias: {
        '$lib': path.resolve(root, 'src/lib')
      }
    },
    server: {
      proxy: {
        '/blog.contact.v1.ContactService': contactServiceProxy
      }
    },
    preview: {
      proxy: {
        '/blog.contact.v1.ContactService': contactServiceProxy
      }
    }
  }
});
