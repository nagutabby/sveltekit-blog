<script lang="ts">
  import type { PageData } from "./$types";
  import Pagination from "$lib/components/Pagination.svelte";
  import { page } from "$app/stores";

  export let data: PageData;
</script>

<Pagination
  totalArticles={data.totalCount}
  currentPage={Number($page.url.pathname.split("/")[2])}
/>
