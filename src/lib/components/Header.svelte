<script lang="ts">
  import HeaderIcons from "./HeaderIcons.svelte";

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
    <div class="container header-icons-container">
      <HeaderIcons />
    </div>
    <div class="container header-group">
      <h1>
        {title}
      </h1>
      <h2>
        {description}
      </h2>
    </div>
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
      flex-flow: column;
      height: 100%;
      width: 100%;
      background: rgba(0, 0, 0, 0.5);
      backdrop-filter: blur(3px);
      justify-content: center;
      align-items: center;
      position: absolute;
      left: 50%;
      top: 50%;
      transform: translate(-50%, -50%);
      display: flex;
      & .header-group {
        display: flex;
        flex-flow: column;
        gap: 1rem 0;
        justify-content: center;
        width: 100%;
        position: relative;

        & h1 {
          font-size: 2.2rem;
        }
        & h2 {
          font-size: 1.4rem;
        }
      }
    }
    & .header-icons-container {
      position: fixed;
      top: 2rem;
      z-index: 1;
    };
  }
</style>
