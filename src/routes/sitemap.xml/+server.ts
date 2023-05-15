import { getList } from "../../lib/microcms";

function create_entry(path: string, lastmod: string) {
  return `<url>
    <loc>${new URL(path, "https://blog.nagutabby.uk").href}</loc>
    ${lastmod ? `<lastmod>${lastmod}</lastmod>` : ''}
  </url>`;
}

export async function GET({ setHeaders }) {
  setHeaders({
    'Content-Type': 'application/xml'
  });

  const response = await getList();

  const posts = response.contents.map((post) => create_entry(post.id, post.revisedAt));

  const sitemap = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${posts.join('\n')}
</urlset>`;


  return new Response(sitemap);

}
