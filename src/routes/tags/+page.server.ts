import { getAllContents } from "$lib/microcms";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ url }) => {
  const response = await getAllContents();
  let tagSet: Set<string> = new Set()
  response.contents.forEach(content => {
    if (content.tags !== undefined) {
      const tags = content.tags.split(",")
      for (let i = 0; i < tags.length; i++) {
        tagSet.add(tags[i]);
      }
    }
  })
  return {tagSet};
};

export const prerender = true;
