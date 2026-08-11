<script lang="ts">
  import Card from "$lib/components/Card.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import Timeline from "$lib/components/Timeline.svelte";
  import OpenGraph from "$lib/components/OpenGraph.svelte";
  import Header from "$lib/components/Header.svelte";
  import { page } from "$app/state";

  interface Item {
    id: string;
    image: string;
    title: string;
  }

  interface Props {
    items: Item[] | undefined;
    urlPrefix: string;
    pagination?: {
      totalPages: number;
    };
  }

  const { items, urlPrefix, pagination }: Props = $props();
</script>

<OpenGraph url={page.data.image} title={page.data.title} body={page.data.body} />

<Header url={page.data.image} title={page.data.title} />

<main class="container px-3 md:px-10 py-10 mx-auto">
  <div class="flex flex-col lg:flex-row items-start lg:relative gap-10 lg:gap-0">
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-5 w-full lg:w-[63%] h-fit">
      {#if items}
        {#each items as item}
          <Card
            id={item.id}
            url={`${urlPrefix}/${item.id}`}
            image={item.image}
            title={item.title}
          />
        {/each}
      {/if}
    </div>
    <div class="w-full lg:w-[33%] lg:absolute lg:top-0 lg:bottom-0 lg:right-0">
      <Timeline />
    </div>
  </div>
  {#if pagination}
    <Pagination totalPages={pagination.totalPages} />
  {/if}
</main>
