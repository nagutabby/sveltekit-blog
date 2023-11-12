<script lang="ts">
  import Toc from "$lib/components/Toc.svelte";
  import type { PageData } from "./$types";
  import TagDropdown from "$lib/components/TagDropdown.svelte";
  import ShareButton from "$lib/components/ShareButton.svelte";
  export let data: PageData;
</script>

<div class="row article-group">
  <div class="col-12 col-lg-4 right-group">
    {#if data.tags !== undefined}
      <div class="col-12">
        <TagDropdown tags={data.tags.split(",")} />
      </div>
    {/if}
    <div class="sticky-top">
      <div class="col-12">
        <div class="toc">
          <Toc />
        </div>
      </div>
      <div class="col-12">
        <ShareButton />
      </div>
    </div>
  </div>
  <div class="col-12 col-lg-8" id="content">
    <article>
      {@html data.body}
    </article>
  </div>
</div>

<style>
  .article-group {
    display: flex;
    flex-direction: row-reverse;
  }
  .toc {
    & details[open] > ul {
      max-height: 20vh;
      overflow-y: auto;
    }
  }
  :global(pre > code) {
    font-size: 0.85rem;
  }
  @media screen and (min-width: 992px) {
    :global(.toc > article) {
      margin: var(--block-spacing-vertical) 0 !important;
    }
    :global(.toc > article > details[open] > ul) {
      max-height: 40vh !important;
    }
    .sticky-top {
      position: sticky;
      top: 0;
    }
  }
  @media (prefers-color-scheme: light) {
    :global(code) {
      color: #545454;
    }
    :root {
      --primary-hover: #0d47a1;
    }
    :global(.active) {
      background: linear-gradient(transparent 60%, #b2ebf2 60%);
    }
  }
  @media (prefers-color-scheme: dark) {
    :global(code) {
      color: #f8f8f2;
    }
    :root {
      --primary-hover: #29b6f6;
    }
    :global(.active) {
      background: linear-gradient(transparent 60%, #455a64 60%);
    }
  }
  :global(a) {
    text-decoration: none;
  }
  :root {
    --primary: #1e88e5 !important;
  }
</style>