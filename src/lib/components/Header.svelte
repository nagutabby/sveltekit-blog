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
