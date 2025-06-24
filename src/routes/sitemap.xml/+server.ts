import { getAllRawData } from '$lib/utils.js';

function create_entry(path: string, lastmod: Date) {
  return `<url>
    <loc>${new URL(`articles/${path}`, "https://blog.nagutabby.uk").href}</loc>
    ${lastmod ? `<lastmod>${lastmod.toISOString()}</lastmod>` : ''}
  </url>`;
}

export async function GET({ setHeaders }) {
  setHeaders({
    'Content-Type': 'application/xml'
  });

  const allArticles = await getAllRawData("articles")

  const posts = allArticles.map((post) => create_entry(post.id, post.updatedAt));

  const sitemap = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${posts.join('\n')}
</urlset>`;


  return new Response(sitemap);

}
