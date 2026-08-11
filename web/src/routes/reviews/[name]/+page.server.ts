import type { Review } from "$lib/types/blog";
import { getAllRawData, getHTMLData } from "$lib/server/content";
import type { PageServerLoad } from "./$types";

export async function entries() {
  const allReviews = await getAllRawData("reviews") as Review[];
  return allReviews.map((review) => ({ name: review.id }));
}

export const load: PageServerLoad = async ({ params }) => {
  const review = await getHTMLData(params.name, "reviews") as Review;
  return review;
};

