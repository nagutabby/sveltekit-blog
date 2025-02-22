import type { MicroCMSQueries } from "microcms-js-sdk";
import type { PageLoad } from "./$types";
import type { BlogInputData } from "../types/blog";
import { getList } from "$lib/microcms";

export const load: PageLoad = async ({ url }) => {
  const page = Number(url.searchParams.get("page")) || 1;
  const limit = 10;
  const startIndex = limit * (page - 1);
  let pageQueries: MicroCMSQueries;
  pageQueries = {
    limit: limit,
    offset: startIndex
  };
  const articleData = await getList(pageQueries);

  const blogData: BlogInputData = {
    image: {
      url: "https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/ca63de19468e45b2833ebf325dbfd56c/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png"
    },
    title: "nagutabbyの考え事",
    body: "学んだことをまとめるブログ"
  };
  const data = {
    ...articleData,
    ...blogData
  };
  return data;
};
