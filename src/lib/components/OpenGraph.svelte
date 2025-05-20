<script lang="ts">
  import { page } from "$app/state";
  import { generateDescriptionFromText } from "$lib/utils";

  const baseURL = "https://blog.nagutabby.uk";
  const defaultTitle = "nagutabbyの考え事";
  const defaultDescription = "学んだことをまとめるブログ";

  type Props = {
    title: string;
    body: string;
    url: string;
  };
  const { title, body, url }: Props = $props();

  let description = $derived(generateDescriptionFromText(body));
</script>

<svelte:head>
  <meta property="og:url" content={baseURL + page.url.pathname} />
  <meta property="og:site_name" content={defaultTitle} />
  <meta property="og:locale" content="ja_JP" />
  <meta name="twitter:card" content="summary" />
  <link rel="canonical" href={baseURL + page.url.pathname} />
  <meta name="theme-color" content="black" />
  <meta name="apple-mobile-web-app-status-bar-style" content="black" />
  <meta name="mobile-web-app-capable" content="yes" />
  <meta name="apple-mobile-web-app-capable" content="yes" />

  {#if title === defaultTitle}
    <title>{defaultTitle}</title>
    <meta property="og:tile" content={defaultTitle} />
    <meta name="twitter:title" content={defaultTitle} />
    <meta property="og:type" content="website" />
  {:else}
    <title>{title} - {defaultTitle}</title>
    <meta property="og:title" content="{title} - {defaultTitle}" />
    <meta name="twitter:title" content="{title} - {defaultTitle}" />
    <meta property="og:type" content="article" />
  {/if}

  {#if description === defaultDescription}
    <meta name="description" content={defaultDescription} />
    <meta property="og:description" content={defaultDescription} />
  {:else}
    <meta name="description" content={description} />
    <meta property="og:description" content={description} />
  {/if}

  <meta property="og:image" content={url} />
  <meta name="twitter:image" content={url} />
</svelte:head>
