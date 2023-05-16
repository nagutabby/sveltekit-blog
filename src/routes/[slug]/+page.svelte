<script lang="ts">
  import Breadcrumb from "$lib/components/Breadcrumb.svelte";
  import type { PageData } from "./$types";
  export let data: PageData;
</script>

<Breadcrumb title={data.title} />

<article>
  {@html data.body}
</article>

<style>
  article {
    margin-top: 0;
  }
</style>
