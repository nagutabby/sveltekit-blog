import fs from 'fs';
import path from 'path';
import { error } from '@sveltejs/kit';
import { dev } from '$app/environment';
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params }) => {
  const fileName = `${params.name}.pdf`;

  // adapter-nodeの本番実行時はstatic/がコンテナに存在せず、build時にコピーされる
  // build/client/content/slidesを見る必要がある
  const filePath = path.join(process.cwd(), dev ? 'static/content/slides' : 'build/client/content/slides', fileName);

  if (!fs.existsSync(filePath)) {
    throw error(404, `スライドが見つかりません: ${params.name}`);
  }

  const url = `/content/slides/${fileName}`;

  return {
    url: url,
  };
};

