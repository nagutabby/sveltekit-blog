<script lang="ts">
  const baseURL = "https://blog.nagutabby.uk";
  const defaultTitle = "nagutabbyの考え事";
  const defaultDescription = "学んだことをまとめるブログ";
  import { page } from "$app/stores";

  export let title: string;
  export let description: string;
  export let url: string;
</script>

<meta name="viewport" content="width=device-width" />
<meta property="og:url" content={baseURL + $page.url.pathname} />
<meta property="og:site_name" content={defaultTitle} />
<meta property="og:locale" content="ja_JP" />
<meta name="twitter:card" content="summary" />
<link rel="canonical" href={baseURL + $page.url.pathname} />
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

<meta property="og:image" content="{url}?fm=webp&fit=clip&w=1000" />
<meta name="twitter:image" content="{url}?fm=webp&fit=clip&w=1000" />
