<script lang="ts">
  export let tags: string[];
</script>

<details class="dropdown">
  <summary>タグ</summary>
  <ul>
    {#each tags as tag}
      <li>
        <a href="/tags/{tag.replaceAll(' ', '-').toLowerCase()}">{tag}</a>
      </li>
    {/each}
  </ul>
</details>
