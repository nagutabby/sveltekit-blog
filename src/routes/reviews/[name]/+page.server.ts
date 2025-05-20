import { GOOGLE_CLOUD_API_KEY } from "$env/static/private";
import { getBook } from "$lib/utils";
import fs from 'fs';
import path from 'path';
import matter from 'gray-matter';
import { convertMarkdownToHtml } from '$lib/markdown';
import { error } from '@sveltejs/kit';
import type { Review, ReviewFrontMatter } from '$lib/types/blog';
import type { PageServerLoad } from "./$types";

const convertHalfToFullWidthPunctuations = (text: string): string => {
  return text.replace(/\?/g, '？').replace(/\!/g, '！');
};

export const load: PageServerLoad = async ({ params }) => {
  const fileName = `${params.name}.md`;

  const filePath = path.join('static/content/reviews', fileName);

  if (!fs.existsSync(filePath)) {
    throw error(404, `記事が見つかりません: ${params.name}`);
  }

  const fileContents = fs.readFileSync(filePath, 'utf8');

  const { data, content } = matter(fileContents);

  const frontMatter = data as ReviewFrontMatter;

  const body = await convertMarkdownToHtml(content);

  const slug = fileName.replace(/\.md$/, '');

  const review = {
    id: slug,
    body,
    ...frontMatter
  } as Review;

  const apiKey = GOOGLE_CLOUD_API_KEY;
  const book = await getBook(review.isbn_13, apiKey);

  if (book.thumbnailUrl.startsWith('http:')) {
    book.thumbnailUrl = book.thumbnailUrl.replace('http:', 'https:');
  }

  if (book.description.includes('!') || book.description.includes('?')) {
    book.description = convertHalfToFullWidthPunctuations(book.description);
  }

  return {
    review,
    book,
  };
};
