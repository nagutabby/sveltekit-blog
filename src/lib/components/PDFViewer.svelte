<script lang="ts">
  import { onMount } from "svelte";
  import { browser } from "$app/environment";

  let url: string = $props();
  let canvas = $state<HTMLCanvasElement>();
  let pdfDoc = $state<any>(null);
  let pageNum = $state(1);
  let totalPages = $state(0);
  let loading = $state(true);
  let error = $state<string | null>(null);
  let renderTask: any = null;

  let isFullscreenSupported = $state(false);
  let isFullscreen = $state(false);
  let viewerContainer = $state<HTMLDivElement>();

  const DISPLAY_SCALE = 1;
  const QUALITY_SCALE = 2;

  async function loadPdfLibrary() {
    if (browser) {
      const pdfjs = await import("pdfjs-dist");

      const workerModule = await import("pdfjs-dist/build/pdf.worker.mjs?url");
      pdfjs.GlobalWorkerOptions.workerSrc = workerModule.default;

      await loadPDF(pdfjs);
    }
  }

  async function cancelCurrentRender() {
    if (renderTask) {
      await renderTask.cancel();
      renderTask = null;
    }
  }

  $effect(() => {
    if (canvas && pdfDoc) {
      renderPage();
    }
  });

  async function renderPage() {
    if (!pdfDoc || !canvas) return;

    await cancelCurrentRender();

    const page = await pdfDoc.getPage(pageNum);

    const viewport = page.getViewport({ scale: DISPLAY_SCALE * QUALITY_SCALE });

    canvas.height = viewport.height;
    canvas.width = viewport.width;

    canvas.style.width = `${viewport.width / QUALITY_SCALE}px`;
    canvas.style.height = `${viewport.height / QUALITY_SCALE}px`;

    const context = canvas.getContext("2d");
    if (!context) return;

    context.clearRect(0, 0, canvas.width, canvas.height);

    const renderContext = {
      canvasContext: context,
      viewport: viewport,
      enableWebGL: true,
      renderInteractiveForms: true,
    };

    renderTask = page.render(renderContext);
    await renderTask.promise;

    renderTask = null;
  }

  async function loadPDF(pdfjs: any) {
    if (!pdfjs) return;

    loading = true;
    try {
      const loadingTask = pdfjs.getDocument(url);
      pdfDoc = await loadingTask.promise;
      totalPages = pdfDoc.numPages;
    } catch (err) {
      if (err instanceof Error) {
        error = `PDFの読み込みに失敗しました: ${err.message}`;
      }
    } finally {
      loading = false;
    }
  }

  function toggleFullscreen() {
    if (!viewerContainer) return;

    if (!document.fullscreenElement) {
      viewerContainer.requestFullscreen().catch((err) => {
        console.error(`フルスクリーンエラー: ${err.message}`);
      });
    } else {
      document.exitFullscreen();
    }
  }

  function handleFullscreenChange() {
    isFullscreen = !!document.fullscreenElement;
  }

  onMount(() => {
    loadPdfLibrary();

    if (browser) {
      isFullscreenSupported = !!document.documentElement.requestFullscreen;

      window.addEventListener("keydown", handleKeyDown);
      document.addEventListener("fullscreenchange", handleFullscreenChange);

      return () => {
        window.removeEventListener("keydown", handleKeyDown);
        document.removeEventListener("fullscreenchange", handleFullscreenChange);
      };
    }
  });

  function handleKeyDown(event: KeyboardEvent) {
    if (event.key === "ArrowRight" || event.key === "Right") {
      nextPage();
    } else if (event.key === "ArrowLeft" || event.key === "Left") {
      prevPage();
    }
  }

  async function nextPage() {
    if (pageNum < totalPages) {
      pageNum++;
      await renderPage();
    }
  }

  async function prevPage() {
    if (pageNum > 1) {
      pageNum--;
      await renderPage();
    }
  }
</script>

<div
  class="container flex flex-col items-center justify-center gap-y-8 px-3 lg:px-12 py-5"
