<script lang="ts">
  import type { PageData } from "./$types";
  import { onMount } from "svelte";
  import ShareButton from "$lib/components/ShareButton.svelte";
  import Date from "$lib/components/Date.svelte";
  import Toc from "$lib/components/Toc.svelte";
  import { page } from "$app/state";
  import Header from "$lib/components/Header.svelte";
  import OpenGraph from "$lib/components/OpenGraph.svelte";
  import StarRating from "$lib/components/StarRating.svelte";
  import BookInfo from "$lib/components/BookInfo.svelte";

  const { data }: { data: PageData } = $props();

  let headingList = $state();
  let isLoading = $state(true);

  onMount(async () => {
    headingList = Array.from(
      document.getElementById("content")!.querySelectorAll("h1, h2, h3"),
    );
    isLoading = false;
  });
</script>

<OpenGraph
  url={page.data.image}
  title={page.data.title}
  body={page.data.body}
></OpenGraph>

<Header url={page.data.image} title={page.data.title}
></Header>

<main class="container px-3 md:px-10 py-10 mx-auto">
  <div class="flex lg:flex-row items-start justify-between flex-col-reverse">
    <div
      class="flex flex-wrap gap-5 justify-center w-full lg:w-[63%]"
      id="content"
    >
      <article class="prose max-w-full">
        <div class="flex flex-col gap-y-4">
          <Date
            publishedAt={data.publishedAt}
            updatedAt={data.updatedAt}
          />

          <div class="flex items-center gap-2">
            <StarRating rating={data.rating}></StarRating>
            <span class="font-bold text-lg">{data.rating}/5</span>
          </div>
          <div>
            {@html data.body}
          </div>
          <BookInfo
            title={data.title}
            description={data.description}
            jpECode={data.jp_e_code}
          />
        </div>
      </article>
    </div>
    <div class="w-full lg:w-[33%] prose lg:top-0 lg:sticky py-5">
      <div class="toc">
        {#if isLoading}
          <div class="flex justify-center">
            <span class="loading loading-spinner loading-lg"></span>
          </div>
        {:else}
          <Toc {headingList} />
        {/if}
      </div>
      <hr />
      <ShareButton />
    </div>
  </div>
</main>
