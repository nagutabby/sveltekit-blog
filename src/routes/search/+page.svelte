<script lang="ts">
  import type { PageData } from "./$types";
  import Card from "$lib/components/Card.svelte";
  import Pagination from "$lib/components/Pagination.svelte";

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
        />
        <Card
          url={data.contents[i - 1].id}
          image={data.contents[i - 1].image.url}
          title={data.contents[i - 1].title}
        />
        <Card
          url={data.contents[i].id}
          image={data.contents[i].image.url}
          title={data.contents[i].title}
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
        />
        <Card
          url={data.contents[i - 1].id}
          image={data.contents[i - 1].image.url}
          title={data.contents[i - 1].title}
        />
        <Card
          url={data.contents[i].id}
          image={data.contents[i].image.url}
          title={data.contents[i].title}
        />
      </div>
    {:else if i + 1 == data.contents.length}
      <div class="grid">
        <Card
          url={data.contents[i - 1].id}
          image={data.contents[i - 1].image.url}
          title={data.contents[i - 1].title}
        />
        <Card
          url={data.contents[i].id}
          image={data.contents[i].image.url}
          title={data.contents[i].title}
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
        />
        <Card
          url={data.contents[i - 1].id}
          image={data.contents[i - 1].image.url}
          title={data.contents[i - 1].title}
        />
        <Card
          url={data.contents[i].id}
          image={data.contents[i].image.url}
          title={data.contents[i].title}
        />
      </div>
    {:else if i + 1 == data.contents.length}
      <div class="grid">
        <Card
          url={data.contents[i].id}
          image={data.contents[i].image.url}
          title={data.contents[i].title}
        />
      </div>
    {/if}
  {/each}
{/if}

<Pagination totalArticles={data.totalCount} />
