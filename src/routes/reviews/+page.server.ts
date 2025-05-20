import { GOOGLE_CLOUD_API_KEY } from "$env/static/private";
import { getBook } from "$lib/utils";
import fs from 'fs';
import path from 'path';
import matter from 'gray-matter';
import { convertMarkdownToHtml } from '$lib/markdown';
import { error } from '@sveltejs/kit';
import type { Review, ReviewFrontMatter } from '$lib/types/blog';
import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ url }) => {
  const page = Number(url.searchParams.get("page")) || 1;
  const perPage = 10;
  const reviewsDir = path.join('static/content/reviews');

  if (!fs.existsSync(reviewsDir)) {
    throw error(500, 'レビューディレクトリが見つかりません');
  }

  const fileNames = fs.readdirSync(reviewsDir).filter(file => file.endsWith('.md'));

  const reviewsPromises = fileNames.map(async (fileName) => {
    const filePath = path.join(reviewsDir, fileName);

    const fileContents = fs.readFileSync(filePath, 'utf8');

    const { data, content } = matter(fileContents);

    const frontMatter = data as ReviewFrontMatter;

    const body = await convertMarkdownToHtml(content);

    const slug = fileName.replace(/\.md$/, '');

    return {
      id: slug,
      body,
      ...frontMatter
    } as Review;
  });

  const allReviews = await Promise.all(reviewsPromises);

  const sortedReviews = allReviews.sort((a, b) => {
    const dateA = a.publishedAt ? new Date(a.publishedAt).getTime() : 0;
    const dateB = b.publishedAt ? new Date(b.publishedAt).getTime() : 0;
    return dateB - dateA;
  });

  const totalArticles = sortedReviews.length;
  const totalPages = Math.ceil(totalArticles / perPage);

  if (page < 1 || page > totalPages && totalPages > 0) {
    redirect(302, '/');
  }

  const startIndex = (page - 1) * perPage;
  const endIndex = startIndex + perPage;
  const paginatedReviews = sortedReviews.slice(startIndex, endIndex);

  const booksPromises = paginatedReviews.map(async (review) => {
    const isbn = review.isbn_13;
    const apiKey = GOOGLE_CLOUD_API_KEY;
    return await getBook(isbn, apiKey);
  });

  let books = await Promise.all(booksPromises);
  books = books.map(book => {
    if (book.thumbnailUrl.startsWith('http:')) {
      book.thumbnailUrl = book.thumbnailUrl.replace('http:', 'https:');
    }
    return book;
  });

  const mergedReviewData = paginatedReviews.map((review: Review, index: number) => ({
    review: review,
    book: books[index]
  }));

  return {
    image: "/images/Microsoft-Fluentui-Emoji-3d-Open-Book-3d.1024.png",
    title: "本のレビュー",
    body: "本を読んだ感想を掲載しています",
    content: mergedReviewData,
    pagination: {
      currentPage: page,
      totalPages,
      hasNextPage: page < totalPages,
      hasPrevPage: page > 1
    }
  };
};
