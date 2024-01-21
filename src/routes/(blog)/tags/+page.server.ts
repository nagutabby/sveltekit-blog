import { getAllContents } from "$lib/microcms";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ parent }) => {
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
  const { titles } = await parent();
  const blogData: Blog = {
    tagSet: tagSet,
    image: {
      url: "https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/ca63de19468e45b2833ebf325dbfd56c/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png"
    },
    title: titles!.slice(-1)[0],
    description: ""
  }
  return blogData;
};

export const prerender = false;
