import type { APIRoute } from 'astro';

export const prerender = true;

export const GET: APIRoute = () => new Response('ok');
