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

{#if data.contents.length % 3 == 0}
  {#each data.contents as { }, i}
    {#if (i + 1) % 3 == 0}
      <div class="grid">
        <Card
          url={data.contents[i - 2].id}
          image={data.contents[i - 2].image.url}
          title={data.contents[i - 2].title}
          date={parseDate(data.contents[i - 1].createdAt)}
        />
        <Card
          url={data.contents[i - 1].id}
          image={data.contents[i - 1].image.url}
          title={data.contents[i - 1].title}
          date={parseDate(data.contents[i].createdAt)}
        />
        <Card
          url={data.contents[i].id}
          image={data.contents[i].image.url}
          title={data.contents[i].title}
          date={parseDate(data.contents[i].createdAt)}
        />
      </div>
    {/if}
  {/each}
{:else if data.contents.length % 3 == 2}
  {#each data.contents as { }, i}
    {#if (i + 1) % 3 == 0}
      <div class="grid">
        <Card
          url={data.contents[i - 2].id}
          image={data.contents[i - 2].image.url}
          title={data.contents[i - 2].title}
          date={parseDate(data.contents[i - 1].createdAt)}
        />
        <Card
          url={data.contents[i - 1].id}
          image={data.contents[i - 1].image.url}
          title={data.contents[i - 1].title}
          date={parseDate(data.contents[i].createdAt)}
        />
        <Card
          url={data.contents[i].id}
          image={data.contents[i].image.url}
          title={data.contents[i].title}
          date={parseDate(data.contents[i].createdAt)}
        />
      </div>
    {:else if i + 1 == data.contents.length}
      <div class="grid">
        <Card
          url={data.contents[i - 1].id}
          image={data.contents[i - 1].image.url}
          title={data.contents[i - 1].title}
          date={parseDate(data.contents[i].createdAt)}
        />
        <Card
          url={data.contents[i].id}
          image={data.contents[i].image.url}
          title={data.contents[i].title}
          date={parseDate(data.contents[i].createdAt)}
        />
      </div>
    {/if}
  {/each}
{:else if data.contents.length % 3 == 1}
  {#each data.contents as { }, i}
    {#if (i + 1) % 3 == 0}
      <div class="grid">
        <Card
          url={data.contents[i - 2].id}
          image={data.contents[i - 2].image.url}
          title={data.contents[i - 2].title}
          date={parseDate(data.contents[i - 1].createdAt)}
        />
        <Card
          url={data.contents[i - 1].id}
          image={data.contents[i - 1].image.url}
          title={data.contents[i - 1].title}
          date={parseDate(data.contents[i].createdAt)}
        />
        <Card
          url={data.contents[i].id}
          image={data.contents[i].image.url}
          title={data.contents[i].title}
          date={parseDate(data.contents[i].createdAt)}
        />
      </div>
    {:else if i + 1 == data.contents.length}
      <div class="grid">
        <Card
          url={data.contents[i].id}
          image={data.contents[i].image.url}
          title={data.contents[i].title}
          date={parseDate(data.contents[i].createdAt)}
        />
      </div>
    {/if}
  {/each}
{/if}

<Pagination
  totalArticles={data.totalCount}
  currentPage={Number($page.url.pathname.split("/")[2])}
/>
