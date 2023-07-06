<script lang="ts">
  export let tags: string;
  const tagList: string[] = tags.split(", ");
</script>

<ul>
  {#each tagList as tag}
    <li>
      <a href="/?tag={tag}">#{tag}</a>
    </li>
  {/each}
</ul>

<style>
  li {
    list-style: none;
  }
  ul {
    display: flex;
    gap: 0 1rem;
  }
</style>
