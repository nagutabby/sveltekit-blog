import { getAllArticleContents } from "$lib/microcms";
import { generateDescriptionFromText } from "$lib/utils"
let latestDate: Date;

function create_entry(title: string, body: string, path: string, publishedAt: string, revisedAt: string) {
  const publishedDate = new Date(publishedAt).toISOString().substring(0, 10);
  if (latestDate === undefined) {
    latestDate = new Date(revisedAt);
  } else {
    if (latestDate < new Date(revisedAt)) {
      latestDate = new Date(revisedAt);
    }
  }

  return `<entry>
  <title>${title}</title>
  <summary>${generateDescriptionFromText(body)}</summary>
  <link href="${new URL(`/articles/${path}`, 'https://blog.nagutabby.uk').href}" />
  ${revisedAt ? `<updated>${revisedAt}</updated>` : ''}
  <published>${publishedAt}</published>
  <id>tag:blog.nagutabby.uk,${publishedDate}:/articles/${path}</id>
  </entry>`;
}

export async function GET({ setHeaders }) {
  setHeaders({
    'Content-Type': 'application/xml'
  });

  const response = await getAllArticleContents();

  const posts = response.contents.map((post) => create_entry(post.title, post.body, post.id, post.publishedAt, post.revisedAt));

  const atom = `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
<link rel="self" href="https://blog.nagutabby.uk/atom.xml" type="application/rss+xml" />
<title>nagutabbyの考え事</title>
<link href="https://blog.nagutabby.uk" />
<updated>${latestDate.toISOString()}</updated>
<author><name>nagutabby</name></author>
<id>tag:blog.nagutabby.uk,2023-01-01:/</id>
${posts.join('\n')}
</feed>`;


  return new Response(atom);

}
