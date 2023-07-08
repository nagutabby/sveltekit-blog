<script lang="ts">
  import { page } from "$app/stores";
  export let totalArticles: number;
  export let currentPage: number | null = 1;
  export let numberOfArticlesPerPage: number;
  let pages: number;
  if (totalArticles % numberOfArticlesPerPage !== 0) {
    pages = Math.floor(totalArticles / numberOfArticlesPerPage) + 1;
  } else {
    pages = Math.floor(totalArticles / numberOfArticlesPerPage);
  }
</script>

<ul>
  {#each { length: pages } as _, i}
    <li>
      {#if i + 1 === currentPage}
        <button class="secondary outline" disabled>{i + 1}</button>
      {:else if $page.url.searchParams.get("tag") !== null}
        <a
          href="?tag={$page.url.searchParams.get('tag')}&page={i + 1}"
          role="button"
          class="secondary">{i + 1}</a
        >
      {:else}
        <a href="?page={i + 1}" role="button" class="secondary">{i + 1}</a>
      {/if}
    </li>
  {/each}
</ul>

<style>
  ul {
    margin: 1rem 0;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0 0.5rem;
  }
  li {
    list-style: none;
    height: 1vh;
  }
</style>
