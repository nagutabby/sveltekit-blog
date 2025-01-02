<script lang="ts">
  import Navigation from "./Navigation.svelte";

  type Props = {
    url: string;
    title: string;
    description: string;
  };
  const { url, title, description }: Props = $props();
</script>

<svelte:head>
  <link
    rel="preload"
    as="image"
    imagesrcset="{url}?fm=webp&fit=clip&w=500 500w, {url}?fm=webp&fit=clip&w=1000 1000w"
  />
</svelte:head>

<header class="relative">
  <img
    alt=""
    srcset="{url}?fm=webp&fit=clip&w=500 500w, {url}?fm=webp&fit=clip&w=1000 1000w"
  />
  <div class="title w-full h-full">
    <div class="container px-3 lg:px-12 py-3 relative h-full flex flex-col gap-y-[25%]">
      <div class="asbolute w-full">
        <Navigation />
      </div>
      <div class="flex flex-col max-w-full prose">
        <h1 class="text-4xl">
          {title}
        </h1>
        <h2 class="text-2xl my-3">
          {description}
        </h2>
      </div>
    </div>
  </div>
</header>

<style lang="scss">
  header {
    position: relative;
    padding: 0 !important;
    & img {
      width: 100%;
      height: 60vh;
      object-fit: cover;
      background-color: white;
    }
    & .title {
      flex-flow: column;
      background: rgba(0, 0, 0, 0.5);
      backdrop-filter: blur(3px);
      justify-content: center;
      align-items: center;
      position: absolute;
      left: 50%;
      top: 50%;
      transform: translate(-50%, -50%);
      display: flex;
    }
  }
</style>
