<script lang="ts">
  import type { PageData } from "./$types";
  import ListLayout from "$lib/components/ListLayout.svelte";

  const { data }: { data: PageData } = $props();
</script>

<ListLayout
  items={data.reviews}
  urlPrefix="reviews"
  pagination={data.pagination}
/>
