<script lang="ts">
  import type { PageData } from "./$types";
  import ListLayout from "$lib/components/ListLayout.svelte";

  const { data }: { data: PageData } = $props();

  function hrefFor(pageNumber: number) {
    return pageNumber === 1 ? "/" : `/page/${pageNumber}`;
  }
</script>

<ListLayout
  items={data.articles}
  urlPrefix="articles"
  pagination={data.pagination}
  {hrefFor}
/>
