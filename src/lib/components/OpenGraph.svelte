<script lang="ts">
  const baseTitle = "nagutabbyの考え事";
  const baseURL = "https://blog.nagutabby.uk";
  const homeDescription = "学んだことをまとめるブログ";
  import { page } from "$app/stores";
  import cat from "$lib/assets/images/cat.jpg";
  $: pathname = $page.url.pathname;

  export let subTitle: string;
  export let description: string;
  export let url: string;
</script>

{#if pathname === "/"}
  <title>{baseTitle}</title>
  <meta property="og:title" content={baseTitle} />
  <meta name="twitter:title" content={baseTitle} />
  <meta property="og:type" content="blog" />
  <meta name="description" content={homeDescription} />
  <meta property="og:description" content={homeDescription} />
  <meta property="og:image" content={cat} />
  <meta name="twitter:image" content={cat} />
{:else}
  <title>{subTitle} - {baseTitle}</title>
  <meta property="og:title" content="{subTitle} - {baseTitle}" />
  <meta name="twitter:title" content="{subTitle} - {baseTitle}" />
  <meta property="og:type" content="article" />
  <meta name="description" content={description} />
  <meta property="og:description" content={description} />
  <meta property="og:image" content="{url}?fit=clip&w=1200" />
  <meta name="twitter:image" content="{url}?fit=clip&w=1200" />
{/if}

<meta property="og:url" content={baseURL + pathname} />
<meta property="og:site_name" content={baseTitle} />
<meta property="og:locale" content="ja_JP" />
<meta name="twitter:card" content="summary" />
<meta name="twitter:site" content="@nagutabby" />
<link rel="canonical" href={baseURL + pathname} />
