import type { APIRoute } from 'astro';
import { getAllRawData } from '$lib/server/content';
import { createSitemap } from '$lib/server/feeds';
import type { Article, Review } from '$lib/types/blog';

export const prerender = true;

export const GET: APIRoute = async () => {
  const articles = await getAllRawData('articles') as Article[];
  const reviews = await getAllRawData('reviews') as Review[];
  return new Response(createSitemap(articles, reviews), {
    headers: { 'Content-Type': 'application/xml' }
  });
};
