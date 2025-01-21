<script lang="ts">
  import "../app.scss";
  import { page } from "$app/state";
  import Header from "$lib/components/Header.svelte";
  import Footer from "$lib/components/Footer.svelte";
  import OpenGraph from "$lib/components/OpenGraph.svelte";

  const { children } = $props();
</script>

{#if page.data.image && page.data.title && page.data.description}
  <OpenGraph
    url={page.data.image.url}
    title={page.data.title}
    description={page.data.description}
  ></OpenGraph>

  <Header url={page.data.image.url} title={page.data.title}></Header>

  <main class="container px-3 md:px-10 py-10 mx-auto">
    {@render children()}
  </main>
{:else}
  {@render children()}
{/if}

<Footer></Footer>
