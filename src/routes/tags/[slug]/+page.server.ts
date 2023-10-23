import { getAllContents } from "$lib/microcms";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ url }) => {
  const response = await getAllContents();
  let matchedContents: any = []
  response.contents.forEach(content => {
    if (content.tags !== undefined) {
      const tags = content.tags.split(",")
      for (let i = 0; i < tags.length; i++) {
        if (url.pathname.split("/")[2] === tags[i].replaceAll(" ", "-").toLowerCase()) {
          matchedContents.push(content);
          break;
        }
      }
    }
  })
  response.contents = matchedContents
  response.totalCount = matchedContents.length
  return response;
};

export const prerender = true;
