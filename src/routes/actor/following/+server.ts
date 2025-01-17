
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async () => {

    const response = new Response(
        JSON.stringify({
            "@context": "https://www.w3.org/ns/activitystreams",
            "id": "https://blog.nagutabby.uk/actor/following",
            "type": "OrderedCollection",
            "totalItems": 0,
            "first": "https://blog.nagutabby.uk/actor/following?page=1",
            "last": "https://blog.nagutabby.uk/actor/following?page=1"
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
