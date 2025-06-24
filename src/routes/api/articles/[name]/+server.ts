import type { RequestHandler } from '@sveltejs/kit';
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
    "@context": "https://www.w3.org/ns/activitystreams",
    "id": `https://blog.nagutabby.uk/api/articles/${params.name}`,
    "type": "Note",
    "attributedTo": "https://blog.nagutabby.uk/actor",
    "name": article.title,
    "content": `<p>${article.title}</p><a href="https://blog.nagutabby.uk/articles/${params.name}" target="_blank">https://blog.nagutabby.uk/articles/${params.name}</a>`,
    "published": article.publishedAt,
    "url": `https://blog.nagutabby.uk/api/articles/${params.name}`,
    "to": ["https://www.w3.org/ns/activitystreams#Public"]
  };

  return new Response(JSON.stringify(note), {
    status: 200,
    headers: {
      'Content-Type': 'application/activity+json'
    }
  });
};
