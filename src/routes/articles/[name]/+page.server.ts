import fs from 'fs';
import path from 'path';
import matter from 'gray-matter';
import { convertMarkdownToHtml, transformImagePath } from '$lib/markdown';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from "./$types";
import type { Article, ArticleFrontMatter } from '$lib/types/blog';

export const load: PageServerLoad = async ({ params }) => {
  const fileName = `${params.name}.md`;

  const filePath = path.join('static/content/articles', fileName);

  if (!fs.existsSync(filePath)) {
    throw error(404, `記事が見つかりません: ${params.name}`);
  }

  const fileContents = fs.readFileSync(filePath, 'utf8');

  const { data, content } = matter(fileContents);
  if (data.image) {
    data.image = transformImagePath(data.image);
  }

  const frontMatter = data as ArticleFrontMatter;

  const body = await convertMarkdownToHtml(content);

  const article: Article = { id: params.name, body, ...frontMatter };

  return article;
};

