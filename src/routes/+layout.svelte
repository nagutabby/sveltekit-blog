<script lang="ts">
  import "$lib/assets/stylesheets/app.scss";
  import { page } from "$app/state";
  import Header from "$lib/components/Header.svelte";
  import Footer from "$lib/components/Footer.svelte";
  import OpenGraph from "$lib/components/OpenGraph.svelte";

  const { children } = $props();
</script>

{#if page.route.id?.includes("/slides")}
  {@render children()}
{:else}
  <OpenGraph
    url={page.data.image.url}
    title={page.data.title}
    body={page.data.body}
  ></OpenGraph>

  <Header url={page.data.image.url} title={page.data.title}></Header>

  <main class="container px-3 md:px-10 py-10 mx-auto">
    {@render children()}
  </main>
{/if}

<Footer></Footer>
