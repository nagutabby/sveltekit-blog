import { error } from '@sveltejs/kit';
import { getAllHTMLData } from '$lib/server/content';
import { paginate, totalPagesFor } from '$lib/server/pagination';
import type { PageServerLoad } from "./$types";
import type { Review } from '$lib/types/blog';

const PER_PAGE = 10;

export async function entries() {
  const allReviews = await getAllHTMLData("reviews") as Review[];
  const totalPages = totalPagesFor(allReviews.length, PER_PAGE);

  return Array.from({ length: Math.max(totalPages - 1, 0) }, (_, i) => ({ page: String(i + 2) }));
}

export const load: PageServerLoad = async ({ params }) => {
  const page = Number(params.page);
  const allReviews = await getAllHTMLData("reviews") as Review[];
  const totalPages = totalPagesFor(allReviews.length, PER_PAGE);

  if (!Number.isInteger(page) || page < 2 || page > totalPages) {
    throw error(404, 'ページが見つかりません');
  }

  const { items: reviews, ...pagination } = paginate(allReviews, page, PER_PAGE);

  return {
    image: "/images/Microsoft-Fluentui-Emoji-3d-Open-Book-3d.1024.png",
    title: "本のレビュー",
    body: "本を読んだ感想を掲載しています",
    reviews,
    pagination
  };
};
