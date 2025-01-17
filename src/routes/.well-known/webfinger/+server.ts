import type { RequestEvent } from '@sveltejs/kit';

export const GET = async ({ request }: RequestEvent): Promise<Response> => {
  const url = new URL(request.url);
  const resource = url.searchParams.get('resource');

  if (!resource) {
    return new Response('Resource parameter required', { status: 400 });
  }

  if (resource !== 'acct:article@blog.nagutabby.uk') {
    return new Response('User not found', { status: 404 });
  }

  return new Response(
    JSON.stringify({
      subject: 'acct:article@blog.nagutabby.uk',
      links: [
        {
          rel: 'self',
          type: 'application/activity+json',
          href: 'https://blog.nagutabby.uk/actor'
        },
        {
          rel: 'http://webfinger.net/rel/profile-page',
          type: 'text/html',
          href: 'https://blog.nagutabby.uk'
        }
      ]
    }),
    {
      headers: {
        'Content-Type': 'application/jrd+json',
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Methods': 'GET',
        'Access-Control-Allow-Headers': 'Content-Type'
      }
    }
  );
};
