<script lang="ts">
  import { page } from "$app/state";

  const { totalArticles }: { totalArticles: number } = $props();

  let currentPage = $derived(Number(page.url.searchParams.get("page")) || 1);

  const numberOfArticlesPerPage = 9;
  const pages = $derived(
    totalArticles % numberOfArticlesPerPage !== 0
    ? Math.floor(totalArticles / numberOfArticlesPerPage) + 1
    : Math.floor(totalArticles / numberOfArticlesPerPage)
  )

  const query = page.url.searchParams.get("q");

</script>

<ul>
  {#each { length: pages } as _, i}
    <li>
      {#if i + 1 === currentPage}
        <button class="outline" aria-disabled="true" disabled>{i + 1}</button>
      {:else if query}
        <a href="/search?q={query}&page={i + 1}" role="button" class="outline"
          >{i + 1}</a
        >
      {:else}
        <a href="/?page={i + 1}" role="button" class="outline">{i + 1}</a>
      {/if}
    </li>
  {/each}
</ul>

<style lang="scss">
  ul {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0 0.5rem;
    padding: 0;
  }
  li {
    list-style: none;
    margin: 1rem 0;
  }
  button {
    margin-bottom: 0;
  }
  a,
  button {
    border-radius: 10px;
    padding: 0.8rem 1.2rem;
    font-size: 1.2rem;
  }
</style>
