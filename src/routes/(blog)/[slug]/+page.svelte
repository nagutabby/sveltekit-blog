<script lang="ts">
  import Toc from "$lib/components/Toc.svelte";
  import type { PageData } from "./$types";
  import TagDropdown from "$lib/components/TagDropdown.svelte";
  import ShareButton from "$lib/components/ShareButton.svelte";
  export let data: PageData;
</script>

<div class="article-group">
  <div class="tag-dropdown-and-toc">
    {#if data.tags !== undefined}
        <TagDropdown tags={data.tags.split(",")} />
    {/if}
    <div class="sticky-top">
      <div class="toc">
        <Toc />
      </div>
      <ShareButton />
    </div>
  </div>
  <div id="content" class="article-content">
    <article>
      {@html data.body}
    </article>
  </div>
</div>
