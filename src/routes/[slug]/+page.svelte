<script lang="ts">
  import Breadcrumb from "$lib/components/Breadcrumb.svelte";
  import Toc from "$lib/components/Toc.svelte";
  import type { PageData } from "./$types";
  import { load } from "cheerio";

  export let data: PageData;

  const loadedBody = load(data.body);
  const headings = loadedBody("h1, h2, h3").toArray();
  const toc = headings.map((data: any) => ({
    text: data.children[0].data,
    id: data.attribs.id,
    name: data.name,
  }));

  let rawToc: string = "";
  let previousTag: string = "";

  function getRawToc() {
    for (let i = 0; i < toc.length; i++) {
      if (previousTag === "") {
        rawToc += `<ul><li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
      } else if (toc[i].name === previousTag) {
        rawToc += `<li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
      } else if (previousTag === "h1") {
        if (toc[i].name === "h2" || "h3") {
          rawToc += `<ul><li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
        }
      } else if (previousTag === "h2") {
        if (toc[i].name === "h3") {
          rawToc += `<ul><li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
        } else if (toc[i].name === "h1") {
          rawToc += `</ul><li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
        }
      } else if (previousTag === "h3") {
        if (toc[i].name === "h2") {
          rawToc += `</ul><li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
        } else if (toc[i].name === "h1") {
          rawToc += `</ul></ul><li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
        }
      }
      previousTag = toc[i].name;
    }
    console.log(toc);
    return rawToc;
  }
</script>

<Breadcrumb title={data.title} />

<Toc toc={getRawToc()} />

<article>
  {@html data.body}
</article>

<style>
  article {
    margin-top: 0;
  }
</style>
