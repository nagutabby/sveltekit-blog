import fs from 'fs';
import path from 'path';
import matter from 'gray-matter';
import { convertMarkdownToHtml, transformImagePath } from '$lib/markdown';
import { error } from '@sveltejs/kit';
import type { Article, ArticleFrontMatter, Review, ReviewFrontMatter } from '$lib/types/blog';
import memoize from 'lodash.memoize';
import { dev } from '$app/environment';

export const generateDescriptionFromText = (body: string) => {
  const hasHTMLTags = /<[a-z][\s\S]*>/i.test(body);

  const plainText: string = hasHTMLTags ? body.replace(/<[^>]+>/g, "") : body;

  const normalizedText = plainText.replace(/\s+/g, " ").trim();
  let description = normalizedText.slice(0, 100);

  if (normalizedText.length > 100) {
    const lastChar = description.slice(-1);
    if (!/[。、．，!！?？]/.test(lastChar)) {
      description += "…";
    }
  }

  return description;
};

const getAllRawDataImpl = async (type: "articles" | "reviews") => {
  const dataDir = path.join(process.cwd(), `static/content/${type}`);

  try {
    await fs.promises.access(dataDir, fs.promises.constants.F_OK);
  } catch {
    throw error(500, 'ディレクトリが見つかりません');
  }


  const allFiles = await fs.promises.readdir(dataDir);
  const markdownFiles = allFiles.filter(file => file.endsWith('.md'));


  const dataPromises = markdownFiles.map(async (fileName) => {
    const filePath = path.join(dataDir, fileName);

    const fileContents = await fs.promises.readFile(filePath, 'utf8');

    const { data, content } = matter(fileContents);

    if (data.image) {
      data.image = transformImagePath(data.image, type);
    }

    const frontMatter = data as ArticleFrontMatter;

    const slug = fileName.replace(/\.md$/, '');

    return {
      id: slug,
      body: content,
      ...frontMatter
    } as Article | Review;
  });

  const allData = await Promise.all(dataPromises);

  const sortedData = allData.sort((a, b) => {
    const dateA = a.publishedAt ? new Date(a.publishedAt).getTime() : 0;
    const dateB = b.publishedAt ? new Date(b.publishedAt).getTime() : 0;
    return dateB - dateA;
  });

  return sortedData;
};

export const getAllRawData = dev ? getAllRawDataImpl : memoize(getAllRawDataImpl);

export const getAllHTMLData = async (type: "articles" | "reviews") => {
  const allData = await getAllRawData(type);
  await Promise.all(
    allData.map(async data => {
      data.body = await convertMarkdownToHtml(data.body);
    })
  );

  return allData;
};

export const getHTMLData = async (id: string, type: "articles" | "reviews"): Promise<Article | Review> => {
  const fileName = `${id}.md`;

  const filePath = path.join(process.cwd(), `static/content/${type}`, fileName);

  try {
    fs.promises.access(filePath, fs.promises.constants.F_OK);
  } catch {
    throw error(404, `記事が見つかりません: ${id}`);
  }

  const fileContents = await fs.promises.readFile(filePath, 'utf8');

  const { data, content } = matter(fileContents);
  if (data.image) {
    data.image = transformImagePath(data.image, type);
  }

  const body = await convertMarkdownToHtml(content);

  const isReviewFrontMatter = (obj: any) => {
    return 'rating' in obj;
  }

  let frontMatter, HTMLData;
  if (isReviewFrontMatter(data)) {
    frontMatter = data as ReviewFrontMatter
    HTMLData = { id, body, ...frontMatter } as Review
  } else {
    frontMatter = data as ArticleFrontMatter
    HTMLData = { id, body, ...frontMatter } as Article
  }

  return HTMLData;
};

export function getWebpPath(src: string) {
  if (!src) return '';

  // URLやデータURIの場合はそのまま返す
  if (src.startsWith('http') || src.startsWith('data:')) {
    return src;
  }

  const lastDotIndex = src.lastIndexOf('.');
  if (lastDotIndex === -1) return src;

  return src.substring(0, lastDotIndex) + '.webp';
}
