<script lang="ts">
  import { getWebpPath } from "$lib/utils";

  type Props = {
    url: string;
    image: string;
    title: string;
  };
  const { url, image, title }: Props = $props();
</script>

<a href="/{url}" class="flex h-fit w-full">
  <article class="card card-side bg-base-100 shadow-xl basis-full">
    <figure class="w-[35%] md:w-[30%] lg:w-[25%]">
      <img alt="" height="250" width="250" src={getWebpPath(image)} loading="lazy" />
    </figure>
    <div class="card-body flex justify-center items-center basis-full">
      <h2 class="card-title text-lg">{title}</h2>
    </div>
  </article>
</a>

<style lang="scss">
  article {
    padding: 0;
    &:hover {
      background-color: var(--mt-color-bg-hover);
    }
    & img {
      width: 100%;
      aspect-ratio: 1;
      object-fit: cover;
      border-radius: 5px 5px 0px 0px;
      background-color: hsl(205, 20%, 90%);
    }
  }
</style>
