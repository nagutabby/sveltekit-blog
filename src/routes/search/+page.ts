import type { MicroCMSQueries } from "microcms-js-sdk";
import type { PageLoad } from "./$types";
import type { ArticleInputData } from "../../types/blog";
import { getArticleList } from "$lib/microcms";

export const load: PageLoad = async ({ url }) => {
  const query = url.searchParams.get("q") || "";
  const page = Number(url.searchParams.get("page")) || 1;

  const limit = 9;
  const startIndex =
    limit *
    (page - 1);

  const pageQueries: MicroCMSQueries = {
    q: query,
    limit: limit,
    offset: startIndex
  };

  const articleData = await getArticleList(pageQueries);
  const blogData: ArticleInputData = {
    image: {
      url: "https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/ca63de19468e45b2833ebf325dbfd56c/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png"
    },
    title: `「${query}」を含む記事`,
    body: `「${query}」を含む記事の検索結果を表示しています`
  };
  const data = {
    ...articleData,
    ...blogData
  };
  return data;
};
