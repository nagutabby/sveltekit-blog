<script lang="ts">
  export let toc: { text: any; id: string; name: string }[];
</script>

<article>
  <details>
    <summary>目次</summary>
    <ul>
      {#each toc as item}
        <li class={item.name}>
          <a href="#{item.id}">{item.text}</a>
        </li>
      {/each}
    </ul>
  </details>
</article>

<style>
  ul {
    padding: 0;
    margin-bottom: 0;
  }

  details {
    border: none;
    padding-bottom: 0;
    margin-bottom: 0;
  }

  li::before {
    content: "•";
    position: relative;
    top: 2px;
    padding-right: 5px;
    font-size: 1.3rem;
  }
  li {
    list-style-type: none;
    margin: 0;
  }
  .h2 {
    padding-left: 1.2rem;
  }
  .h3 {
    padding-left: 2.4rem;
  }
</style>
