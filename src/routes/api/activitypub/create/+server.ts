import prisma from '$lib/prisma';
import type { RequestHandler } from '@sveltejs/kit';

export const POST: RequestHandler = async ({ request, fetch }) => {
  const data = await request.json();

  if (data.id) {
    const response = await fetch(`/articles/${data.id}/create`);
    console.log(response.json())
  }
  return new Response(null, { status: 200 });
}
