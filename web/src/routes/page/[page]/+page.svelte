<script lang="ts">
  import type { PageData } from "./$types";
  import ArticlesPage from "$lib/components/ArticlesPage.svelte";

  const { data }: { data: PageData } = $props();
</script>

<ArticlesPage articles={data.articles} pagination={data.pagination} />
