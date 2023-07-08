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
      <a href="/?tag={tag.id}" role="button" class="outline">
        #{tag.name}
      </a>
    </li>
  {/each}
</ul>

<style>
  li {
    list-style: none;
    margin: 0;
  }
  ul {
    display: flex;
    padding: 0;
    gap: 0 1rem;
  }
</style>
