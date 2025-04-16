<script lang="ts">
  type Props = {
    url: string;
    image: string;
    title: string;
  };
  const { url, image, title }: Props = $props();
</script>

<a href="/{url}" class="flex h-fit w-full">
  <article class="card card-side bg-base-100 shadow-xl basis-full">
    <figure class="w-[30%] md:w-[25%] lg:w-[20%]">
      <img
        alt=""
        height="250"
        width="250"
        srcset="{image}?fm=webp&fit=clip&w=250&h=250 250w, {image}?fm=webp&fit=clip&w=500&h=500 500w"
        loading="lazy"
      />
    </figure>
    <div class="card-body flex justify-center items-center">
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
