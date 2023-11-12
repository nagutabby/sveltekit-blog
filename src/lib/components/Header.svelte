<script lang="ts">
  export let url: string;
  export let title: string;
  export let description: string;
</script>

<svelte:head>
  <link
    rel="preload"
    as="image"
    imagesrcset="{url}?fm=webp&fit=clip&w=500 500w, {url}?fm=webp&fit=clip&w=1000 1000w"
  />
</svelte:head>

<header>
  <img
    alt=""
    srcset="{url}?fm=webp&fit=clip&w=500 500w, {url}?fm=webp&fit=clip&w=1000 1000w"
  />
  <div class="title">
    <hgroup class="container">
      <h1>
        {title}
      </h1>
      <h2>
        {description}
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
      background-color: white;
    }
    & .title {
      height: 100%;
      width: 100%;
      background: rgba(0, 0, 0, 0.5);
      backdrop-filter: blur(3px);
      position: absolute;
      left: 50%;
      top: 50%;
      transform: translate(-50%, -50%);
      display: flex;
      & hgroup {
        display: flex;
        flex-flow: column;
        gap: 1rem 0;
        justify-content: center;
        width: 100%;
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
