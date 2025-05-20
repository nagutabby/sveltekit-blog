
import { getAllHTMLArticles } from '$lib/utils';
import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ url }) => {
  const page = Number(url.searchParams.get("page")) || 1;
  const perPage = 10;

  const allArticles = await getAllHTMLArticles()
  const totalArticles = allArticles.length;
  const totalPages = Math.ceil(totalArticles / perPage);

  if (page < 1 || page > totalPages && totalPages > 0) {
    redirect(302, '/');
  }

  const startIndex = (page - 1) * perPage;
  const endIndex = startIndex + perPage;
  const paginatedArticles = allArticles.slice(startIndex, endIndex);

  return {
    image: "images/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png",
    title: "nagutabbyの考え事",
    body: "学んだことをまとめるブログ",
    articles: paginatedArticles,
    pagination: {
      currentPage: page,
      totalPages,
      hasNextPage: page < totalPages,
      hasPrevPage: page > 1
    }
  };
};
