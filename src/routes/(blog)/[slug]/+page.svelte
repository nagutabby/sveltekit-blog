<script lang="ts">
  import hljs from "highlight.js/lib/core";
  import python from "highlight.js/lib/languages/python";
  import ruby from "highlight.js/lib/languages/ruby";
  import erb from "highlight.js/lib/languages/erb";
  import bash from "highlight.js/lib/languages/bash"
  import css from "highlight.js/lib/languages/css";
  import scss from "highlight.js/lib/languages/scss";
  import javascript from "highlight.js/lib/languages/javascript";
  import typescript from "highlight.js/lib/languages/typescript";
  import yaml from "highlight.js/lib/languages/yaml";
  import json from "highlight.js/lib/languages/json";
  import dockerfile from "highlight.js/lib/languages/dockerfile";
  import Toc from "$lib/components/Toc.svelte";
  import type { PageData } from "./$types";
  import TagDropdown from "$lib/components/TagDropdown.svelte";
  import ShareButton from "$lib/components/ShareButton.svelte";
  import Date from "$lib/components/Date.svelte";
  import { onMount } from "svelte";

  export let data: PageData;

  hljs.registerLanguage("python", python)
  hljs.registerLanguage("ruby", ruby)
  hljs.registerLanguage("erb", erb)
  hljs.registerLanguage("bash", bash)
  hljs.registerLanguage("css", css)
  hljs.registerLanguage("scss", scss)
  hljs.registerLanguage("javascript", javascript)
  hljs.registerLanguage("typescript", typescript)
  hljs.registerLanguage("json", json)
  hljs.registerLanguage("dockerfile", dockerfile)
  hljs.registerLanguage("yaml", yaml)

  onMount(async () => {
    document.querySelectorAll("pre code").forEach((element) => {
      hljs.highlightElement(element as HTMLElement);
    });
  })
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
      <Date createdAt={data.createdAt} revisedAt={data.revisedAt} />
      {@html data.body}
    </article>
  </div>
</div>
