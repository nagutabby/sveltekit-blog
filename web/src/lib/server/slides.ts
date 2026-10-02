import fs from 'node:fs';
import path from 'node:path';
import { ContentError } from './content';

const slidesDir = path.join(process.cwd(), 'static/content/slides');

export type Slide = { id: string; url: string };

export function slidesFromFiles(files: string[]): Slide[] {
  return files
    .filter((file) => file.endsWith('.pdf'))
    .sort()
    .map((file) => ({ id: file.slice(0, -4), url: `/content/slides/${file}` }));
}

export function getSlides(): Slide[] {
  if (!fs.existsSync(slidesDir)) {
    throw new ContentError(500, 'スライドディレクトリが見つかりません');
  }
  return slidesFromFiles(fs.readdirSync(slidesDir));
}

export function getSlide(id: string): Slide {
  const slide = getSlides().find((entry) => entry.id === id);
  if (!slide) throw new ContentError(404, `スライドが見つかりません: ${id}`);
  return slide;
}
