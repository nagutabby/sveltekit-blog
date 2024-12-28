<script lang="ts">
  import "../app.scss";
  import { page } from "$app/state";
  import Header from "$lib/components/Header.svelte";
  import Footer from "$lib/components/Footer.svelte";
  import OpenGraph from "$lib/components/OpenGraph.svelte";

  const { children } = $props()
</script>

{#if page.url.pathname === "/" || page.url.pathname === "/contact" || page.url.pathname === "/search"}
  <OpenGraph
    url={page.data.image.url}
    title={page.data.title}
    description={page.data.description}
  ></OpenGraph>

  <Header
    url={page.data.image.url}
    title={page.data.title}
    description={page.data.description}
  ></Header>

  <main class="container">
    {@render children()}
  </main>

  <Footer></Footer>
{:else}
  {@render children()}
{/if}
