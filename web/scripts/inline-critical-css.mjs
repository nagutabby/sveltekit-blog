import Beasties from 'beasties';
import { readFile, writeFile } from 'node:fs/promises';
import { glob } from 'tinyglobby';

const files = await glob('build/**/*.html');

await Promise.all(files.map(async (file) => {
  const beasties = new Beasties({
    path: 'build',
    preload: 'media',
    pruneSource: false,
    mergeStylesheets: true,
  });
  const html = await readFile(file, 'utf-8');
  await writeFile(file, await beasties.process(html));
}));

console.log(`クリティカルCSSを${files.length}件のHTMLにインライン化しました`);
