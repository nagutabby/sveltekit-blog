<script lang="ts">
  import Breadcrumb from "$lib/components/Breadcrumb.svelte";
  import Toc from "$lib/components/Toc.svelte";
  import type { PageData } from "./$types";
  import { load } from "cheerio";
  import { onMount } from "svelte";
  import TagList from "$lib/components/TagList.svelte";

  export let data: PageData;

  let rawToc: string = "";

  onMount(async () => {
    const loadedBody = load(data.body);
    let headings: any = loadedBody("h1, h2, h3").toArray();
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
    let activeHeading: any;
    const updateElements = (entries: any) => {
      entries.forEach((entry: any) => {
        if (entry.isIntersecting) {
          if (activeHeading !== undefined && activeHeading !== null) {
            activeHeading.classList.remove("active");
          }
          activeHeading = document.querySelector(
            `a[href="#${entry.target.id}"]`
          );

          if (activeHeading !== null) {
            activeHeading.classList.add("active");
          }
        }
      });
    };
    headings = document.querySelectorAll("h1, h2, h3");
    const options = {
      threshold: 1,
    };
    const observer = new IntersectionObserver(updateElements, options);
    headings.forEach((heading: any) => {
      observer.observe(heading);
    });
  });
</script>

<Breadcrumb title={data.title} />

<TagList tags={data.tags} />

<div class="row article-group">
  <div class="col-12 col-lg-4">
    <div class="sticky-top toc">
      <Toc toc={rawToc} />
    </div>
  </div>
  <div class="col-12 col-lg-8">
    <article>
      {@html data.body}
    </article>
  </div>
</div>

<style>
  .article-group {
    display: flex;
    flex-direction: row-reverse;
  }
  .sticky-top {
    position: sticky;
    top: 0;
  }
  :global(pre > code) {
    font-size: 0.85rem;
  }
  @media (prefers-color-scheme: light) {
    :global(code) {
      color: #545454;
    }
    :root {
      --primary-hover: #0d47a1;
    }
    :global(.active) {
      background: linear-gradient(transparent 60%, #B2EBF2 60%);
    }
  }
  @media (prefers-color-scheme: dark) {
    :global(code) {
      color: #f8f8f2;
    }
    :root {
      --primary-hover: #29b6f6;
    }
    :global(.active) {
      background: linear-gradient(transparent 60%, #455a64 60%);
    }
  }
  :global(a) {
    text-decoration: none;
  }
  :root {
    --primary: #1e88e5 !important;
  }
</style>
