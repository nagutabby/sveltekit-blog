<script>
  import Breadcrumb from "$lib/components/Breadcrumb.svelte";
  import { page } from "$app/stores";
</script>

{#key $page.url.pathname}
  <Breadcrumb titles={$page.data.titles} />
{/key}
<slot />
