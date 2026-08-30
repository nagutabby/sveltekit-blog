<script lang="ts">
  interface Props {
    totalPages: number;
    currentPage: number;
    hrefFor: (pageNumber: number) => string;
  }

  const { totalPages, currentPage, hrefFor }: Props = $props();

  const visiblePages = $derived(
    Array.from({ length: totalPages }, (_, i) => i + 1).filter(
      (pageNumber) => Math.abs(pageNumber - currentPage) <= 1,
    ),
  );
</script>

<ul class="flex justify-center items-center gap-x-5 p-4">
  <li>
    {#if currentPage === 1}
      <button class="btn btn-outline btn-secondary btn-circle text-lg" aria-disabled="true" disabled aria-label="最初のページに戻る">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-chevron-double-left" viewBox="0 0 16 16">
          <path fill-rule="evenodd" d="M8.354 1.646a.5.5 0 0 1 0 .708L2.707 8l5.647 5.646a.5.5 0 0 1-.708.708l-6-6a.5.5 0 0 1 0-.708l6-6a.5.5 0 0 1 .708 0z" />
          <path fill-rule="evenodd" d="M12.354 1.646a.5.5 0 0 1 0 .708L6.707 8l5.647 5.646a.5.5 0 0 1-.708.708l-6-6a.5.5 0 0 1 0-.708l6-6a.5.5 0 0 1 .708 0z" />
        </svg>
      </button>
    {:else}
      <a href={hrefFor(1)} role="button" class="btn btn-outline btn-secondary btn-circle text-lg" aria-label="最初のページに戻る">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-chevron-double-left" viewBox="0 0 16 16">
          <path fill-rule="evenodd" d="M8.354 1.646a.5.5 0 0 1 0 .708L2.707 8l5.647 5.646a.5.5 0 0 1-.708.708l-6-6a.5.5 0 0 1 0-.708l6-6a.5.5 0 0 1 .708 0z" />
          <path fill-rule="evenodd" d="M12.354 1.646a.5.5 0 0 1 0 .708L6.707 8l5.647 5.646a.5.5 0 0 1-.708.708l-6-6a.5.5 0 0 1 0-.708l6-6a.5.5 0 0 1 .708 0z" />
        </svg>
      </a>
    {/if}
  </li>
  {#each visiblePages as pageNumber (pageNumber)}
    <li>
      {#if pageNumber === currentPage}
        <button class="btn btn-outline btn-secondary btn-circle text-lg" aria-disabled="true" disabled>{pageNumber}</button>
      {:else}
        <a href={hrefFor(pageNumber)} role="button" class="btn btn-outline btn-secondary btn-circle text-lg"
          >{pageNumber}</a
        >
      {/if}
    </li>
  {/each}
  <li>
    {#if currentPage === totalPages}
      <button class="btn btn-outline btn-secondary btn-circle text-lg" aria-disabled="true" disabled aria-label="最後のページに進む">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-chevron-double-right" viewBox="0 0 16 16">
          <path fill-rule="evenodd" d="M3.646 1.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1 0 .708l-6 6a.5.5 0 0 1-.708-.708L9.293 8 3.646 2.354a.5.5 0 0 1 0-.708z" />
          <path fill-rule="evenodd" d="M7.646 1.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1 0 .708l-6 6a.5.5 0 0 1-.708-.708L13.293 8 7.646 2.354a.5.5 0 0 1 0-.708z" />
        </svg>
      </button>
    {:else}
      <a href={hrefFor(totalPages)} role="button" class="btn btn-outline btn-secondary btn-circle text-lg" aria-label="最後のページに進む">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-chevron-double-right" viewBox="0 0 16 16">
          <path fill-rule="evenodd" d="M3.646 1.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1 0 .708l-6 6a.5.5 0 0 1-.708-.708L9.293 8 3.646 2.354a.5.5 0 0 1 0-.708z" />
          <path fill-rule="evenodd" d="M7.646 1.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1 0 .708l-6 6a.5.5 0 0 1-.708-.708L13.293 8 7.646 2.354a.5.5 0 0 1 0-.708z" />
        </svg>
      </a>
    {/if}
  </li>
</ul>

<style lang="scss">

  li {
    list-style: none;
    margin: 1rem 0;
  }
</style>
