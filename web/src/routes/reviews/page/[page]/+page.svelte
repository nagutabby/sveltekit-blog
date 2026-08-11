<script lang="ts">
  import type { PageData } from "./$types";
  import ListLayout from "$lib/components/ListLayout.svelte";

  const { data }: { data: PageData } = $props();

  function hrefFor(pageNumber: number) {
    return pageNumber === 1 ? "/reviews" : `/reviews/page/${pageNumber}`;
  }
</script>

<ListLayout
  items={data.reviews}
  urlPrefix="reviews"
  pagination={data.pagination}
  {hrefFor}
/>
