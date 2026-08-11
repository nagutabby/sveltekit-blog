import type { Article } from "$lib/types/blog";
import { getAllRawData, getHTMLData } from "$lib/server/content";
import type { PageServerLoad } from "./$types";

export async function entries() {
  const allArticles = await getAllRawData("articles") as Article[];
  return allArticles.map((article) => ({ name: article.id }));
}

export const load: PageServerLoad = async ({ params }) => {
  const article = await getHTMLData(params.name, "articles") as Article;
  return article;
};

