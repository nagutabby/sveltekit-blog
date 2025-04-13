import { getReviewDetail } from "$lib/microcms";
import type { PageServerLoad } from "./$types";
import { GOOGLE_CLOUD_API_KEY } from "$env/static/private";
import { getBook } from "$lib/utils";

export const load: PageServerLoad = async ({ params }) => {
  const reviewData = await getReviewDetail(params.name);

  const apiKey = GOOGLE_CLOUD_API_KEY;
  const bookData = await getBook(reviewData.isbn_13, apiKey);

  if (bookData.thumbnailUrl.startsWith('http:')) {
    bookData.thumbnailUrl = bookData.thumbnailUrl.replace('http:', 'https:');
  }

  const data = {
    review: reviewData,
    book: bookData
  };
  return data;
};
