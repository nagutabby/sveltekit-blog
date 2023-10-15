<script lang="ts">
  import { page } from "$app/stores";
  $: pathname = $page.url.pathname;

  export let url: string;
  export let title: string;
  export let description: string;
</script>

<svelte:head>
  {#if pathname === "/"}
    <link
      rel="preload"
      as="image"
      imagesrcset="https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/ca63de19468e45b2833ebf325dbfd56c/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png?fm=webp&fit=clip&w=600 600w, https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/ca63de19468e45b2833ebf325dbfd56c/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png?fm=webp&fit=clip&w=1200 1200w"
    />
  {:else}
    <link
      rel="preload"
      as="image"
      imagesrcset="{url}?fm=webp&fit=clip&w=600 600w, {url}?fm=webp&fit=clip&w=1200 1200w"
    />
  {/if}
</svelte:head>

<header>
  {#if pathname === "/"}
    <img
      alt=""
      srcset="https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/ca63de19468e45b2833ebf325dbfd56c/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png?fm=webp&fit=clip&w=600 600w, https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/ca63de19468e45b2833ebf325dbfd56c/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png?fm=webp&fit=clip&w=1200 1200w"
    />
  {:else}
    <img
      alt=""
      srcset="{url}?fm=webp&fit=clip&w=600 600w, {url}?fm=webp&fit=clip&w=1200 1200w"
    />
  {/if}
  <div class="title container">
    <hgroup>
      <h1>
        {#if pathname === "/"}
          nagutabbyの考え事
        {:else}
          {title}
        {/if}
      </h1>
      <h2>
        {#if pathname === "/"}
          学んだことをまとめるブログ
        {:else}
          {description}
        {/if}
      </h2>
    </hgroup>
  </div>
</header>

<style>
  header {
    position: relative;
    padding: 0 !important;
    & img {
      width: 100%;
      height: 50vh;
      object-fit: cover;
      filter: brightness(70%);
      background-color: white;
    }
    & .title {
      position: absolute;
      left: 50%;
      top: 50%;
      transform: translate(-50%, -50%);
      & hgroup {
        display: flex;
        flex-flow: column;
        gap: 1rem 0;
      }
      & hgroup > h1 {
        font-size: 2.3rem;
      }
      & hgroup > h2 {
        font-size: 1.3rem;
      }
    }
  }
  @media (prefers-color-scheme: light) {
    h1,
    h2 {
      color: hsl(205, 20%, 94%);
    }
  }
  @media (prefers-color-scheme: dark) {
    h2 {
      color: var(--h1-color);
    }
  }
</style>
