import type { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
  appId: 'uk.nagutabby.blog',
  appName: 'sveltekit-blog',
  webDir: 'build',
  server: {
    url: 'http://localhost:5173',
    cleartext: true
  }
};

export default config;
