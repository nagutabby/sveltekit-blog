import type { RequestEvent } from '@sveltejs/kit';
import prisma from '$lib/prisma';
import crypto from 'crypto';
import { PRIVATE_KEY } from '$env/static/private';
import { signRequest } from '$lib/signRequest';

class ActivityPubError extends Error {
  constructor(
    message: string,
    public statusCode: number = 500,
    public details?: unknown
  ) {
    super(message);
    this.name = 'ActivityPubError';
  }
}

class ValidationError extends ActivityPubError {
  constructor(message: string, details?: unknown) {
    super(message, 400, details);
    this.name = 'ValidationError';
  }
}

interface ActorInfo {
  inbox: string;
  publicKey: {
    id: string;
    owner: string;
    publicKeyPem: string;
  };
}

async function fetchActor(actorUrl: string): Promise<ActorInfo | null> {
  try {
    const response = await fetch(actorUrl, {
      headers: {
        'Accept': 'application/activity+json'
      }
    });

    if (!response.ok) {
      console.log(`Actor fetch failed: ${response.status} ${response.statusText}`);
      throw new Error(`${response.status}`);
    }

    const actor = await response.json();
    if (!actor.inbox || !actor.publicKey?.publicKeyPem) {
      console.error('Invalid actor data:', actor);
      throw new Error('Missing inbox or public key in actor object');
    }

    return {
      inbox: actor.inbox,
      publicKey: actor.publicKey
    };
  } catch (error) {
    console.error('Error fetching actor:', error);
    return null;
  }
}

async function fetchOutboxPage(url: string): Promise<any> {
  const response = await fetch(url, {
    headers: {
      'Accept': 'application/activity+json'
    }
  });

  if (!response.ok) {
    throw new Error(`Failed to fetch outbox page: ${response.status}`);
  }

  return response.json();
}

async function sendAcceptActivity(activity: any, targetInbox: string, privateKeyPem: string) {
  const accept = {
    '@context': 'https://www.w3.org/ns/activitystreams',
    id: `https://blog.nagutabby.uk/activities/${crypto.randomUUID()}`,
    type: 'Accept',
    actor: 'https://blog.nagutabby.uk/actor',
    object: activity
  };

  const body = JSON.stringify(accept);

  try {
    const headers = await signRequest(
      targetInbox,
      'POST',
      body,
      privateKeyPem
    );

    const response = await fetch(targetInbox, {
      method: 'POST',
      headers,
      body
    });

    if (!response.ok) {
      const responseText = await response.text();
      console.error('Accept activity failed:', response.status, response.statusText);
      console.error('Response body:', responseText);
      throw new Error(`HTTP ${response.status}: ${responseText}`);
    }

    return true;
  } catch (error) {
    console.error('Error sending Accept:', error);
    throw error;
  }
}

export const POST = async ({ request }: RequestEvent) => {
  try {
    const activity = await request.json();

    if (!activity['@context'] || !activity.type || !activity.actor) {
      throw new ValidationError('Invalid activity: missing required fields');
    }

    // アクターの情報を取得
    const actorInfo = await fetchActor(activity.actor);
    if (!actorInfo) {
      return new Response('Could not fetch actor information', { status: 400 });
    }

    const privateKeyPem = PRIVATE_KEY.split(String.raw`\n`).join('\n');

    try {
      switch (activity.type) {
        case 'Follow': {
          await prisma.follower.upsert({
            where: { actorId: activity.actor },
            update: {
              following: true,
              inbox: actorInfo.inbox,
              publicKeyPem: actorInfo.publicKey.publicKeyPem
            },
            create: {
              actorId: activity.actor,
              inbox: actorInfo.inbox,
              publicKeyPem: actorInfo.publicKey.publicKeyPem,
              following: true
            }
          });

          // Acceptを送信
          await sendAcceptActivity(activity, actorInfo.inbox, privateKeyPem);

          return new Response(JSON.stringify({
            '@context': 'https://www.w3.org/ns/activitystreams',
            id: `https://blog.nagutabby.uk/activities/${crypto.randomUUID()}`,
            type: 'Accept',
            actor: 'https://blog.nagutabby.uk/actor',
            object: activity
          }), {
            headers: {
              'Content-Type': 'application/activity+json'
            }
          });
        }

        case 'Undo': {
          if (activity.object?.type === 'Follow') {
            await prisma.follower.update({
              where: { actorId: activity.actor },
              data: {
                following: false,
                inbox: actorInfo.inbox,
                publicKeyPem: actorInfo.publicKey.publicKeyPem
              }
            });

            await sendAcceptActivity(activity, actorInfo.inbox, privateKeyPem);

            return new Response(JSON.stringify({
              '@context': 'https://www.w3.org/ns/activitystreams',
              type: 'Accept',
              actor: 'https://blog.nagutabby.uk/actor',
              object: activity
            }), {
              headers: {
                'Content-Type': 'application/activity+json'
              }
            });
          }
          return new Response('Invalid Undo activity', { status: 400 });
        }

        case 'Accept': {
          if (activity.object?.type === 'Follow' || activity.object?.type === 'Subscribe') {
            // アクターの情報を取得
            const actorInfo = await fetchActor(activity.actor);
            if (!actorInfo) {
              return new Response('Could not fetch actor information', { status: 400 });
            }

            await prisma.relayConnection.upsert({
              where: {
                actorId: activity.actor
              },
              update: {
                connected: true,
                lastAcceptedAt: new Date(),
                inbox: actorInfo.inbox
              },
              create: {
                actorId: activity.actor,
                connected: true,
                lastAcceptedAt: new Date(),
                inbox: actorInfo.inbox
              }
            });

            return new Response('', { status: 200 });
          }
          return new Response('Invalid Accept activity', { status: 400 });
        }

        default:
          const message = `${activity.type} activity is not supported`;
          console.log(message);
          return new Response(message, { status: 422 });
      }
    } catch (error) {
      console.error('Error processing activity:', error);
      throw error;
    }
  } catch (error: unknown) {
    console.error('Error processing request:', error);

    if (error instanceof ActivityPubError) {
      return new Response(error.message, {
        status: error.statusCode,
        headers: { 'Content-Type': 'application/json' }
      });
    }

    if (error instanceof Error) {
      return new Response(error.message, {
        status: 500,
        headers: { 'Content-Type': 'application/json' }
      });
    }

    return new Response('An unknown error occurred', {
      status: 500,
      headers: { 'Content-Type': 'application/json' }
    });
  }
};
