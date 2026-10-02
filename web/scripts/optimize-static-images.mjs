import fs from 'node:fs/promises';
import path from 'node:path';
import sharp from 'sharp';

const outputDir = path.resolve('dist');
const imageExtensions = new Set(['.jpg', '.jpeg', '.png', '.gif']);

/** @param {string} directory */
async function processDirectory(directory) {
  const entries = await fs.readdir(directory, { withFileTypes: true });
  for (const entry of entries) {
    const filePath = path.join(directory, entry.name);
    if (entry.isDirectory()) {
      await processDirectory(filePath);
      continue;
    }
    if (!imageExtensions.has(path.extname(entry.name).toLowerCase())) continue;
    const outputPath = path.join(directory, `${path.basename(entry.name, path.extname(entry.name))}.webp`);
    try {
      const source = await fs.stat(filePath);
      const existing = await fs.stat(outputPath).catch(() => null);
      if (existing && existing.mtimeMs >= source.mtimeMs) continue;
      await sharp(filePath).webp({ quality: 80 }).toFile(outputPath);
      console.log(`Optimized ${path.relative(outputDir, filePath)} to WebP`);
    } catch (error) {
      console.error(`Image optimization failed for ${filePath}`, error);
    }
  }
}

await processDirectory(outputDir);
