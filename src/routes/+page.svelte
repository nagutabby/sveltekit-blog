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

<OpenGraph
  url={page.data.image.url}
  title={page.data.title}
  body={page.data.body}
></OpenGraph>

<Header url={page.data.image.url} title={page.data.title}></Header>

<main class="container px-3 md:px-10 py-10 mx-auto">
  <div
    class="flex flex-col md:flex-row items-start md:relative gap-10 md:gap-0"
  >
    <div class="flex flex-wrap gap-5 justify-center w-full md:w-[63%]">
      {#each data.contents as article}
        <Card
          url={`articles/${article.id}`}
          image={article.image.url}
          title={article.title}
        />
      {/each}
    </div>
    <div class="w-full md:w-[33%] md:absolute md:top-0 md:bottom-0 md:right-0">
      <Timeline></Timeline>
    </div>
  </div>
  <Pagination totalArticles={data.totalCount} />
</main>
