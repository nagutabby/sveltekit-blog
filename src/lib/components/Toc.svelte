<script lang="ts">
  export let toc: { text: any; id: string; name: string }[];

  let tmpItemName = "";
  let count = 0;
  function clearCount() {
    count = 0;
  }
  function incrementCount() {
    count += 1;
  }
  function updateTmpItemName(itemName: string) {
    tmpItemName = itemName;
  }
</script>

<article>
  <details>
    <summary>目次</summary>
    <ul>
      {#each toc as item, i}
        {#if i === 0}
          <li>
            <a href="#{item.id}">{item.text}</a>
          </li>
          {updateTmpItemName(item.name)}
          a
        {:else if tmpItemName === item.name}
          {incrementCount()}
        {:else if tmpItemName !== item.name}
          {clearCount()}
        {/if}
      {/each}
    </ul>
  </details>
</article>
