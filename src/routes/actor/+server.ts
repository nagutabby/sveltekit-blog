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
      "summary": '<p>ブログ記事を投稿するBotアカウントです。</p><p>運用者: <a href="https://mastodon.social/@nagutabby" target="_blank">@nagutabby</a></p>',
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
      },
      "icon": {
        "type": "Image",
        "mediaType": "image/png",
        "url": "https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/ca63de19468e45b2833ebf325dbfd56c/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png?fit=clip&w=500"
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
