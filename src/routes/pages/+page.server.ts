import type { MicroCMSQueries } from "microcms-js-sdk";
import { getList } from "$lib/microcms";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ url }) => {
  const limit = 1;
  const startIndex = 0;
  let pageQueries: MicroCMSQueries;
  pageQueries = {
    limit: limit,
    offset: startIndex,
    filters: "id[not_equals]about-me"
  }
  let response = await getList(pageQueries)
  return response;
};

export const prerender = true;
