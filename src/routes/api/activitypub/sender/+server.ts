import prisma from '$lib/prisma';
import type { RequestHandler } from '@sveltejs/kit';
import { PRIVATE_KEY } from '$env/static/private';
import { signRequest } from '$lib/signRequest';

export const POST: RequestHandler = async ({ request, fetch }) => {
  const data = await request.json();

  if (data.id) {
    console.log(data);
    try {
      const createActivity = await fetch(`/api/articles/${data.id}/create`);
      if (!createActivity.ok) {
        console.error('Error fetching create activity:', await createActivity.text());
        return new Response(null, { status: 500 });
      }
      const createActivityData = await createActivity.json();

      const relayConnections = await prisma.relayConnection.findMany();
      if (!relayConnections.length) {
        console.log('No relay connections found');
        return new Response(null, { status: 200 });
      }

      const privateKeyPem = PRIVATE_KEY.split(String.raw`\n`).join('\n');

      for (const relayConnection of relayConnections) {
        try {
          const body = JSON.stringify(createActivityData);
          const headers = await signRequest(
            relayConnection.inbox,
            'POST',
            body,
            privateKeyPem
          );

          const response = await fetch(relayConnection.inbox, {
            method: 'POST',
            headers,
            body
          });

          const responseText = await response.text();

          if (!response.ok) {
            console.error('Error response from', relayConnection.inbox, responseText);
            continue;
          }

          try {
            if (responseText && response.headers.get('content-type')?.includes('json')) {
              const responseData = JSON.parse(responseText);
              console.log('Create activity data:', responseData);
            }
          } catch (parseError) {
            console.error('Error parsing JSON response:', parseError);
          }
        } catch (relayError) {
          console.error('Error sending to relay:', relayConnection.inbox, relayError);
        }
      }

      return new Response(null, { status: 200 });

    } catch (error) {
      console.error('Error in main process:', error);
      return new Response(null, { status: 500 });
    }
  }

  return new Response(null, { status: 400 });
};
