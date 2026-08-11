<script lang="ts">
  import PDFViewer from "$lib/components/PDFViewer.svelte";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();
</script>

<main class="container px-3 lg:px-12 py-10 mx-auto">
  <PDFViewer url={data.url} />
</main>
