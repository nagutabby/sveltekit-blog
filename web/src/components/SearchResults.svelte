<script lang="ts">
  import { onMount } from 'svelte';
  import Card from '$lib/components/Card.svelte';
  import Pagination from '$lib/components/Pagination.svelte';
  import type { SearchableArticle } from '$lib/types/search';

  let { articles }: { articles: SearchableArticle[] } = $props();
  const perPage = 10;
  let locationSearch = $state('');

  onMount(() => {
    const sync = () => { locationSearch = window.location.search; };
    sync();
    window.addEventListener('popstate', sync);
    document.addEventListener('astro:page-load', sync);
    return () => {
      window.removeEventListener('popstate', sync);
      document.removeEventListener('astro:page-load', sync);
    };
  });

  const query = $derived(new URLSearchParams(locationSearch).get('q') ?? '');
  const requestedPage = $derived(Number(new URLSearchParams(locationSearch).get('page')) || 1);
  const filteredArticles = $derived(query
    ? articles.filter((article) => article.rawBody.toLowerCase().includes(query.toLowerCase()))
    : articles);
  const totalPages = $derived(Math.max(Math.ceil(filteredArticles.length / perPage), 1));
  const currentPage = $derived(Math.min(Math.max(requestedPage, 1), totalPages));
  const paginatedArticles = $derived(filteredArticles.slice((currentPage - 1) * perPage, currentPage * perPage));
  const hrefFor = (pageNumber: number) => `/search?q=${encodeURIComponent(query)}&page=${pageNumber}`;

  $effect(() => {
    if (typeof document === 'undefined') return;
    const title = query ? `「${query}」を含む記事` : '記事を検索';
    const description = query ? `「${query}」を含む記事の検索結果を表示しています` : '記事を検索できます';
    document.title = `${title} - nagutabbyの考え事`;
    document.querySelector('meta[name="description"]')?.setAttribute('content', description);
    document.querySelector('meta[property="og:title"]')?.setAttribute('content', `${title} - nagutabbyの考え事`);
    document.querySelector('meta[property="og:description"]')?.setAttribute('content', description);
    document.querySelector('meta[name="twitter:title"]')?.setAttribute('content', `${title} - nagutabbyの考え事`);
  });
</script>

<div class="grid grid-cols-1 sm:grid-cols-2 gap-5 w-full lg:w-[63%] h-fit">
  {#each paginatedArticles as article}
    <Card id={article.id} url={`articles/${article.id}`} image={article.image} title={article.title} />
  {/each}
</div>
{#if filteredArticles.length > 0}
  <Pagination {totalPages} {currentPage} {hrefFor} />
{/if}
