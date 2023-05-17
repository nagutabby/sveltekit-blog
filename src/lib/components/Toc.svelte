<script lang="ts">
  export let toc: { text: any; id: string; name: string }[];
</script>

<article>
  <details>
    <summary>目次</summary>
    <ul>
      {#each toc as item}
        <li>
          <a href="#{item.id}" class={item.name}>{item.text}</a>
        </li>
      {/each}
    </ul>
  </details>
</article>
