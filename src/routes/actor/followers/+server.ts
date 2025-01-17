import type { RequestHandler } from './$types';
import prisma from '$lib/prisma';

export const GET: RequestHandler = async () => {
  const followerCount = await prisma.follower.count({
    where: {
      following: true
    }
  });

  const response = new Response(
    JSON.stringify({
      "@context": "https://www.w3.org/ns/activitystreams",
      "id": "https://blog.nagutabby.uk/actor/followers",
      "type": "OrderedCollection",
      "totalItems": followerCount,
      "first": "https://blog.nagutabby.uk/actor/followers?page=1",
      "last": "https://blog.nagutabby.uk/actor/followers?page=1"
    }),
    {
      headers: {
        'Content-Type': 'application/activity+json',
        'Cache-Control': 'max-age=0, private, must-revalidate'
      }
    }
  );

  return response;
};
