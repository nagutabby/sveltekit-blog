import type { MicroCMSQueries } from "microcms-js-sdk";
import { getList } from "../lib/microcms";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ url }) => {
  const numberOfArticlesPerPage: number = 4;
  let startIndex: number;
  if (url.searchParams.get("page") !== null) {
    startIndex =
      numberOfArticlesPerPage *
      (Number(url.searchParams.get("page")) - 1);
  } else {
    startIndex = 0;
  }
  const pageQueries: MicroCMSQueries = {
    limit: numberOfArticlesPerPage,
    offset: startIndex
  }
  return await getList(pageQueries);
};

export const prerender = false;
