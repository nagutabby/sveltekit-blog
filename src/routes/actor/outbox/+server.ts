import type { RequestEvent } from '@sveltejs/kit';
import { XMLParser } from 'fast-xml-parser';

async function fetchAtomFeed() {
  const response = await fetch('https://blog.nagutabby.uk/atom.xml');
  const xmlData = await response.text();

  const parser = new XMLParser({
    ignoreAttributes: false,
    attributeNamePrefix: "@_"
  });

  const feed = parser.parse(xmlData);
  return feed.feed.entry || [];
}

function convertTagUriToHttps(tagUri: string): string {
  const match = tagUri.match(/^tag:blog\.nagutabby\.uk,\d{4}-\d{2}-\d{2}:(.+)$/);
  if (match) {
    return `https://blog.nagutabby.uk${match[1]}`;
  }
  return tagUri;
}

function convertAtomToActivityPub(atomEntries: any[]) {
  return atomEntries.map(entry => ({
    id: `${convertTagUriToHttps(entry.id)}/create`,
    type: 'Create',
    actor: 'https://blog.nagutabby.uk/actor',
    published: entry.published,
    to: ['https://www.w3.org/ns/activitystreams#Public'],
    object: {
      id: convertTagUriToHttps(entry.id),
      type: 'Note',
      attributedTo: 'https://blog.nagutabby.uk/actor',
      name: entry.title,
      content: `${entry.title}\n\n${convertTagUriToHttps(entry.id)}`,
      published: entry.published,
      url: convertTagUriToHttps(entry.id),
      to: ['https://www.w3.org/ns/activitystreams#Public']
    }
  }));
}

export const GET = async ({ url }: RequestEvent) => {
  const page = url.searchParams.get('page');
  const pageSize = 10;

  try {
    const atomEntries = await fetchAtomFeed();
    const totalItems = atomEntries.length;

    if (!page) {
      // OutboxのCollection
      return new Response(JSON.stringify({
        '@context': 'https://www.w3.org/ns/activitystreams',
        id: 'https://blog.nagutabby.uk/actor/outbox',
        type: 'OrderedCollection',
        totalItems: totalItems,
        first: 'https://blog.nagutabby.uk/actor/outbox?page=1',
        last: `https://blog.nagutabby.uk/actor/outbox?page=${Math.ceil(totalItems / pageSize)}`
      }), {
        headers: {
          'Content-Type': 'application/activity+json',
          'Cache-Control': 'max-age=0, private, must-revalidate'
        }
      });
    }

    const pageNum = parseInt(page);
    const start = (pageNum - 1) * pageSize;
    const end = start + pageSize;
    const pageEntries = atomEntries.slice(start, end);
    const activities = convertAtomToActivityPub(pageEntries);

    return new Response(JSON.stringify({
      '@context': 'https://www.w3.org/ns/activitystreams',
      type: 'OrderedCollectionPage',
      id: `https://blog.nagutabby.uk/actor/outbox?page=${page}`,
      partOf: 'https://blog.nagutabby.uk/actor/outbox',
      orderedItems: activities,
      prev: pageNum > 1
        ? `https://blog.nagutabby.uk/actor/outbox?page=${pageNum - 1}`
        : null,
      next: end < totalItems
        ? `https://blog.nagutabby.uk/actor/outbox?page=${pageNum + 1}`
        : null
    }), {
      headers: {
        'Content-Type': 'application/activity+json',
        'Cache-Control': 'max-age=0, private, must-revalidate'
      }
    });
  } catch (error) {
    console.error('Error fetching or parsing atom feed:', error);
    return new Response('Internal Server Error', { status: 500 });
  }
};

// POSTは無効化（読み取り専用）
export const POST = async () => {
  return new Response('Method not allowed', { status: 405 });
};
