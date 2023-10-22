<script lang="ts">
  import Breadcrumb from "$lib/components/Breadcrumb.svelte";
  import Toc from "$lib/components/Toc.svelte";
  import type { PageData } from "./$types";
  import { onMount } from "svelte";
  import TagList from "$lib/components/TagList.svelte";

  export let data: PageData;

  let rawToc: string = "";

  onMount(async () => {
    let headings: any = Array.from(document.getElementById("content")!.querySelectorAll("h1, h2, h3"));
    const toc = headings.map((heading: Element) => ({
      text: heading.textContent,
      id: heading.getAttribute("id"),
      name: heading.tagName,
    }));
    let previousTag: string = "";
    for (let i = 0; i < toc.length; i++) {
      if (previousTag === "") {
        rawToc += `<ul><li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
      } else if (toc[i].name === previousTag) {
        rawToc += `<li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
      } else if (previousTag === "H1") {
        if (toc[i].name === "H2" || "H3") {
          rawToc += `<ul><li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
        }
      } else if (previousTag === "H2") {
        if (toc[i].name === "H3") {
          rawToc += `<ul><li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
        } else if (toc[i].name === "H1") {
          rawToc += `</ul><li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
        }
      } else if (previousTag === "H3") {
        if (toc[i].name === "H2") {
          rawToc += `</ul><li><a href="#${toc[i].id}">${toc[i].text}</a></li>`;
        } else if (toc[i].name === "H1") {
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
      rootMargin: "0% 0% -100% 0%",
      threshold: 0,
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
  <div class="col-12 col-lg-4 sticky-top">
    <div class="toc">
      <Toc toc={rawToc} />
    </div>
  </div>
  <div class="col-12 col-lg-8" id="content">
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
    align-self: flex-start;
  }
  .toc {
    & article {
      margin: 0;
    }
    & details[open] > ul {
      max-height: 20vh;
      overflow-y: auto;
    }
  }
  :global(pre > code) {
    font-size: 0.85rem;
  }
  @media screen and (min-width: 992px) {
    :global(.toc > article) {
      margin: var(--block-spacing-vertical) 0 !important;
    }
    :global(.toc > article > details[open] > ul) {
      max-height: 40vh !important;
    }
  }
  @media (prefers-color-scheme: light) {
    :global(code) {
      color: #545454;
    }
    :root {
      --primary-hover: #0d47a1;
    }
    :global(.active) {
      background: linear-gradient(transparent 60%, #b2ebf2 60%);
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
