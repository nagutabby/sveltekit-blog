<script lang="ts">
  import type { PageData } from "./$types";
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
  export let data: PageData;
</script>

{#each data.contents as { }, i}
  {#if (i + 1) % 2 === 0}
    <div class="grid">
      <Card
        url={data.contents[i - 1].id}
        image={data.contents[i - 1].image.url}
        title={data.contents[i - 1].title}
        date={parseDate(data.contents[i - 1].createdAt)}
      />
      <Card
        url={data.contents[i].id}
        image={data.contents[i].image.url}
        title={data.contents[i].title}
        date={parseDate(data.contents[i].createdAt)}
      />
    </div>
  {:else if i + 1 === data.contents.length}
    <Card
      url={data.contents[i].id}
      image={data.contents[i].image.url}
      title={data.contents[i].title}
      date={parseDate(data.contents[i].createdAt)}
    />
  {/if}
{/each}

{#if $page.url.searchParams.get("page") !== null}
  <Pagination
    numberOfArticlesPerPage={data.numberOfArticlesPerPage}
    totalArticles={data.totalCount}
    currentPage={Number($page.url.searchParams.get("page"))}
  />
{:else}
  <Pagination
    numberOfArticlesPerPage={data.numberOfArticlesPerPage}
    totalArticles={data.totalCount}
  />
{/if}
