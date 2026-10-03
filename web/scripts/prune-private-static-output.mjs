import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { readdir, rm } from 'node:fs/promises';

const outputDirectory = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../dist');
const privateNames = new Set(['.obsidian', '.DS_Store']);
let removed = 0;

/** @param {string} directory */
async function prune(directory) {
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    const entryPath = path.join(directory, entry.name);
    if (privateNames.has(entry.name)) {
      await rm(entryPath, { recursive: true, force: true });
      removed += 1;
    } else if (entry.isDirectory()) {
      await prune(entryPath);
    }
  }
}

await prune(outputDirectory);
console.log(`Removed ${removed} local metadata entries from the static build.`);
