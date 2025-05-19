import type { PageServerLoad } from "./$types";
import fs from 'fs';
import path from 'path';
import matter from 'gray-matter';
import { convertMarkdownToHtml } from '$lib/markdown';
import { error } from '@sveltejs/kit';
import type { Article, ArticleFrontMatter } from '$lib/types/blog';

const transformImagePath = (imagePath: string): string => {
  if (imagePath && typeof imagePath === 'string' && imagePath.startsWith('images/')) {
    const fileName = path.basename(imagePath);
    return `/content/articles/images/${fileName}`;
  }
  return imagePath;
};

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

  const article: Article = { body, ...frontMatter };

  return article;
};

