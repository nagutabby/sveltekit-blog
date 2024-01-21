import type { PageServerLoad } from "./$types";
import { getAllContents } from "$lib/microcms";



export const load: PageServerLoad = async ({ url, parent }) => {
  const articleData = await getAllContents();
  let matchedTagName = "";
  let matchedContents: any = [];

  articleData.contents.forEach(content => {
    if (content.tags !== undefined) {
      const tags = content.tags.split(",");
      for (let i = 0; i < tags.length; i++) {
        if (url.pathname.split("/")[2] === tags[i].replaceAll(" ", "-").toLowerCase()) {
          if (matchedTagName === "") {
            matchedTagName = tags[i];
          }
          matchedContents.push(content);
          break;
        }
      }
    }
  });
  articleData.contents = matchedContents;
  const { titles } = await parent();
  const blogData: Blog = {
    image: {
      url: "https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/ca63de19468e45b2833ebf325dbfd56c/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png"
    },
    title: `${matchedTagName}タグが付けられた記事`,
    titles: titles!.concat(`${matchedTagName}タグが付けられた記事`),
    description: ""
  };
  const data = {
    ...articleData,
    ...blogData
  };
  return data;
};

export const prerender = false;
