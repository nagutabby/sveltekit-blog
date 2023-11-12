<script lang="ts">
  import type { PageData } from "./$types";
  import TagList from "$lib/components/TagList.svelte";

  export let data: PageData;
</script>

{#if data.tagSet !== undefined}
  <div class="tag-list">
    <TagList tags={Array.from(data.tagSet)} />
  </div>
{/if}

<style>
  .tag-list {
    margin: 5vh 0;
  }
  :global(ul) {
    display: flex;
    justify-content: center;
  }
</style>
