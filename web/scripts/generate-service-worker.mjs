import { createHash } from 'node:crypto';
import fs from 'node:fs/promises';
import path from 'node:path';

const outputDir = path.resolve('dist');
const workerPath = path.join(outputDir, 'service-worker.js');

/** @param {string} directory @param {string} [relative] @returns {Promise<string[]>} */
async function listFiles(directory, relative = '') {
  const entries = await fs.readdir(directory, { withFileTypes: true });
  const files = [];
  for (const entry of entries) {
    const childRelative = path.posix.join(relative, entry.name);
    if (entry.isDirectory()) {
      files.push(...await listFiles(path.join(directory, entry.name), childRelative));
    } else if (childRelative !== 'service-worker.js') {
      files.push(`/${childRelative}`);
    }
  }
  return files;
}

const outputFiles = await listFiles(outputDir);
const files = new Set(outputFiles);
for (const file of outputFiles) {
  if (!file.endsWith('.html')) continue;
  const route = file.endsWith('/index.html')
    ? file.slice(0, -'index.html'.length) || '/'
    : file.slice(0, -'.html'.length);
  files.add(route);
}
const precacheFiles = [...files].sort();
const version = createHash('sha256').update(precacheFiles.join('\n')).digest('hex').slice(0, 12);
const precache = JSON.stringify(precacheFiles);
const source = `
const CACHE = 'cache-${version}';
const VERSION = '${version}';
const PRECACHE = ${precache};
const STATIC_ASSETS = new Set(PRECACHE);
self.addEventListener('install', (event) => {
  event.waitUntil(caches.open(CACHE).then((cache) => cache.addAll(PRECACHE)).then(() => self.skipWaiting()));
});
self.addEventListener('activate', (event) => {
  event.waitUntil(caches.keys().then(async (keys) => {
    for (const key of keys) if (key !== CACHE) await caches.delete(key);
    await self.clients.claim();
  }));
});
async function fetchAndCache(request) {
  const cache = await caches.open('offline-' + VERSION);
  try {
    const response = await fetch(request);
    if (response.ok || response.type === 'opaque') await cache.put(request, response.clone());
    return response;
  } catch (error) {
    const response = await cache.match(request);
    if (response) return response;
    throw error;
  }
}
self.addEventListener('fetch', (event) => {
  if (event.request.method !== 'GET' || event.request.headers.has('range')) return;
  const url = new URL(event.request.url);
  if (!url.protocol.startsWith('http')) return;
  const isStatic = url.origin === self.location.origin && STATIC_ASSETS.has(url.pathname);
  if (event.request.cache === 'only-if-cached' && !isStatic) return;
  event.respondWith((async () => {
    const cached = isStatic && await caches.match(event.request);
    return cached || fetchAndCache(event.request);
  })());
});
`;

await fs.writeFile(workerPath, source);
console.log(`Generated service worker with ${precacheFiles.length} precached URLs`);
