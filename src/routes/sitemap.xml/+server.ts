import type { Article, Review } from '$lib/types/blog.js';
import { getAllRawData } from '$lib/utils.js';

function createEntry(path: string, lastmod: Date, type: "articles" | "reviews") {
  return `<url>
    <loc>${new URL(`${type}/${path}`, "https://blog.nagutabby.uk").href}</loc>
    ${lastmod ? `<lastmod>${lastmod.toISOString()}</lastmod>` : ''}
  </url>`;
}

export async function GET({ setHeaders }) {
  setHeaders({
    'Content-Type': 'application/xml'
  });

  const allArticles = await getAllRawData("articles") as Article[]
  const allReviews = await getAllRawData("reviews") as Review[]

  const articlePosts = allArticles.map((post) => createEntry(post.id, post.updatedAt, "articles"));
  const reviewPosts = allReviews.map((post) => createEntry(post.id, post.updatedAt, "reviews"));

  const posts = [...articlePosts, ...reviewPosts]
  
  const sitemap = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${posts.join('\n')}
</urlset>`;


  return new Response(sitemap);

}
