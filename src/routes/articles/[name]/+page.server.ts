import type { Article } from "$lib/types/blog";
import { getHTMLData } from "$lib/utils";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params }) => {
  const article = await getHTMLData(params.name, "articles") as Article;
  return article;
};

