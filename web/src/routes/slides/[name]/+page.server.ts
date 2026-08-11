import fs from 'fs';
import path from 'path';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params }) => {
  const fileName = `${params.name}.pdf`;

  const filePath = path.join(process.cwd(), 'static/content/slides', fileName);

  if (!fs.existsSync(filePath)) {
    throw error(404, `スライドが見つかりません: ${params.name}`);
  }

  const url = `/content/slides/${fileName}`;

  return {
    url: url,
  };
};

