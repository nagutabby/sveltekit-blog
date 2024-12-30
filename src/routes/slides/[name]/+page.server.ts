export async function load({ params }) {
  if (!params.name) {
    throw new Error('Name parameter is required');
  }

  return {
    pdfUrl: `/api/slides/${params.name}`
  };
}
