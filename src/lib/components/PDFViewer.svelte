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
  let pdfjs: any = null;

  const DISPLAY_SCALE = 1;
  const QUALITY_SCALE = 2;

  async function loadPdfLibrary() {
    if (browser) {
      const pdfModule = await import("pdfjs-dist");
      pdfjs = pdfModule;

      const workerModule = await import("pdfjs-dist/build/pdf.worker.mjs?url");
      pdfjs.GlobalWorkerOptions.workerSrc = workerModule.default;

      await loadPDF();
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

  async function loadPDF() {
    if (!pdfjs) return;

    loading = true;
    try {
      const loadingTask = pdfjs.getDocument(url);
      pdfDoc = await loadingTask.promise;
      totalPages = pdfDoc.numPages;
    } catch (err) {
      error = `PDFの読み込みに失敗しました: ${err.message}`;
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    loadPdfLibrary();

    if (browser) {
      window.addEventListener("keydown", handleKeyDown);

      return () => {
        window.removeEventListener("keydown", handleKeyDown);
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
    <canvas
      bind:this={canvas}
      class="block !w-full !h-full rounded-md"
    ></canvas>
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
