import prisma from '$lib/prisma';
import type { Config } from '@sveltejs/adapter-vercel';
import type { RequestHandler } from '@sveltejs/kit';
import { PRIVATE_KEY } from '$env/static/private';
import { signRequest } from '$lib/signRequest';

export const config: Config = {
  maxDuration: 30
};

const sendActivity = async (articleId: string, activityType: string, fetch: Function) => {
  try {
    const activity = await fetch(`/api/articles/${articleId}/${activityType}`);
    if (!activity.ok) {
      console.error(`Error fetching ${activityType} activity:`, await activity.text());
      return new Response(null, { status: 500 });
    }
    const activityData = await activity.json();

    const relayConnections = await prisma.relayConnection.findMany();
    if (!relayConnections.length) {
      console.log('No relay connections found');
      return new Response(null, { status: 200 });
    }

    const privateKeyPem = PRIVATE_KEY.split(String.raw`\n`).join('\n');

    for (const relayConnection of relayConnections) {
      const body = JSON.stringify(activityData);
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

      if (responseText && response.headers.get('content-type')?.includes('json')) {
        const responseData = JSON.parse(responseText);
        console.log(`${activityType} activity data:`, responseData);
      }
    }
    return new Response(null, { status: 200 });
  } catch (error) {
    console.error('Error:', error);
    return new Response(null, { status: 500 });
  }
};

export const POST: RequestHandler = async ({ request, fetch }) => {
  const data = await request.json();
  if (!data.id) {
    return new Response(null, { status: 400 });
  }
  switch (data.type) {
    case 'new':
      await sendActivity(data.id, 'create', fetch);
      break;
    case 'edit':
      await sendActivity(data.id, 'update', fetch);
      break;
    case 'delete':
      await sendActivity(data.id, data.type, fetch);
      break;
    default:
      return new Response(null, { status: 400 });
  }
  return new Response(null, { status: 200 });
};
