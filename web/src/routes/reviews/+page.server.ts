import { getAllHTMLData } from '$lib/server/content';
import { paginate } from '$lib/server/pagination';
import type { PageServerLoad } from "./$types";
import type { Review } from '$lib/types/blog';

const PER_PAGE = 10;

export const load: PageServerLoad = async () => {
  const allReviews = await getAllHTMLData("reviews") as Review[];
  const { items: reviews, ...pagination } = paginate(allReviews, 1, PER_PAGE);

  return {
    image: "/images/Microsoft-Fluentui-Emoji-3d-Open-Book-3d.1024.png",
    title: "本のレビュー",
    body: "本を読んだ感想を掲載しています",
    reviews,
    pagination
  };
};
