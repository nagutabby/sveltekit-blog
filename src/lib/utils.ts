import fs from 'fs';
import path from 'path';
import matter from 'gray-matter';
import { convertMarkdownToHtml, transformImagePath } from '$lib/markdown';
import { error } from '@sveltejs/kit';
import type { Article, ArticleFrontMatter } from '$lib/types/blog';
import type { Book, GoogleBooksResponse } from "$lib/types/book";

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

export async function getBook(isbn: string, apiKey: string): Promise<Book> {
  let url = `https://www.googleapis.com/books/v1/volumes?q=isbn:${isbn}&key=${apiKey}`;

  const response = await fetch(url);

  if (!response.ok) {
    throw new Error(`API request failed with status: ${response.status}`);
  }

  const data: GoogleBooksResponse = await response.json();

  const volumeInfo = data.items[0].volumeInfo;

  const book = {
    thumbnailUrl: volumeInfo.imageLinks?.thumbnail,
    infoLink: volumeInfo.infoLink,
    description: volumeInfo.description
  };
  return book;
}

export const getAllRawArticles = async () => {
  const articlesDir = path.join(process.cwd(), 'static/content/articles');

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

    const slug = fileName.replace(/\.md$/, '');

    return {
      id: slug,
      body: content,
      ...frontMatter
    } as Article;
  });

  const allArticles = await Promise.all(articlesPromises);

  const sortedArticles = allArticles.sort((a, b) => {
    const dateA = a.publishedAt ? new Date(a.publishedAt).getTime() : 0;
    const dateB = b.publishedAt ? new Date(b.publishedAt).getTime() : 0;
    return dateB - dateA;
  });

  return sortedArticles;
};

export const getAllHTMLArticles = async () => {
  const allArticles = await getAllRawArticles();
  allArticles.forEach(async article => {
    article.body = await convertMarkdownToHtml(article.body);
  });

  return allArticles;
};

export const getHTMLArticle = async (id: string) => {
  const fileName = `${id}.md`;

  const filePath = path.join(process.cwd(), 'static/content/articles', fileName);

  if (!fs.existsSync(filePath)) {
    throw error(404, `記事が見つかりません: ${id}`);
  }

  const fileContents = fs.readFileSync(filePath, 'utf8');

  const { data, content } = matter(fileContents);
  if (data.image) {
    data.image = transformImagePath(data.image);
  }

  const frontMatter = data as ArticleFrontMatter;

  const body = await convertMarkdownToHtml(content);

  const article: Article = { id, body, ...frontMatter };

  return article;
};
