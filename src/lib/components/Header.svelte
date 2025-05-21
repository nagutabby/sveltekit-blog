<script lang="ts">
  import Navigation from "./Navigation.svelte";
  import {getWebpPath} from "$lib/utils"

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
    href={getWebpPath(url)}
  />
</svelte:head>

<header class="relative">
  <img
    alt=""
    src={getWebpPath(url)}
    class="w-full h-[55vh] lg:h-[45vh]"
  />
  <div class="title w-full h-full left-[50%] top-[50%] absolute flex flex-col">
    <div class="container sm:px-4 md:px-8 py-3 relative h-full flex flex-col gap-y-[30%]">
      <div class="w-full">
        <Navigation />
      </div>
      <div class="flex flex-col mx-4">
        <h1 class="text-2xl sm:text-3xl font-semibold">
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
