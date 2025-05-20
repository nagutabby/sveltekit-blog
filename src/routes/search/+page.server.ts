import { getAllRawArticles } from '$lib/utils';
import { redirect } from '@sveltejs/kit';
import { convertMarkdownToHtml } from '$lib/markdown';
import type { Article } from '$lib/types/blog';
import type { PageServerLoad } from "./$types";

interface ArticleWithRawBody extends Article {
  rawBody: string;
}

export const load: PageServerLoad = async ({ url }) => {
  const query = url.searchParams.get("q") || "";
  const page = Number(url.searchParams.get("page")) || 1;
  const perPage = 10;

  const allArticles = await getAllRawArticles();

  const allArticlesWithRawBody: ArticleWithRawBody[] = await Promise.all(
    allArticles.map(async (article) => {
      const articleWithRawBody = { ...article } as ArticleWithRawBody;
      articleWithRawBody.rawBody = article.body;
      articleWithRawBody.body = await convertMarkdownToHtml(article.body);
      return articleWithRawBody;
    })
  );

  const filteredArticles = query
    ? allArticlesWithRawBody.filter(article => {
      const content = article.rawBody.toLowerCase();
      const searchQuery = query.toLowerCase();

      return content.includes(searchQuery);
    })
    : allArticlesWithRawBody;

  const totalArticles = filteredArticles.length;
  const totalPages = Math.ceil(totalArticles / perPage);

  if (page < 1 || (page > totalPages && totalPages > 0)) {
    redirect(302, '/');
  }

  const startIndex = (page - 1) * perPage;
  const endIndex = startIndex + perPage;
  const paginatedArticles = filteredArticles.slice(startIndex, endIndex).map(article => {
    const { rawBody, ...rest } = article;
    return rest;
  });

  return {
    image: "/images/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png",
    title: `「${query}」を含む記事`,
    body: `「${query}」を含む記事の検索結果を表示しています`,
    articles: paginatedArticles,
    pagination: {
      currentPage: page,
      totalPages,
      hasNextPage: page < totalPages,
      hasPrevPage: page > 1
    }
  };
};
