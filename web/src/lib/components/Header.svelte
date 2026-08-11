<script lang="ts">
  import Navigation from "./Navigation.svelte";
  import { getWebpPath } from "$lib/utils";

  type Props = {
    id?: string;
    url: string;
    title: string;
  };
  const { id, url, title }: Props = $props();
</script>

<svelte:head>
  <link rel="preload" as="image" href={getWebpPath(url)} fetchpriority="high" />
</svelte:head>

<header
  class="card bg-base-100 image-full relative w-full h-[60vh] md:h-[55vh] lg:h-[50vh] rounded-none"
>
  <figure class="w-full h-full">
    <img
      alt=""
      src={getWebpPath(url)}
      class="w-full h-full"
      fetchpriority="high"
      style={id ? `view-transition-name: article-img-${id};` : ""}
    />
  </figure>

  <div
    class="card-body title w-full h-full absolute z-10 px-0 sm:px-10 md:px-15 lg:px-25 xl:px-30"
  >
    <div class="py-1 relative h-full flex flex-col">
      <div class="w-full">
        <Navigation />
      </div>
      <div class="flex-1 flex flex-col justify-center text-center sm:text-left">
        <div>
          <h1
            class="text-2xl sm:text-3xl font-semibold text-white inline-block"
            style={id ? `view-transition-name: article-title-${id};` : ""}
          >
            {title}
          </h1>
        </div>
      </div>
    </div>
  </div>
</header>

<style lang="scss">
  header {
    padding: 0 !important;

    &::before {
      z-index: 1;
    }

    & img {
      aspect-ratio: 1 / 1;
      object-fit: cover;
    }
  }
</style>
