import { getDetail } from "$lib/microcms";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params }) => {
  const articleData = await getDetail(params.name);
  const data = {
    ...articleData,
  };
  return data;
};

