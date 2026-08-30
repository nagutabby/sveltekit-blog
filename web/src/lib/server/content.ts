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

// Mirrors backend/internal/content.transformImagePath: a frontmatter
// image path like "images/foo.png" is rewritten to the URL web's static
// file server exposes it at.
const transformImagePath = (imagePath: string | undefined, contentType: ContentType): string => {
  if (imagePath?.startsWith('images/')) {
    return `/content/${contentType}/images/${path.basename(imagePath)}`;
  }
  return imagePath ?? '';
};

// filenameDatePattern mirrors backend/internal/content.Loader's
// filenameDatePattern: articles/reviews are named by their publish date
// (YYYY-MM-DD.md), with an optional "-N" suffix to disambiguate multiple
// posts published on the same date.
const filenameDatePattern = /^(\d{4}-\d{2}-\d{2})(-\d+)?\.md$/;

// parsePublishedAtFromFilename derives a post's publish date from its
// filename, which is the source of truth now that frontmatter no longer
// carries a publishedAt field.
const parsePublishedAtFromFilename = (filename: string): Date => {
  const match = filenameDatePattern.exec(filename);
  if (!match) {
    throw new Error(`filename "${filename}" is not in YYYY-MM-DD[-N].md format`);
  }
  return new Date(match[1]);
};

type MarkdownEntry = {
  id: string;
  publishedAt: Date;
  parsed: ReturnType<typeof matter>;
};

const readMarkdownFile = (contentType: ContentType, filename: string) => {
  const filePath = path.join(CONTENT_DIR, contentType, filename);
  const raw = fs.readFileSync(filePath, 'utf-8');
  return matter(raw);
};

const listMarkdownEntries = (contentType: ContentType): MarkdownEntry[] =>
  fs
    .readdirSync(path.join(CONTENT_DIR, contentType))
    .filter((name) => name.endsWith('.md'))
    .sort()
    .map((filename) => {
      const publishedAt = parsePublishedAtFromFilename(filename);
      const parsed = readMarkdownFile(contentType, filename);
      return { id: parsed.data.id ?? '', publishedAt, parsed };
    });

const findMarkdownEntryById = (contentType: ContentType, id: string): MarkdownEntry => {
  const entry = listMarkdownEntries(contentType).find((e) => e.id === id);
  if (!entry) {
    throw new ContentNotFoundError(`content not found: ${contentType}/${id}`);
  }
  return entry;
};

const toArticle = (entry: MarkdownEntry): Article => ({
  id: entry.id,
  body: entry.parsed.content,
  title: entry.parsed.data.title ?? '',
  image: transformImagePath(entry.parsed.data.image, 'articles'),
  publishedAt: entry.publishedAt
});

const toReview = (entry: MarkdownEntry): Review => ({
  id: entry.id,
  body: entry.parsed.content,
  title: entry.parsed.data.title ?? '',
  description: entry.parsed.data.description ?? '',
  jp_e_code: entry.parsed.data.jp_e_code ?? '',
  image: transformImagePath(entry.parsed.data.image, 'reviews'),
  rating: entry.parsed.data.rating ?? 0,
  publishedAt: entry.publishedAt
});

// Newest first, matching backend/internal/content.Loader's
// sort.SliceStable(... PublishedAt.After ...).
const byPublishedAtDescending = <T extends { publishedAt: Date }>(a: T, b: T) =>
  b.publishedAt.getTime() - a.publishedAt.getTime();

const listArticles = (): Article[] =>
  listMarkdownEntries('articles').map(toArticle).sort(byPublishedAtDescending);

const listReviews = (): Review[] =>
  listMarkdownEntries('reviews').map(toReview).sort(byPublishedAtDescending);

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
      const article = toArticle(findMarkdownEntryById('articles', id));
      article.body = await convertMarkdownToHtml(article.body);
      return article;
    }

    const review = toReview(findMarkdownEntryById('reviews', id));
    review.body = await convertMarkdownToHtml(review.body);
    return review;
  } catch (err) {
    if (err instanceof ContentNotFoundError) {
      throw error(404, `記事が見つかりません: ${id}`);
    }
    throw error(500, 'コンテンツの取得に失敗しました');
  }
};
