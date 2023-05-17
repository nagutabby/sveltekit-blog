<script lang="ts">
  import Breadcrumb from "$lib/components/Breadcrumb.svelte";
  import Toc from "$lib/components/Toc.svelte";
  import type { PageData } from "./$types";
  export let data: PageData;

  import { load } from "cheerio";

  const loadedBody = load(data.body);
  const headings = loadedBody("h1, h2, h3").toArray();
  const toc = headings.map((data) => ({
    text: (data.children[0] as any).data,
    id: data.attribs.id,
    name: data.name,
  }));
</script>

<Breadcrumb title={data.title} />

<Toc {toc} />

<article>
  {@html data.body}
</article>

<style>
  article {
    margin-top: 0;
  }
</style>
