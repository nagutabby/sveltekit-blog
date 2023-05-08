<script lang="ts">
  import type { PageData } from "./$types";
  export let data: PageData;
</script>

<svelte:head>
  <title>ブログ</title>
  <meta name="description" content={data.description} />
</svelte:head>

<section>
  <h1>{data.title}</h1>
  <img src={data.image.url} alt="" />
  <div>{@html data.body}</div>
</section>

<style>
  section {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
  }
</style>
