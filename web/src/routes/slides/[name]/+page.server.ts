import fs from 'fs';
import path from 'path';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from "./$types";

export function entries() {
  const slidesDir = path.join(process.cwd(), 'static/content/slides');

  if (!fs.existsSync(slidesDir)) {
    return [];
  }

  return fs
    .readdirSync(slidesDir)
    .filter((file) => file.endsWith('.pdf'))
    .map((file) => ({ name: file.replace('.pdf', '') }));
}

export const load: PageServerLoad = async ({ params }) => {
  const fileName = `${params.name}.pdf`;

  // このページはprerender(SSG)対象なのでloadはビルド時にしか実行されない。
  // entries()と同じくstatic/を見ればよい(理由はslides/+page.server.ts参照)。
  const filePath = path.join(process.cwd(), 'static/content/slides', fileName);

  if (!fs.existsSync(filePath)) {
    throw error(404, `スライドが見つかりません: ${params.name}`);
  }

  const url = `/content/slides/${fileName}`;

  return {
    url: url,
  };
};

