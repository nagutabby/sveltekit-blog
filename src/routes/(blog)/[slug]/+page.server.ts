import { getDetail } from "$lib/microcms";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params, parent }) => {
  const articleData = await getDetail(params.slug);
  const { titles } = await parent();
  const blogData: Blog = {
    titles: titles?.concat(articleData.title)
  }
  const data = {
    ...articleData,
    ...blogData
  }
  return data
};

export const prerender = true;
