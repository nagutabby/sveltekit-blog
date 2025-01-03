<script lang="ts">
  import { onMount } from "svelte";
  import * as pdfjs from "pdfjs-dist";
  // @ts-ignore
  import workerUrl from "pdfjs-dist/build/pdf.worker.mjs?url";
  import { browser } from "$app/environment";

  const { pdfUrl }: string = $props();
  let canvas = $state<HTMLCanvasElement>();
  let loading = $state(true);
  let error = $state<string | null>();

  onMount(async () => {
    if (canvas) {
      try {
        if (browser) {
          pdfjs.GlobalWorkerOptions.workerSrc = workerUrl;
        }
        const response = await fetch(`/api/slides/${pdfUrl}`);
        if (!response.ok) throw new Error("PDF fetch failed");

        const pdfData = await response.arrayBuffer();

        const pdf = await pdfjs.getDocument({ data: pdfData }).promise;
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
  {/if}

  <canvas
    bind:this={canvas}
    class="w-full h-full object-contain block {loading || error
      ? 'hidden'
      : ''}"
  ></canvas>
</div>
