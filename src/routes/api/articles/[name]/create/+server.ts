import type { RequestHandler } from '@sveltejs/kit';
import { signActivity } from '$lib/signRequest';
import { PRIVATE_KEY } from '$env/static/private';
import { getHTMLData } from '$lib/utils';
import type { Article } from '$lib/types/blog';

export const GET: RequestHandler = async ({ params }) => {
  let article: Article;

  if (params.name) {
    article = await getHTMLData(params.name, "articles");
  } else {
    return new Response('Article not found', { status: 404 });
  }

  const note = {
    "id": `https://blog.nagutabby.uk/api/articles/${params.name}`,
    "type": "Note",
    "attributedTo": "https://blog.nagutabby.uk/actor",
    "name": article.title,
    "content": `<p>${article.title}</p><a href="https://blog.nagutabby.uk/articles/${params.name}" target="_blank">https://blog.nagutabby.uk/articles/${params.name}</a>`,
    "published": article.publishedAt,
    "url": `https://blog.nagutabby.uk/api/articles/${params.name}`,
    "to": ["https://www.w3.org/ns/activitystreams#Public"]
  };

  const activity = {
    "@context": [
      "https://www.w3.org/ns/activitystreams",
      "https://w3id.org/security/v1"
    ],
    "id": `https://blog.nagutabby.uk/api/articles/${params.name}/create`,
    "type": "Create",
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
