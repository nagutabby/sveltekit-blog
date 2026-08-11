import { error } from '@sveltejs/kit';
import { getAllHTMLData } from '$lib/server/content';
import { paginate, totalPagesFor } from '$lib/server/pagination';
import type { PageServerLoad } from "./$types";
import type { Article } from '$lib/types/blog';

const PER_PAGE = 10;

export async function entries() {
  const allArticles = await getAllHTMLData("articles") as Article[];
  const totalPages = totalPagesFor(allArticles.length, PER_PAGE);

  return Array.from({ length: Math.max(totalPages - 1, 0) }, (_, i) => ({ page: String(i + 2) }));
}

export const load: PageServerLoad = async ({ params }) => {
  const page = Number(params.page);
  const allArticles = await getAllHTMLData("articles") as Article[];
  const totalPages = totalPagesFor(allArticles.length, PER_PAGE);

  if (!Number.isInteger(page) || page < 2 || page > totalPages) {
    throw error(404, 'ページが見つかりません');
  }

  const { items: articles, ...pagination } = paginate(allArticles, page, PER_PAGE);

  return {
    image: "/images/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png",
    title: "nagutabbyの考え事",
    body: "学んだことをまとめるブログ",
    articles,
    pagination
  };
};
