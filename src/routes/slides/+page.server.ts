import fs from 'fs';
import path from 'path';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async () => {
  const slidesDir = path.join('static/content/slides');

  if (!fs.existsSync(slidesDir)) {
    throw error(500, 'スライドディレクトリが見つかりません');
  }

  const fileNames = fs.readdirSync(slidesDir).filter(file => file.endsWith('.pdf'));

  const slidesPromises = fileNames.map(async (fileName) => {
    const id = fileName.replace('.pdf', '');
    return {
      id: id,
      url: `/content/slides/${fileName}`
    };
  });

  const pdfs = await Promise.all(slidesPromises);

  return {
    image: "images/Microsoft-Fluentui-Emoji-3d-Package-3d.1024.png",
    title: "スライド一覧",
    body: "LT会で使ったスライドをまとめています",
    pdfs: pdfs,
  };
};
