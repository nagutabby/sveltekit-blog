<script lang="ts">
  import type { PageData } from "./$types";
  export let data: PageData;
  import Card from "../lib/components/Card.svelte";
</script>

{#each data.contents as content, i}
  {#if i % 2 === 0}
    <div class="grid">
      <Card
        url={data.contents[i].id}
        image={data.contents[i].image.url}
        title={data.contents[i].title}
        date={data.contents[i].createdAt}
      />
      <Card
        url={data.contents[i + 1].id}
        image={data.contents[i + 1].image.url}
        title={data.contents[i + 1].title}
        date={data.contents[i + 1].createdAt}
      />
    </div>
  {/if}
{/each}
