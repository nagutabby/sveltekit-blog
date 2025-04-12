import type { RequestHandler } from '@sveltejs/kit';
import { getArticleDetail, type Blog } from '$lib/microcms';

export const GET: RequestHandler = async ({ params }) => {
  const { name } = params;
  let articleData: Blog | null;

  // 記事データの取得
  if (name) {
    articleData = await getArticleDetail(name);
    if (!articleData) {
      return new Response('Article not found', { status: 404 });
    }
  } else {
    return new Response('Article not found', { status: 404 });
  }

  // Noteオブジェクトの定義
  const note = {
    "@context": "https://www.w3.org/ns/activitystreams",
    "id": `https://blog.nagutabby.uk/api/articles/${name}`,
    "type": "Note",
    "attributedTo": "https://blog.nagutabby.uk/actor",
    "name": articleData.title,
    "content": `<p>${articleData.title}</p><a href="https://blog.nagutabby.uk/articles/${name}" target="_blank">https://blog.nagutabby.uk/articles/${name}</a>`,
    "published": articleData.publishedAt,
    "url": `https://blog.nagutabby.uk/api/articles/${name}`,
    "to": ["https://www.w3.org/ns/activitystreams#Public"]
  };

  return new Response(JSON.stringify(note), {
    status: 200,
    headers: {
      'Content-Type': 'application/activity+json'
    }
  });
};
