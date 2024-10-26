import { getDetail } from "$lib/microcms";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params }) => {
  const articleData = await getDetail(params.slug);
  const data = {
    ...articleData,
  };
  return data;
};

export const prerender = false;
