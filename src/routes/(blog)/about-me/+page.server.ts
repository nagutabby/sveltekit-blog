import type { PageServerLoad } from "./$types";
import { getDetail } from "$lib/microcms";

export const load: PageServerLoad = async ({ parent }) => {
  const articleData = await getDetail("about-me")
  const { titles } = await parent();
  const blogData: Blog = {
    titles: titles?.concat(articleData.title)
  }
  const data = {
    ...articleData,
    ...blogData
  }
  return data;
}

export const prerender = true;
