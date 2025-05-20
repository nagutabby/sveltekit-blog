<script lang="ts">
  import { onMount } from "svelte";
  import { browser } from "$app/environment";

  const { url }: { url: string } = $props();
  let canvas = $state<HTMLCanvasElement>();
  let loading = $state(true);
  let error = $state<string | null>(null);

  onMount(async () => {
    if (canvas && browser) {
      try {
        const pdfjs = await import("pdfjs-dist");
        const workerModule = await import(
          "pdfjs-dist/build/pdf.worker.mjs?url"
        );
        const workerUrl = workerModule.default;

        pdfjs.GlobalWorkerOptions.workerSrc = workerUrl;

        const pdf = await pdfjs.getDocument(url).promise;
        const page = await pdf.getPage(1);
        const viewport = page.getViewport({ scale: 1 });

        canvas.width = viewport.width;
        canvas.height = viewport.height;

        const context = canvas.getContext("2d");
        if (!context) throw new Error("Canvas context is null");

        await page.render({
          canvasContext: context,
          viewport: viewport,
        }).promise;

        loading = false;
      } catch (err) {
        console.error("PDF rendering error:", err);
        error = err instanceof Error ? err.message : "Failed to load PDF";
        loading = false;
      }
    }
  });
</script>

<div
  class="w-full aspect-video bg-base-100 rounded-lg overflow-hidden flex justify-center items-center"
>
  {#if loading}
    <span class="loading loading-spinner"></span>
  {:else if error}
    <span class="text-error">{error}</span>
  {/if}

  <canvas
    bind:this={canvas}
    class="w-full h-full object-contain block {loading || error
      ? 'hidden'
      : ''}"
  ></canvas>
</div>
