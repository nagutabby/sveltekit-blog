import type { Review } from "$lib/types/blog";
import { getHTMLData } from "$lib/utils";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params }) => {
  const review = await getHTMLData(params.name, "reviews") as Review;
  return review;
};

