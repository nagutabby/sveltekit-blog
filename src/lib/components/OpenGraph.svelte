<script lang="ts">
  const baseTitle = "nagutabbyの考え事";
  const baseURL = "https://blog.nagutabby.uk";
  const homeDescription = "学んだことをまとめるブログ";
  import { page } from "$app/stores";

  export let subTitle: string;
  export let description: string;
  export let url: string;
</script>

<meta name="viewport" content="width=device-width" />
<meta property="og:url" content={baseURL + $page.url.pathname} />
<meta property="og:site_name" content={baseTitle} />
<meta property="og:locale" content="ja_JP" />
<meta name="twitter:card" content="summary" />
<link rel="canonical" href={baseURL + $page.url.pathname} />
<meta name="theme-color" content="black" />
<meta name="apple-mobile-web-app-status-bar-style" content="black" />
<meta name="mobile-web-app-capable" content="yes" />
<meta name="apple-mobile-web-app-capable" content="yes" />

{#if subTitle === ""}
  <title>{baseTitle}</title>
  <meta property="og:title" content={baseTitle} />
  <meta name="twitter:title" content={baseTitle} />
  <meta property="og:type" content="website" />
{:else}
  <title>{subTitle} - {baseTitle}</title>
  <meta property="og:title" content="{subTitle} - {baseTitle}" />
  <meta name="twitter:title" content="{subTitle} - {baseTitle}" />
  <meta property="og:type" content="article" />
{/if}

{#if description === ""}
  <meta name="description" content={homeDescription} />
  <meta property="og:description" content={homeDescription} />
{:else}
  <meta name="description" content={description} />
  <meta property="og:description" content={description} />
{/if}

{#if url === ""}
  <meta
    property="og:image"
    content="https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/ca63de19468e45b2833ebf325dbfd56c/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png?fm=webp&fit=clip&w=1000"
  />
  <meta
    name="twitter:image"
    content="https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/ca63de19468e45b2833ebf325dbfd56c/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png?fm=webp&fit=clip&w=1000"
  />
{:else}
  <meta property="og:image" content="{url}?fm=webp&fit=clip&w=1000" />
  <meta name="twitter:image" content="{url}?fm=webp&fit=clip&w=1000" />
{/if}
