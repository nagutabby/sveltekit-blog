import { generateDescriptionFromText } from "$lib/utils";
import fs from 'fs';
import path from 'path';
import matter from 'gray-matter';
import { convertMarkdownToHtml, transformImagePath } from '$lib/markdown';
import { error } from '@sveltejs/kit';
import type { Article, ArticleFrontMatter } from '$lib/types/blog';

let latestDate: Date | undefined;

function create_entry(title: string, body: string, path: string, publishedAt: Date, updatedAt: Date) {
  const publishedDate = new Date(publishedAt).toISOString().substring(0, 10);

  if (latestDate === undefined) {
    latestDate = new Date(updatedAt);
  } else if (latestDate < new Date(updatedAt)) {
    latestDate = new Date(updatedAt);
  }

  const formattedPublishedAt = new Date(publishedAt).toISOString();
  const formattedUpdatedAt = new Date(updatedAt).toISOString();

  return `<entry>
  <title>${title}</title>
  <summary type="text">${generateDescriptionFromText(body)}</summary>
  <link href="${new URL(`/articles/${path}`, 'https://blog.nagutabby.uk').href}" rel="alternate" />
  <updated>${formattedUpdatedAt}</updated>
  <published>${formattedPublishedAt}</published>
  <id>tag:blog.nagutabby.uk,${publishedDate}:/articles/${path}</id>
  </entry>`;
}

export async function GET({ setHeaders }) {
  setHeaders({
    'Content-Type': 'application/xml'
  });

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

  const posts = sortedArticles.map((post) => create_entry(post.title, post.body, post.id, post.publishedAt, post.updatedAt));

  const atom = `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
<link rel="self" href="https://blog.nagutabby.uk/atom.xml" type="application/rss+xml" />
<title>nagutabbyの考え事</title>
<link href="https://blog.nagutabby.uk" />
<updated>${latestDate ? latestDate.toISOString() : new Date().toISOString()}</updated>
<author><name>nagutabby</name></author>
<id>tag:blog.nagutabby.uk,2023-01-01:/</id>
${posts.join('\n')}
</feed>`;


  return new Response(atom);

}
