<script lang="ts">
  const { rating }: { rating: number } = $props();
</script>

<div class="rating flex gap-x-1">
  {#each Array(5) as _, i}
    <div
      class="mask mask-star"
      aria-label="{i + 1} star"
      aria-current={i + 1 === rating ? "true" : undefined}
    ></div>
  {/each}
</div>
