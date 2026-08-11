<script lang="ts">
  interface Props {
    totalPages: number;
    currentPage: number;
    hrefFor: (pageNumber: number) => string;
  }

  const { totalPages, currentPage, hrefFor }: Props = $props();
</script>

<ul class="flex justify-center items-center gap-x-5 p-4">
  {#each { length: totalPages } as _, i}
    <li>
      {#if i + 1 === currentPage}
        <button class="btn btn-outline btn-secondary btn-circle text-lg" aria-disabled="true" disabled>{i + 1}</button>
      {:else}
        <a href={hrefFor(i + 1)} role="button" class="btn btn-outline btn-secondary btn-circle text-lg"
          >{i + 1}</a
        >
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
