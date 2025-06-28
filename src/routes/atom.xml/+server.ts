import type { Article, Review } from "$lib/types/blog.js";
import { generateDescriptionFromText } from "$lib/utils";
import { getAllHTMLData } from "$lib/utils";

export async function GET({ setHeaders }) {
  setHeaders({
    'Content-Type': 'application/xml'
  });

  const allArticles = await getAllHTMLData("articles") as Article[];
  const allReviews = await getAllHTMLData("reviews") as Review[];

  let latestDate: Date | undefined;

  function createEntry(title: string, body: string, path: string, publishedAt: Date, updatedAt: Date, type: "articles" | "reviews") {
    const publishedDate = new Date(publishedAt).toISOString().substring(0, 10);

    if (latestDate === undefined) {
      latestDate = new Date(updatedAt);
    } else if (latestDate < new Date(updatedAt)) {
      latestDate = new Date(updatedAt);
    }

    const formattedPublishedAt = new Date(publishedAt).toISOString();
    const formattedUpdatedAt = new Date(updatedAt).toISOString();

    return `<entry>
  <title>${title}</title>
  <summary type="text"><![CDATA[${generateDescriptionFromText(body)}]]></summary>
  <link href="${new URL(`/${type}/${path}`, 'https://blog.nagutabby.uk').href}" rel="alternate" />
  <updated>${formattedUpdatedAt}</updated>
  <published>${formattedPublishedAt}</published>
  <id>tag:blog.nagutabby.uk,${publishedDate}:/${type}/${path}</id>
  </entry>`;
  }

  const articlePosts = allArticles.map((post) => createEntry(post.title, post.body, post.id, post.publishedAt, post.updatedAt, "articles"));
  const reviewPosts = allReviews.map((post) => createEntry(post.title, post.body, post.id, post.publishedAt, post.updatedAt, "reviews"));

  const posts = [...articlePosts, ...reviewPosts]
  
  const atom = `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
<link rel="self" href="https://blog.nagutabby.uk/atom.xml" type="application/rss+xml" />
<title>nagutabbyの考え事</title>
<link href="https://blog.nagutabby.uk" />
<updated>${latestDate ? latestDate.toISOString() : new Date().toISOString()}</updated>
<author><name>nagutabby</name></author>
<id>tag:blog.nagutabby.uk,2023-01-01:/</id>
${posts.join('\n')}
</feed>`;


  return new Response(atom);
}
