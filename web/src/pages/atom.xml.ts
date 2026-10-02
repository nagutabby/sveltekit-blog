import type { APIRoute } from 'astro';
import { getAllHTMLData } from '$lib/server/content';
import { createAtomFeed } from '$lib/server/feeds';
import type { Article, Review } from '$lib/types/blog';

export const prerender = true;

export const GET: APIRoute = async () => {
  const articles = await getAllHTMLData('articles') as Article[];
  const reviews = await getAllHTMLData('reviews') as Review[];
  return new Response(createAtomFeed(articles, reviews), {
    headers: { 'Content-Type': 'application/xml' }
  });
};
