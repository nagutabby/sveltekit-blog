<script lang="ts">
  export let tags: string[];
</script>

<details role="list">
  <summary aria-haspopup="listbox">タグ</summary>
  <ul role="listbox">
    {#each tags as tag}
      <li>
      <a href="/tags/{tag.replaceAll(' ', '-').toLowerCase()}">{tag}</a>
      </li>
    {/each}
  </ul>
</details>

<style>
  details {
    margin-top: var(--block-spacing-vertical);
    margin-bottom: 0;
    display: flex;
    & summary {
      display: flex;
      justify-content: space-between;
      padding: var(--form-element-spacing-vertical)
      var(--form-element-spacing-horizontal) !important;
    }
  }
</style>
