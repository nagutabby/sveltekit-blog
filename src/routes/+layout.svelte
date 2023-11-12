<script lang="ts">
  import Header from "$lib/components/Header.svelte";
  import { page } from "$app/stores";
  import Footer from "$lib/components/Footer.svelte";
  import OpenGraph from "$lib/components/OpenGraph.svelte";
</script>

<svelte:head>
  <OpenGraph
    url={$page.data.image.url}
    title={$page.data.title}
    description={$page.data.description}
  />
</svelte:head>

<Header
  url={$page.data.image.url}
  title={$page.data.title}
  description={$page.data.description}
/>

<main class="container">
  <slot />
</main>

<Footer />

<style>
  :global(body) {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }
  main {
    flex-grow: 1;
  }
</style>
