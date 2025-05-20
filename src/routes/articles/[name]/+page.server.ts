import { getHTMLArticle } from "$lib/utils";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params }) => {
  const article = await getHTMLArticle(params.name);
  return article;
};

