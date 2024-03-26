<script lang="ts">
  import "../app.scss";
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
  ></OpenGraph>
</svelte:head>

<Header
  url={$page.data.image.url}
  title={$page.data.title}
  description={$page.data.description}
></Header>

<main class="container">
  <slot />
</main>

<Footer></Footer>

<style>
  .container {
    margin-top: 2rem;
  }
</style>
