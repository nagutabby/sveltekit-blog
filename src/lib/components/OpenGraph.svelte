<script lang="ts">
  const baseTitle = "nagutabbyの考え事";
  const baseURL = "https://blog.nagutabby.uk";
  const homeDescription = "学んだことをまとめるブログ";
  import { page } from "$app/stores";
  $: pathname = $page.url.pathname;

  export let subTitle: string;
  export let description: string;
  export let url: string;
</script>

{#if pathname === "/" || pathname.split("/")[1] === "pages" || pathname.split("/")[1] === "tags"}
  <title>{baseTitle}</title>
  <meta property="og:title" content={baseTitle} />
  <meta name="twitter:title" content={baseTitle} />
  <meta property="og:type" content="blog" />
  <meta name="description" content={homeDescription} />
  <meta property="og:description" content={homeDescription} />
  <meta
    property="og:image"
    content="https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/ca63de19468e45b2833ebf325dbfd56c/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png?fm=webp&fit=clip&w=1000"
  />
  <meta
    name="twitter:image"
    content="https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/ca63de19468e45b2833ebf325dbfd56c/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png?fm=webp&fit=clip&w=1000"
  />
{:else}
  <title>{subTitle} - {baseTitle}</title>
  <meta property="og:title" content="{subTitle} - {baseTitle}" />
  <meta name="twitter:title" content="{subTitle} - {baseTitle}" />
  <meta property="og:type" content="article" />
  <meta name="description" content={description} />
  <meta property="og:description" content={description} />
  <meta property="og:image" content="{url}?fm=webp&fit=clip&w=1000" />
  <meta name="twitter:image" content="{url}?fm=webp&fit=clip&w=1000" />
{/if}

<meta property="og:url" content={baseURL + pathname} />
<meta property="og:site_name" content={baseTitle} />
<meta property="og:locale" content="ja_JP" />
<meta name="twitter:card" content="summary" />
<meta name="twitter:site" content="@nagutabby" />
<link rel="canonical" href={baseURL + pathname} />
