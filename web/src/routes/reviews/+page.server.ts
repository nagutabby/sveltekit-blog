
import { getAllHTMLData } from '$lib/utils';
import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from "./$types";
import type { Review } from '$lib/types/blog';

export const load: PageServerLoad = async ({ url }) => {
  const page = Number(url.searchParams.get("page")) || 1;
  const perPage = 10;

  const allReviews = await getAllHTMLData("reviews") as Review[]
  const totalReviews = allReviews.length;
  const totalPages = Math.ceil(totalReviews / perPage);

  if (page < 1 || page > totalPages && totalPages > 0) {
    redirect(302, '/');
  }

  const startIndex = (page - 1) * perPage;
  const endIndex = startIndex + perPage;
  const paginatedReviews = allReviews.slice(startIndex, endIndex);

  return {
    image: "/images/Microsoft-Fluentui-Emoji-3d-Open-Book-3d.1024.png",
    title: "本のレビュー",
    body: "本を読んだ感想を掲載しています",
    reviews: paginatedReviews,
    pagination: {
      currentPage: page,
      totalPages,
      hasNextPage: page < totalPages,
      hasPrevPage: page > 1
    }
  };
};
