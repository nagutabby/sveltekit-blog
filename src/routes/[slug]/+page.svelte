<script lang="ts">
  import Breadcrumb from "$lib/components/Breadcrumb.svelte";
  import Toc from "$lib/components/Toc.svelte";
  import type { PageData } from "./$types";
  import { load } from "cheerio";
  import { onMount } from "svelte";
  import hljs from "highlight.js";

  export let data: PageData;

  let rawToc: string = "";

  onMount(async () => {
    document.querySelectorAll("pre code").forEach((el) => {
      hljs.highlightElement(el as HTMLElement);
    });
    const loadedBody = load(data.body);
    const headings = loadedBody("h1, h2, h3").toArray();
    const toc = headings.map((data: any) => ({
      text: data.children[0].data,
      id: data.attribs.id,
      name: data.name,
    }));
    let previousTag: string = "";
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
  });
</script>

<Breadcrumb title={data.title} />

<Toc toc={rawToc} />

<article>
  {@html data.body}
</article>

<style>
  @import url("https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/styles/a11y-light.min.css")
    screen and (prefers-color-scheme: light);
  @import url("https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/styles/a11y-dark.min.css")
    screen and (prefers-color-scheme: dark);
  article {
    margin-top: 0;
  }
  :global(pre > code) {
    font-size: 0.85rem;
  }
  :global(code) {
    padding: 0.15rem 0.4rem;
  }
  :global(.hljs) {
    background: none;
  }
  @media (prefers-color-scheme: light) {
    :global(code) {
      color: #545454;
    }
  }
  @media (prefers-color-scheme: dark) {
    :global(code) {
      color: #f8f8f2;
    }
  }
</style>
