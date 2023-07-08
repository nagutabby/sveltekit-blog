<script lang="ts">
  import Header from "$lib/components/Header.svelte";
  import { page } from "$app/stores";
  import Footer from "$lib/components/Footer.svelte";
  import OpenGraph from "$lib/components/OpenGraph.svelte";
  $: pathname = $page.url.pathname;
</script>

<svelte:head>
  {#if pathname === "/"}
    <OpenGraph url="" subTitle="" description="" />
  {:else}
    <OpenGraph
      url={$page.data.image.url}
      subTitle={$page.data.title}
      description={$page.data.description}
    />
  {/if}
</svelte:head>

{#if pathname === "/"}
  <Header url="" title="" description="" />
{:else}
  <Header
    url={$page.data.image.url}
    title={$page.data.title}
    description={$page.data.description}
  />
{/if}

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
