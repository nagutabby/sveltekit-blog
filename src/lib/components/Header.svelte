<script lang="ts">
  import { page } from "$app/stores";
  import cat_x1 from "$lib/assets/images/cat_1x.webp";
  import cat_x2 from "$lib/assets/images/cat_2x.webp";
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
      imagesrcset="{cat_x1} 600w, {cat_x2} 1200w"
    />
  {:else}
    <link
      rel="preload"
      as="image"
      imagesrcset="{url}?fit=clip&w=600 600w, {url}?fit=clip&w=1200 1200w"
    />
  {/if}
</svelte:head>

<header>
  {#if pathname === "/"}
    <img alt="" srcset="{cat_x1} 600w, {cat_x2} 1200w" />
  {:else}
    <img
      alt=""
      srcset="{url}?fit=clip&w=600 600w, {url}?fit=clip&w=1200 1200w"
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
