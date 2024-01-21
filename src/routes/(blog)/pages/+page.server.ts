import type { MicroCMSQueries } from "microcms-js-sdk";
import { getList } from "$lib/microcms";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ parent }) => {
  const limit = 1;
  const startIndex = 0;
  let pageQueries: MicroCMSQueries;
  pageQueries = {
    limit: limit,
    offset: startIndex,
    filters: "id[not_equals]about-me"
  }
  const articleData = await getList(pageQueries);
  const { titles } = await parent();
  const blogData: Blog = {
    image: {
      url: "https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/ca63de19468e45b2833ebf325dbfd56c/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png"
    },
    title: titles!.slice(-1)[0],
    description: ""
  };
  const data = {
    ...articleData,
    ...blogData
  }
  return data;
}

export const prerender = false;
