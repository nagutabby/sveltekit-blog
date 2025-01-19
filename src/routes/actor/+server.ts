import { PUBLIC_PUBLIC_KEY } from '$env/static/public';

export const GET = async () => {
  return new Response(
    JSON.stringify({
      "@context": [
        "https://www.w3.org/ns/activitystreams",
        "https://w3id.org/security/v1"
      ],
      "id": "https://blog.nagutabby.uk/actor",
      "type": "Service",
      "preferredUsername": "article",
      "name": "nagutabbyの考え事",
      "summary": "ブログ投稿用アカウント",
      "url": "https://blog.nagutabby.uk",
      "inbox": "https://blog.nagutabby.uk/actor/inbox",
      "outbox": "https://blog.nagutabby.uk/actor/outbox",
      "following": "https://blog.nagutabby.uk/actor/following",
      "followers": "https://blog.nagutabby.uk/actor/followers",
      "discoverable": true,
      "publicKey": {
        "id": "https://blog.nagutabby.uk/actor#main-key",
        "owner": "https://blog.nagutabby.uk/actor",
        "publicKeyPem": PUBLIC_PUBLIC_KEY
      }
    }, null, 2),
    {
      headers: {
        'Content-Type': 'application/activity+json',
        'Cache-Control': 'max-age=0, private, must-revalidate',
      }
    }
  );
};
