import fs from 'fs';
import path from 'path';
import matter from 'gray-matter';
import { convertMarkdownToHtml, transformImagePath } from '$lib/markdown';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from "./$types";
import type { Article, ArticleFrontMatter } from '$lib/types/blog';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ url }) => {
  const page = Number(url.searchParams.get("page")) || 1;
  const perPage = 10;
  const articlesDir = path.join('static/content/articles');

  if (!fs.existsSync(articlesDir)) {
    throw error(500, '記事ディレクトリが見つかりません');
  }

  const fileNames = fs.readdirSync(articlesDir).filter(file => file.endsWith('.md'));

  const articlesPromises = fileNames.map(async (fileName) => {
    const filePath = path.join(articlesDir, fileName);

    const fileContents = fs.readFileSync(filePath, 'utf8');

    const { data, content } = matter(fileContents);

    if (data.image) {
      data.image = transformImagePath(data.image);
    }

    const frontMatter = data as ArticleFrontMatter;

    const body = await convertMarkdownToHtml(content);

    const slug = fileName.replace(/\.md$/, '');

    return {
      id: slug,
      body,
      ...frontMatter
    } as Article;
  });

  const allArticles = await Promise.all(articlesPromises);

  const sortedArticles = allArticles.sort((a, b) => {
    const dateA = a.publishedAt ? new Date(a.publishedAt).getTime() : 0;
    const dateB = b.publishedAt ? new Date(b.publishedAt).getTime() : 0;
    return dateB - dateA;
  });

  const totalArticles = sortedArticles.length;
  const totalPages = Math.ceil(totalArticles / perPage);

  if (page < 1 || page > totalPages && totalPages > 0) {
    redirect(302, '/');
  }

  const startIndex = (page - 1) * perPage;
  const endIndex = startIndex + perPage;
  const paginatedArticles = sortedArticles.slice(startIndex, endIndex);

  return {
    image: "images/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png",
    title: "nagutabbyの考え事",
    body: "学んだことをまとめるブログ",
    articles: paginatedArticles,
    pagination: {
      currentPage: page,
      totalPages,
      hasNextPage: page < totalPages,
      hasPrevPage: page > 1
    }
  };
};
