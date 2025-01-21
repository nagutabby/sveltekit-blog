import type { RequestHandler } from '@sveltejs/kit';
import { getDetail, type Blog } from '$lib/microcms'; // ここで実際のデータファイルをインポートしてください

export const GET: RequestHandler = async ({ params }) => {
  const { name } = params;
  let articleData: Blog | null;

  // 記事データの取得
  if (name) {
    articleData = await getDetail(name);
    if (!articleData) {
      return new Response('Article not found', { status: 404 });
    }
  } else {
    return new Response('Article not found', { status: 404 });
  }

  // Noteオブジェクトの定義
  const note = {
    "id": `https://blog.nagutabby.uk/articles/${name}`,
    "type": "Note",
    "attributedTo": "https://blog.nagutabby.uk/actor",
    "name": articleData.title,
    "content": `<p>${articleData.title}</p><a href="https://blog.nagutabby.uk/articles/${name}" target="_blank">https://blog.nagutabby.uk/articles/${name}</a>`,
    "published": articleData.publishedAt,
    "url": `https://blog.nagutabby.uk/articles/${name}`,
    "to": ["https://www.w3.org/ns/activitystreams#Public"]
  };

  // Createアクティビティの定義
  const createActivity = {
    "@context": "https://www.w3.org/ns/activitystreams",
    "id": `https://blog.nagutabby.uk/articles/${name}/create`,
    "type": "Create",
    "actor": "https://blog.nagutabby.uk/actor",
    "published": new Date().toISOString(),
    "to": ["https://www.w3.org/ns/activitystreams#Public"],
    "object": note
  };

  return new Response(JSON.stringify(createActivity), {
    headers: {
      'Content-Type': 'application/activity+json'
    }
  });
};
