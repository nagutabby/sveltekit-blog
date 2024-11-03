<script lang="ts">
  import { page } from "$app/stores";

  export let totalArticles: number;

  $: currentPage = Number($page.url.searchParams.get("page")) || 1;

  const numberOfArticlesPerPage = 9;
  let pages: number;
  $: {
    if (totalArticles % numberOfArticlesPerPage !== 0) {
      pages = Math.floor(totalArticles / numberOfArticlesPerPage) + 1;
    } else {
      pages = Math.floor(totalArticles / numberOfArticlesPerPage);
    }
  }

  const query = $page.url.searchParams.get("q");
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
