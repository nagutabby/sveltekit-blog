<script lang="ts">
  export let toc: string;
</script>

<article>
  <details open>
    <summary>目次</summary>
    {@html toc}
  </details>
</article>

<style>
  :global(article > details > ul) {
    margin-bottom: 0;
  }
</style>
