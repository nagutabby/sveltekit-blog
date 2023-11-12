import { getAllContents } from "../../lib/microcms";

let latestDate: Date;

function create_entry(title: string, description: string, path: string, createdAt: string, revisedAt: string) {
  const createdDate = new Date(createdAt).toISOString().substring(0, 10);
  if (latestDate === undefined) {
    latestDate = new Date(revisedAt);
  } else {
    if (latestDate < new Date(revisedAt)) {
      latestDate = new Date(revisedAt);
    }
  }

  return `<entry>
  <title>${title}</title>
  <summary>${description}</summary>
  <link href="${new URL(path, 'https://blog.nagutabby.uk').href}" />
  ${revisedAt ? `<updated>${revisedAt}</updated>` : ''}
  <id>tag:blog.nagutabby.uk,${createdDate}:/${path}</id>
  </entry>`;
}

export async function GET({ setHeaders }) {
  setHeaders({
    'Content-Type': 'application/xml'
  });

  const response = await getAllContents();

  const posts = response.contents.map((post) => create_entry(post.title, post.description, post.id, post.createdAt, post.revisedAt));

  const atom = `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xml:lang="ja">
<title>nagutabbyの考え事</title>
<link href="https://blog.nagutabby.uk" />
<updated>${latestDate.toISOString()}</updated>
<author><name>nagutabby</name></author>
<id>tag:blog.nagutabby.uk,2023-01-01:/</id>
${posts.join('\n')}
</feed>`;


  return new Response(atom);

}

export const prerender = true;
