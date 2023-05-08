import { getDetail } from "../../lib/microcms";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params }) => {
  return await getDetail(
    params.slug
  );
};

export const prerender = true;
