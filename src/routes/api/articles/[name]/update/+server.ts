import type { RequestHandler } from '@sveltejs/kit';
import { getDetail, type Blog } from '$lib/microcms';
import { signActivity } from '$lib/signRequest';
import { PRIVATE_KEY } from '$env/static/private';

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
    "id": `https://blog.nagutabby.uk/api/articles/${name}`,
    "type": "Note",
    "attributedTo": "https://blog.nagutabby.uk/actor",
    "name": articleData.title,
    "content": `<p>${articleData.title}</p><a href="https://blog.nagutabby.uk/articles/${name}" target="_blank">https://blog.nagutabby.uk/articles/${name}</a>`,
    "published": articleData.publishedAt,
    "url": `https://blog.nagutabby.uk/api/articles/${name}`,
    "to": ["https://www.w3.org/ns/activitystreams#Public"]
  };

  // Updateアクティビティの定義
  const activity = {
    "@context": [
      "https://www.w3.org/ns/activitystreams",
      "https://w3id.org/security/v1"
    ],
    "id": `https://blog.nagutabby.uk/api/articles/${name}/update`,
    "type": "Update",
    "actor": "https://blog.nagutabby.uk/actor",
    "published": new Date().toISOString(),
    "to": ["https://www.w3.org/ns/activitystreams#Public"],
    "object": note
  };

  const signatureData = await signActivity(activity, PRIVATE_KEY);

  const signedActivity = {
    ...activity,
    signature: signatureData
  };

  return new Response(JSON.stringify(signedActivity), {
    headers: {
      'Content-Type': 'application/activity+json'
    }
  });
};
