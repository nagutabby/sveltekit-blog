<!-- routes/slides/+page.svelte -->
<script lang="ts">
  import type { PageData } from "./$types";
  import PDFThumbnail from "$lib/components/PDFThumbnail.svelte";
  import Header from "$lib/components/Header.svelte";
  import OpenGraph from "$lib/components/OpenGraph.svelte";
  let { data }: { data: PageData } = $props();

  const url = "https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/b5c72b7577a5452d8b902f9909b755ea/Microsoft-Fluentui-Emoji-3d-Package-3d.1024.png"
  const title = "スライド一覧";
  const description = "LT会で使ったスライドをまとめています"
</script>

<OpenGraph
  {url}
  {title}
  body={description}
></OpenGraph>

<Header
  {url}
  {title}
></Header>

<main class="container px-3 lg:px-12 py-10 mx-auto">
    {#if data.pdfUrls.length === 0}
      <p class="text-gray-600">スライドがありません</p>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {#each data.pdfUrls as pdfUrl}
          <div class="card shadow-md hover:shadow-lg transition-shadow">
            <a href="/slides/{pdfUrl}">
              <PDFThumbnail {pdfUrl} />
            </a>
          </div>
        {/each}
      </div>
    {/if}
</main>
