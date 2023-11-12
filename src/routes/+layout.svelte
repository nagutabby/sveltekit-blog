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
  @media (prefers-color-scheme: light) {
    :global(code) {
      color: #545454;
    }
    :root {
      --primary-hover: #0d47a1;
    }
    :global(.active) {
      background: linear-gradient(transparent 60%, #b2ebf2 60%);
    }
  }
  @media (prefers-color-scheme: dark) {
    :global(code) {
      color: #f8f8f2;
    }
    :root {
      --primary-hover: #29b6f6;
    }
    :global(.active) {
      background: linear-gradient(transparent 60%, #455a64 60%);
    }
  }
  :global(a) {
    text-decoration: none;
  }
  :root {
    --primary: #1e88e5 !important;
  }
</style>
