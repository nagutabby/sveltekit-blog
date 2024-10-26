import { getDetail } from "$lib/microcms";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params }) => {
  const articleData = await getDetail(params.slug);
  const data = {
    ...articleData,
  };
  return data;
};

export const prerender = true;
