import type { RequestHandler } from '@sveltejs/kit';

export const GET: RequestHandler = async ({ params, request }) => {
  const acceptHeader = request.headers.get('Accept');
  if (acceptHeader) {
    console.log('acceptHeader: ', acceptHeader);
  }
  return new Response(null, {
    status: 200,
    headers: {
      'Content-Type': 'application/activity+json'
    }
  });
};
