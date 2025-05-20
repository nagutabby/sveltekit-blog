import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
  const response = await fetch('/api/slides');

  if (!response.ok) {
    throw new Error('スライドのメタデータの取得に失敗しました');
  }

  const slides = await response.json();

  const pdfUrls = slides.map((slide: { name: string; }) =>
    slide.name.replace(/\.pdf$/, '')
  );

  return {
    pdfUrls
  };
};
