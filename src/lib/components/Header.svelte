<script lang="ts">
  import Navigation from "./Navigation.svelte";

  type Props = {
    url: string;
    title: string;
  };
  const { url, title }: Props = $props();
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
    class="w-full h-[40vh] lg:h-[50vh]"
  />
  <div class="title w-full h-full left-[50%] top-[50%] absolute flex flex-col">
    <div class="container px-3 lg:px-12 py-3 relative h-full flex flex-col gap-y-[30%]">
      <div class="asbolute w-full">
        <Navigation />
      </div>
      <div class="flex flex-col max-w-full">
        <h1 class="text-3xl sm:text-4xl lg:text-5xl font-semibold">
          {title}
        </h1>
      </div>
    </div>
  </div>
</header>

<style lang="scss">
  header {
    padding: 0 !important;
    & img {
      object-fit: cover;
      background-color: white;
    }
    & .title {
      background: rgba(0, 0, 0, 0.5);
      backdrop-filter: blur(3px);
      justify-content: center;
      align-items: center;
      transform: translate(-50%, -50%);
    }
  }
</style>
