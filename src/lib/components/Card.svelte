<script lang="ts">
  import { getWebpPath } from "$lib/utils";

  type Props = {
    id: string;
    url: string;
    image: string;
    title: string;
  };
  const { id, url, image, title }: Props = $props();
</script>

<a href="/{url}" class="flex h-full w-full">
  <article class="card bg-base-100 shadow-xl basis-full image-full w-full">
    <figure class="w-full">
      <img
        alt=""
        height="250"
        width="250"
        src={getWebpPath(image)}
        loading="lazy"
        style="view-transition-name: article-img-{id};"
      />
    </figure>
    <div
      class="card-body flex justify-center items-center basis-full z-10 p-4 md:p-6"
    >
      <div class="w-full text-center">
        <h2
          class="card-title text-xl lg:text-xl text-white inline-block"
          style="view-transition-name: article-title-{id};"
        >
          {title}
        </h2>
      </div>
    </div>
  </article>
</a>

<style lang="scss">
  article {
    padding: 0;
    width: 100%;
    &:hover {
      opacity: 0.7;
    }

    & img {
      width: 100%;
      aspect-ratio: 1 / 1;
      object-fit: cover;
      border-radius: 5px 5px 0px 0px;
    }
  }
</style>
