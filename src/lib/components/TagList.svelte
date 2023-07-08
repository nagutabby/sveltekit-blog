<script lang="ts">
  export let tags: {
    id: string;
    createdAt: string;
    updatedAt: string;
    publishedAt: string;
    revisedAt: string;
    name: string;
  }[];
</script>

<ul>
  {#each tags as tag}
    <li>
      <a href="/?tag={tag.id}">#{tag.name}</a>
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
