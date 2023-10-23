<script lang="ts">
  export let tags: string[];
</script>

<ul>
  {#each tags as tag}
    <li>
      <a
        href="/tags/{tag.replaceAll(' ', '-').toLowerCase()}"
        role="button"
        class="outline"
      >
        #{tag}
      </a>
    </li>
  {/each}
</ul>

<style>
  ul {
    display: flex;
    flex-wrap: wrap;
    padding: 0;
    gap: 1rem;
    & li {
      list-style: none;
      margin: 0;
    }
  }
</style>