>
  {#if loading}
    <span class="loading loading-spinner loading-lg"></span>
  {:else if error}
    <div class="alert alert-error">{error}</div>
  {:else}
    <div
      bind:this={viewerContainer}
      class="group relative max-w-full bg-base-100 flex items-center justify-center rounded-md overflow-hidden"
    >
      <canvas bind:this={canvas} class="block !w-full !h-full rounded-md"></canvas>

      {#if isFullscreenSupported}
        <button
          onclick={toggleFullscreen}
          class="absolute bottom-4 right-4 p-3 bg-black/40 hover:bg-black/60 text-white rounded-full shadow-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200"
          aria-label={isFullscreen ? "全画面表示を終了" : "全画面表示"}
        >
          {#if isFullscreen}
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" class="bi bi-fullscreen-exit" viewBox="0 0 16 16">  <path d="M5.5 0a.5.5 0 0 1 .5.5v4A1.5 1.5 0 0 1 4.5 6h-4a.5.5 0 0 1 0-1h4a.5.5 0 0 0 .5-.5v-4a.5.5 0 0 1 .5-.5m5 0a.5.5 0 0 1 .5.5v4a.5.5 0 0 0 .5.5h4a.5.5 0 0 1 0 1h-4A1.5 1.5 0 0 1 10 4.5v-4a.5.5 0 0 1 .5-.5M0 10.5a.5.5 0 0 1 .5-.5h4A1.5 1.5 0 0 1 6 11.5v4a.5.5 0 0 1-1 0v-4a.5.5 0 0 0-.5-.5h-4a.5.5 0 0 1-.5-.5m10 1a1.5 1.5 0 0 1 1.5-1.5h4a.5.5 0 0 1 0 1h-4a.5.5 0 0 0-.5.5v4a.5.5 0 0 1-1 0z"/></svg>
          {:else}
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" class="bi bi-arrows-fullscreen" viewBox="0 0 16 16">  <path fill-rule="evenodd" d="M5.828 10.172a.5.5 0 0 0-.707 0l-4.096 4.096V11.5a.5.5 0 0 0-1 0v3.975a.5.5 0 0 0 .5.5H4.5a.5.5 0 0 0 0-1H1.732l4.096-4.096a.5.5 0 0 0 0-.707m4.344 0a.5.5 0 0 1 .707 0l4.096 4.096V11.5a.5.5 0 1 1 1 0v3.975a.5.5 0 0 1-.5.5H11.5a.5.5 0 0 1 0-1h2.768l-4.096-4.096a.5.5 0 0 1 0-.707m0-4.344a.5.5 0 0 0 .707 0l4.096-4.096V4.5a.5.5 0 1 0 1 0V.525a.5.5 0 0 0-.5-.5H11.5a.5.5 0 0 0 0 1h2.768l-4.096 4.096a.5.5 0 0 0 0 .707m-4.344 0a.5.5 0 0 1-.707 0L1.025 1.732V4.5a.5.5 0 0 1-1 0V.525a.5.5 0 0 1 .5-.5H4.5a.5.5 0 0 1 0 1H1.732l4.096 4.096a.5.5 0 0 1 0 .707"/></svg>
          {/if}
        </button>
      {/if}
    </div>

    <div class="flex items-center gap-4 w-full justify-center">
      <button
        onclick={prevPage}
        disabled={pageNum <= 1}
        class="btn btn-neutral disabled:btn-disabled disabled:cursor-not-allowed block w-[30%] md:w-[15%]"
        aria-label="前のページに戻る"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          fill="currentColor"
          class="bi bi-arrow-left flex justify-center w-full"
          viewBox="0 0 16 16"
        >
          <path
            fill-rule="evenodd"
            d="M15 8a.5.5 0 0 0-.5-.5H2.707l3.147-3.146a.5.5 0 1 0-.708-.708l-4 4a.5.5 0 0 0 0 .708l4 4a.5.5 0 0 0 .708-.708L2.707 8.5H14.5A.5.5 0 0 0 15 8"
          />
        </svg>
      </button>
      <span class="text-lg">
        {pageNum} / {totalPages}
      </span>
      <button
        onclick={nextPage}
        disabled={pageNum >= totalPages}
        class="btn btn-neutral disabled:btn-disabled disabled:cursor-not-allowed block w-[30%] md:w-[15%]"
        aria-label="次のページに進む"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          fill="currentColor"
          class="bi bi-arrow-right flex justify-center w-full"
          viewBox="0 0 16 16"
        >
          <path
            fill-rule="evenodd"
            d="M1 8a.5.5 0 0 1 .5-.5h11.793l-3.147-3.146a.5.5 0 0 1 .708-.708l4 4a.5.5 0 0 1 0 .708l-4 4a.5.5 0 0 1-.708-.708L13.293 8.5H1.5A.5.5 0 0 1 1 8"
          />
        </svg>
      </button>
    </div>
  {/if}
</div>
