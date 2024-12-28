<script lang="ts">
  import hljs from "highlight.js/lib/core";
  import python from "highlight.js/lib/languages/python";
  import ruby from "highlight.js/lib/languages/ruby";
  import erb from "highlight.js/lib/languages/erb";
  import bash from "highlight.js/lib/languages/bash";
  import css from "highlight.js/lib/languages/css";
  import scss from "highlight.js/lib/languages/scss";
  import javascript from "highlight.js/lib/languages/javascript";
  import typescript from "highlight.js/lib/languages/typescript";
  import yaml from "highlight.js/lib/languages/yaml";
  import json from "highlight.js/lib/languages/json";
  import dockerfile from "highlight.js/lib/languages/dockerfile";
  import Toc from "$lib/components/Toc.svelte";
  import type { PageData } from "./$types";
  import ShareButton from "$lib/components/ShareButton.svelte";
  import Date from "$lib/components/Date.svelte";
  import { onMount } from "svelte";
  import "../../app.scss";
  import Header from "$lib/components/Header.svelte";
  import Footer from "$lib/components/Footer.svelte";
  import OpenGraph from "$lib/components/OpenGraph.svelte";

  const { data }: { data: PageData } = $props();

  hljs.registerLanguage("python", python);
  hljs.registerLanguage("ruby", ruby);
  hljs.registerLanguage("erb", erb);
  hljs.registerLanguage("bash", bash);
  hljs.registerLanguage("css", css);
  hljs.registerLanguage("scss", scss);
  hljs.registerLanguage("javascript", javascript);
  hljs.registerLanguage("typescript", typescript);
  hljs.registerLanguage("json", json);
  hljs.registerLanguage("dockerfile", dockerfile);
  hljs.registerLanguage("yaml", yaml);

  let headingList = $state();
  let isLoading = $state(true);

  onMount(async () => {
    headingList = Array.from(
      document.getElementById("content")!.querySelectorAll("h1, h2, h3"),
    );
    isLoading = false;

    document.querySelectorAll("pre code").forEach((element) => {
      hljs.highlightElement(element as HTMLElement);
    });
    document.querySelectorAll("pre").forEach((element) => {
      const copyButton = document.createElement("button");
      element.insertAdjacentElement("beforeend", copyButton);
      copyButton.innerHTML = `
      <svg xmlns="http://www.w3.org/2000/svg" width="16px" height="16px" fill="currentColor" class="bi bi-clipboard" viewBox="0 0 16 16">
        <path d="M4 1.5H3a2 2 0 0 0-2 2V14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V3.5a2 2 0 0 0-2-2h-1v1h1a1 1 0 0 1 1 1V14a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V3.5a1 1 0 0 1 1-1h1v-1z"/>
        <path d="M9.5 1a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-3a.5.5 0 0 1-.5-.5v-1a.5.5 0 0 1 .5-.5h3zm-3-1A1.5 1.5 0 0 0 5 1.5v1A1.5 1.5 0 0 0 6.5 4h3A1.5 1.5 0 0 0 11 2.5v-1A1.5 1.5 0 0 0 9.5 0h-3z"/>
      </svg>`;
      copyButton.setAttribute("data-tooltip", "コピー");
      copyButton.setAttribute("data-placement", "left");
      copyButton.setAttribute("aria-label", "コピー");
      copyButton.addEventListener("click", async (event) => {
        const code = copyButton.previousElementSibling?.textContent!;

        navigator.clipboard.writeText(code).then(() => {
          copyButton.innerHTML = `
          <svg xmlns="http://www.w3.org/2000/svg" width="16px" height="16px" fill="currentColor" class="bi bi-check-lg" viewBox="0 0 16 16">
            <path d="M12.736 3.97a.733.733 0 0 1 1.047 0c.286.289.29.756.01 1.05L7.88 12.01a.733.733 0 0 1-1.065.02L3.217 8.384a.757.757 0 0 1 0-1.06.733.733 0 0 1 1.047 0l3.052 3.093 5.4-6.425a.247.247 0 0 1 .02-.022Z"/>
          </svg>`;
          copyButton.setAttribute("data-tooltip", "コピーしました！");
          setTimeout(() => {
            copyButton.innerHTML = `
            <svg xmlns="http://www.w3.org/2000/svg" width="16px" height="16px" fill="currentColor" class="bi bi-clipboard" viewBox="0 0 16 16">
              <path d="M4 1.5H3a2 2 0 0 0-2 2V14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V3.5a2 2 0 0 0-2-2h-1v1h1a1 1 0 0 1 1 1V14a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V3.5a1 1 0 0 1 1-1h1v-1z"/>
              <path d="M9.5 1a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-3a.5.5 0 0 1-.5-.5v-1a.5.5 0 0 1 .5-.5h3zm-3-1A1.5 1.5 0 0 0 5 1.5v1A1.5 1.5 0 0 0 6.5 4h3A1.5 1.5 0 0 0 11 2.5v-1A1.5 1.5 0 0 0 9.5 0h-3z"/>
            </svg>`;
            copyButton.setAttribute("data-tooltip", "コピー");
          }, 3000);
        });
      });
    });
  });
</script>

<OpenGraph
  url={data.image.url}
  title={data.title}
  description={data.description}
></OpenGraph>

<Header url={data.image.url} title={data.title} description={data.description}
></Header>

<main class="container">
  <div class="article-group">
    <div class="sidebar">
      <div class="sticky-top">
        <div class="toc">
          {#if isLoading}
            <div aria-busy="true" class="loading-toc"></div>
          {:else}
            <Toc { headingList } />
          {/if}
        </div>
        <hr />
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
</main>

<Footer></Footer>

<style lang="scss">
  main.container {
    margin-top: 2rem;
  }
  .loading-toc {
    height: 20vh;
  }
</style>
