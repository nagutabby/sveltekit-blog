<script lang="ts">
  import { page } from "$app/state";

  const { totalPages }: { totalPages: number } = $props();

  let currentPage = $derived(Number(page.url.searchParams.get("page")) || 1);

  const query = page.url.searchParams.get("q");

</script>

<ul class="flex justify-center items-center gap-x-5 p-4">
  {#each { length: totalPages } as _, i}
    <li>
      {#if i + 1 === currentPage}
        <button class="btn btn-outline btn-secondary btn-circle text-lg" aria-disabled="true" disabled>{i + 1}</button>
      {:else if query}
        <a href="/search?q={query}&page={i + 1}" role="button" class="btn btn-outline btn-secondary btn-circle text-lg"
          >{i + 1}</a
        >
      {:else}
        <a href="/?page={i + 1}" role="button" class="btn btn-outline btn-secondary btn-circle text-lg">{i + 1}</a>
      {/if}
    </li>
  {/each}
</ul>

<style lang="scss">

  li {
    list-style: none;
    margin: 1rem 0;
  }
</style>
