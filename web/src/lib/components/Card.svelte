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
  <article class="card bg-base-100 shadow-xl overflow-hidden flex">
    <figure class="object-contain w-45 mx-auto pt-5">
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
      class="card-body flex justify-center items-center z-10"
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
  }
</style>
