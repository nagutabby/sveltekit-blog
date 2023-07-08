import type { PageServerLoad } from "./$types";
import { getDetail } from "../../lib/microcms";

export const load: PageServerLoad = async () => {
  let response = await getDetail("about-me")
  return response;
}
