import { getAllRawData } from '$lib/server/content';
import { convertMarkdownToHtml } from '$lib/markdown';
import type { Article } from '$lib/types/blog';
import type { PageServerLoad } from "./$types";

export interface SearchableArticle extends Article {
  rawBody: string;
}

export const load: PageServerLoad = async () => {
  const allArticles = await getAllRawData("articles") as Article[];

  const articles: SearchableArticle[] = await Promise.all(
    allArticles.map(async (article) => ({
      ...article,
      rawBody: article.body,
      body: await convertMarkdownToHtml(article.body)
    }))
  );

  return {
    image: "/images/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png",
    articles
  };
};
