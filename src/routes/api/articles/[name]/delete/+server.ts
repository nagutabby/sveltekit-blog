import type { RequestHandler } from '@sveltejs/kit';
import { getArticleDetail, type Blog } from '$lib/microcms';
import { signActivity } from '$lib/signRequest';
import { PRIVATE_KEY } from '$env/static/private';

export const GET: RequestHandler = async ({ params }) => {
  const { name } = params;

  // Noteオブジェクトの定義
  const note = {
    "id": `https://blog.nagutabby.uk/api/articles/${name}`,
    "type": "Note"
  };

  // Deleteアクティビティの定義
  const activity = {
    "@context": [
      "https://www.w3.org/ns/activitystreams",
      "https://w3id.org/security/v1"
    ],
    "id": `https://blog.nagutabby.uk/api/articles/${name}/delete`,
    "type": "Delete",
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
