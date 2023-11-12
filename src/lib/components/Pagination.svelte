<script lang="ts">
  import { page } from "$app/stores";
  export let totalArticles: number;
  export let currentPage: number | null = 1;
  const numberOfArticlesPerPage = 9;
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
        <button class="outline" aria-disabled="true" disabled>{i + 1}</button>
      {:else if $page.url.pathname.split("/")[1] === "tags"}
        <a
          href="/tags/{$page.url.pathname.split('/')[2]}/pages/{i + 1}"
          role="button"
          class="outline">{i + 1}</a
        >
      {:else if $page.url.pathname.split("/")[1] === "pages" || $page.url.pathname.split("/")[1] === ""}
        <a href="/pages/{i + 1}" role="button" class="outline">{i + 1}</a>
      {/if}
    </li>
  {/each}
</ul>

<style>
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
