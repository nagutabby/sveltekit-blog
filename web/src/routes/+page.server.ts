import { getAllHTMLData } from '$lib/server/content';
import { paginate } from '$lib/server/pagination';
import type { PageServerLoad } from "./$types";
import type { Article } from '$lib/types/blog';

const PER_PAGE = 10;

export const load: PageServerLoad = async () => {
  const allArticles = await getAllHTMLData("articles") as Article[];
  const { items: articles, ...pagination } = paginate(allArticles, 1, PER_PAGE);

  return {
    image: "/images/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png",
    title: "nagutabbyの考え事",
    body: "学んだことをまとめるブログ",
    articles,
    pagination
  };
};
