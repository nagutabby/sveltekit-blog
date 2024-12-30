import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params }) => {
  if (!params.name) {
    throw new Error('Name parameter is required');
  }

  return {
    pdfUrl: `/api/slides/${params.name}`
  };
}
