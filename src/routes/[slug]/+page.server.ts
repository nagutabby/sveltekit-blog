import { getDetail } from "$lib/microcms";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params, parent }) => {
  const getArticle = async () => {
    const articleData = await getDetail(params.slug);
    const data = {
      ...articleData,
    };
    return data;
  };

  return { streamed: { article: getArticle() } };
};

export const prerender = false;
