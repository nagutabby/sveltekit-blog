<script lang="ts">
  import type { PageData } from "./$types";
  import Card from "$lib/components/Card.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import Timeline from "$lib/components/Timeline.svelte";
  import OpenGraph from "$lib/components/OpenGraph.svelte";
  import Header from "$lib/components/Header.svelte";
  import { page } from "$app/state";

  const { data }: { data: PageData } = $props();
</script>

<OpenGraph url={page.data.image} title={page.data.title} body={page.data.body}
></OpenGraph>

<Header url={page.data.image} title={page.data.title}></Header>

<main class="container px-3 md:px-10 py-10 mx-auto">
  <div
    class="flex flex-col lg:flex-row items-start lg:relative gap-10 lg:gap-0"
  >
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-5 w-full lg:w-[63%] h-fit">
      {#if data.articles}
        {#each data.articles as article}
          <Card
            id={article.id}
            url={`articles/${article.id}`}
            image={article.image}
            title={article.title}
          />
        {/each}
      {/if}
    </div>
    <div class="w-full lg:w-[33%] lg:absolute lg:top-0 lg:bottom-0 lg:right-0">
      <Timeline></Timeline>
    </div>
  </div>
  {#if data.pagination}
    <Pagination totalPages={data.pagination.totalPages} />
  {/if}
</main>
