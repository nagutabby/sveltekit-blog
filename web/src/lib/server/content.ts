import { convertMarkdownToHtml } from '$lib/markdown';
import { error } from '@sveltejs/kit';
import matter from 'gray-matter';
import fs from 'node:fs';
import path from 'node:path';
import type { Article, Review } from '$lib/types/blog';
import memoize from 'lodash.memoize';
import { dev } from '$app/environment';

type ContentType = 'articles' | 'reviews';

// Markdown+frontmatter source, shared with the Go backend's
// internal/content.Loader (which reads the same files for ActivityPub
// federation). Overridable via CONTENT_DIR for parity with the Go side;
// otherwise resolved relative to process.cwd(), which vite/SvelteKit set
// to the web/ project root for dev, build, and preview alike. Deliberately
// NOT resolved from import.meta.url: Vite bundles this module into a
// server chunk at build time, at an unrelated path/depth, so a path
// relative to *this source file* breaks once bundled.
const CONTENT_DIR = process.env.CONTENT_DIR ?? path.resolve(process.cwd(), '../backend/content');

class ContentNotFoundError extends Error {}

const toDate = (value: unknown): Date => {
  if (value instanceof Date) return value;
  if (typeof value === 'string' && value) return new Date(value);
  return new Date(0);
};

// Mirrors backend/internal/content.transformImagePath: a frontmatter
// image path like "images/foo.png" is rewritten to the URL web's static
// file server exposes it at.
const transformImagePath = (imagePath: string | undefined, contentType: ContentType): string => {
  if (imagePath?.startsWith('images/')) {
    return `/content/${contentType}/images/${path.basename(imagePath)}`;
  }
  return imagePath ?? '';
};

const readMarkdownFile = (contentType: ContentType, id: string) => {
  const filePath = path.join(CONTENT_DIR, contentType, `${id}.md`);
  let raw: string;
  try {
    raw = fs.readFileSync(filePath, 'utf-8');
  } catch (err) {
    if ((err as NodeJS.ErrnoException).code === 'ENOENT') {
      throw new ContentNotFoundError(`content not found: ${contentType}/${id}`);
    }
    throw err;
  }
  return matter(raw);
};

const listMarkdownIds = (contentType: ContentType): string[] =>
  fs
    .readdirSync(path.join(CONTENT_DIR, contentType))
    .filter((name) => name.endsWith('.md'))
    .map((name) => name.slice(0, -3))
    .sort();

const toArticle = (id: string, parsed: ReturnType<typeof matter>): Article => ({
  id,
  body: parsed.content,
  title: parsed.data.title ?? '',
  image: transformImagePath(parsed.data.image, 'articles'),
  publishedAt: toDate(parsed.data.publishedAt),
  updatedAt: toDate(parsed.data.updatedAt)
});

const toReview = (id: string, parsed: ReturnType<typeof matter>): Review => ({
  id,
  body: parsed.content,
  title: parsed.data.title ?? '',
  description: parsed.data.description ?? '',
  jp_e_code: parsed.data.jp_e_code ?? '',
  image: transformImagePath(parsed.data.image, 'reviews'),
  rating: parsed.data.rating ?? 0,
  publishedAt: toDate(parsed.data.publishedAt),
  updatedAt: toDate(parsed.data.updatedAt)
});

// Newest first, matching backend/internal/content.Loader's
// sort.SliceStable(... PublishedAt.After ...).
const byPublishedAtDescending = <T extends { publishedAt: Date }>(a: T, b: T) =>
  b.publishedAt.getTime() - a.publishedAt.getTime();

const listArticles = (): Article[] =>
  listMarkdownIds('articles')
    .map((id) => toArticle(id, readMarkdownFile('articles', id)))
    .sort(byPublishedAtDescending);

const listReviews = (): Review[] =>
  listMarkdownIds('reviews')
    .map((id) => toReview(id, readMarkdownFile('reviews', id)))
    .sort(byPublishedAtDescending);

const getAllRawDataImpl = async (type: ContentType): Promise<(Article | Review)[]> => {
  try {
    return type === 'articles' ? listArticles() : listReviews();
  } catch {
    throw error(500, 'コンテンツの取得に失敗しました');
  }
};

export const getAllRawData = dev ? getAllRawDataImpl : memoize(getAllRawDataImpl);

export const getAllHTMLData = async (type: ContentType) => {
  const allData = await getAllRawData(type);
  await Promise.all(
    allData.map(async (data) => {
      data.body = await convertMarkdownToHtml(data.body);
    })
  );

  return allData;
};

export const getHTMLData = async (id: string, type: ContentType): Promise<Article | Review> => {
  try {
    if (type === 'articles') {
      const article = toArticle(id, readMarkdownFile('articles', id));
      article.body = await convertMarkdownToHtml(article.body);
      return article;
    }

    const review = toReview(id, readMarkdownFile('reviews', id));
    review.body = await convertMarkdownToHtml(review.body);
    return review;
  } catch (err) {
    if (err instanceof ContentNotFoundError) {
      throw error(404, `記事が見つかりません: ${id}`);
    }
    throw error(500, 'コンテンツの取得に失敗しました');
  }
};
