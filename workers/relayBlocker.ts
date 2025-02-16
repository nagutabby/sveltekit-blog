export default {
  async fetch(
    request: Request
  ): Promise<Response> {
    try {
      const userAgent = request.headers.get('user-agent') || '';

      const url = new URL(request.url);
      const pathname = url.pathname;

      if (userAgent.toLowerCase().includes('relay') && pathname === '/actor/inbox') {
        const clonedRequest = request.clone();
        const activity = await clonedRequest.json();
        if (activity.type && activity.type == 'Accept') {
          return fetch(request);
        } else {
          return new Response(
            JSON.stringify({
              error: 'Forbidden',
              status: 403,
              message: 'Relay server access is not permitted',
              details: {
                reason: 'This server does not accept requests from relay servers',
                path: pathname,
                userAgent: userAgent
              }
            } as const),
            {
              status: 403,
              headers: {
                'Content-Type': 'application/activity+json',
              },
            }
          );
        }
      }

      return fetch(request);

    } catch (error) {
      return new Response(
        JSON.stringify({
          error: 'Internal Server Error',
          message: 'An error occurred while processing your request'
        } as const),
        {
          status: 500,
          headers: {
            'Content-Type': 'application/activity+json',
          }
        }
      );
    }
  },
};
