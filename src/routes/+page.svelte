<script lang="ts">
  import type { PageData } from "./$types";
  export let data: PageData;
  import Card from "$lib/components/Card.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import { page } from "$app/stores";
  function parseDate(date: string | Date) {
    date = new Date(date);
    const year = date.getFullYear();
    const month = ("00" + (date.getMonth() + 1)).slice(-2);
    const day = ("00" + date.getDate()).slice(-2);
    date = year + "-" + month + "-" + day;
    return date;
  }
</script>

{#each data.contents as { }, i}
  {#if i % 2 === 0}
    <div class="grid">
      <Card
        url={data.contents[i].id}
        image={data.contents[i].image.url}
        title={data.contents[i].title}
        date={parseDate(data.contents[i].createdAt)}
      />
      <Card
        url={data.contents[i + 1].id}
        image={data.contents[i + 1].image.url}
        title={data.contents[i + 1].title}
        date={parseDate(data.contents[i + 1].createdAt)}
      />
    </div>
  {/if}
{/each}

{#if $page.url.searchParams.get("page") !== null}
  <Pagination
    numberOfArticlesPerPage={data.contents.length}
    totalArticles={data.totalCount}
    currentPage={Number($page.url.searchParams.get("page"))}
  />
{:else}
  <Pagination
    numberOfArticlesPerPage={data.contents.length}
    totalArticles={data.totalCount}
  />
{/if}
