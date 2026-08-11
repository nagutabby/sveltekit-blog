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

export const GET = async () => {
  try {
    const atomEntries = await fetchAtomFeed();
    const totalItems = atomEntries.length;

    return new Response(JSON.stringify({
      '@context': 'https://www.w3.org/ns/activitystreams',
      id: 'https://blog.nagutabby.uk/actor/outbox',
      type: 'OrderedCollection',
      totalItems: totalItems
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
