<script lang="ts">
  import type { PageData } from "./$types";
  import Card from "$lib/components/Card.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import Timeline from "$lib/components/Timeline.svelte";
  import Header from "$lib/components/Header.svelte";
  import OpenGraph from "$lib/components/OpenGraph.svelte";
  import { page } from "$app/state";
  import { browser } from "$app/environment";

  const { data }: { data: PageData } = $props();

  const PER_PAGE = 10;

  // page.url.searchParams cannot be read while prerendering (SSG), so this
  // only resolves to the real query/page once hydrated in the browser.
  const query = $derived(browser ? page.url.searchParams.get("q") ?? "" : "");
  const requestedPage = $derived(browser ? Number(page.url.searchParams.get("page")) || 1 : 1);

  const filteredArticles = $derived(
    query
      ? data.articles.filter((article) =>
          article.rawBody.toLowerCase().includes(query.toLowerCase())
        )
      : data.articles
  );

  const totalPages = $derived(Math.max(Math.ceil(filteredArticles.length / PER_PAGE), 1));
  const currentPage = $derived(Math.min(Math.max(requestedPage, 1), totalPages));

  const paginatedArticles = $derived(
    filteredArticles.slice((currentPage - 1) * PER_PAGE, currentPage * PER_PAGE)
  );

  const title = $derived(query ? `「${query}」を含む記事` : "記事を検索");
  const body = $derived(
    query ? `「${query}」を含む記事の検索結果を表示しています` : "記事を検索できます"
  );

  function hrefFor(pageNumber: number) {
    return `/search?q=${encodeURIComponent(query)}&page=${pageNumber}`;
  }
</script>

<OpenGraph url={data.image} {title} {body} />

<Header url={data.image} {title} />

<main class="container px-3 md:px-10 py-10 mx-auto">
  <div
    class="flex flex-col md:flex-row items-start md:relative gap-10 md:gap-0"
  >
    <div class="flex flex-wrap gap-5 justify-center w-full md:w-[63%]">
      {#each paginatedArticles as article}
        <Card
          id={article.id}
          url={`articles/${article.id}`}
          image={article.image}
          title={article.title}
        />
      {/each}
    </div>
    <div class="w-full md:w-[33%] md:absolute md:top-0 md:bottom-0 md:right-0">
      <Timeline></Timeline>
    </div>
  </div>
  {#if filteredArticles.length > 0}
    <Pagination totalPages={totalPages} {currentPage} {hrefFor} />
  {/if}
</main>
