<script lang="ts">
  import type { PageData } from "./$types";
  import PDFThumbnail from "$lib/components/PDFThumbnail.svelte";
  import Header from "$lib/components/Header.svelte";
  import OpenGraph from "$lib/components/OpenGraph.svelte";
  let { data }: { data: PageData } = $props();
  import { page } from "$app/state";
</script>

<OpenGraph url={page.data.image} title={page.data.title} body={page.data.body}
></OpenGraph>

<Header url={page.data.image} title={page.data.title}></Header>

<main class="container px-3 lg:px-12 py-10 mx-auto">
  {#if data.pdfs.length === 0}
    <p class="text-gray-600">スライドがありません</p>
  {:else}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each data.pdfs as pdf}
        <div class="card shadow-md hover:shadow-lg transition-shadow">
          <a href="/slides/{pdf.id}">
            <PDFThumbnail url={pdf.url} />
          </a>
        </div>
      {/each}
    </div>
  {/if}
</main>
