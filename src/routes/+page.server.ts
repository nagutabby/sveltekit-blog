import type { MicroCMSQueries } from "microcms-js-sdk";
import { getList } from "../lib/microcms";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ url }) => {
  const numberOfArticlesPerPage = 6;
  let startIndex: number;
  let limit: number;
  if (url.searchParams.get("page") !== null) {
    limit = numberOfArticlesPerPage;
    startIndex =
      numberOfArticlesPerPage *
      (Number(url.searchParams.get("page")) - 1);
  } else {
    limit = numberOfArticlesPerPage;
    startIndex = 0;
  }
  let pageQueries: MicroCMSQueries;
  if (url.searchParams.get("tag") !== null) {
    pageQueries = {
      limit: limit,
      offset: startIndex,
      filters: `tags[contains]${url.searchParams.get("tag")}`
    }
  } else {
    pageQueries = {
      limit: limit,
      offset: startIndex,
      filters: "id[not_equals]about-me"
    }
  }
  let response = await getList(pageQueries)
  response.numberOfArticlesPerPage = numberOfArticlesPerPage;
  return response;
};

export const prerender = false;
