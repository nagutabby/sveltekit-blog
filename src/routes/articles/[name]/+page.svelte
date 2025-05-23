<script lang="ts">
  import { codeToHtml } from "shiki";
  import Toc from "$lib/components/Toc.svelte";
  import type { PageData } from "./$types";
  import ShareButton from "$lib/components/ShareButton.svelte";
  import Date from "$lib/components/Date.svelte";
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import Header from "$lib/components/Header.svelte";
  import OpenGraph from "$lib/components/OpenGraph.svelte";

  const { data }: { data: PageData } = $props();

  // サポートする言語の配列
  const languages = [
    "python",
    "ruby",
    "bash",
    "html",
    "xml",
    "css",
    "javascript",
    "typescript",
    "json",
    "dockerfile",
    "yaml",
    "java",
  ];

  let headingList = $state();
  let isLoading = $state(true);

  onMount(async () => {
    // headingListの取得
    headingList = Array.from(
      document.getElementById("content")!.querySelectorAll("h1, h2, h3"),
    );

    try {
      // codeブロックにシンタックスハイライトを適用
      const codeBlocks = document.querySelectorAll("pre code");
      for (const element of codeBlocks) {
        const code = element.textContent || "";
        const classNames = element.getAttribute("class") || "";
        // 言語名を取得（例: language-javascript → javascript）
        const langMatch = classNames.match(/language-(\w+)/);
        const lang = langMatch ? langMatch[1] : "text";

        // 言語がサポートされているか確認
        if (languages.includes(lang)) {
          try {
            // shikiの新しいAPIを使用してハイライト
            const html = await codeToHtml(code, {
              lang,
              theme: "dracula",
            });

            // HTMLをパースして必要な部分だけを取得
            const temp = document.createElement("div");
            temp.innerHTML = html;
            const highlightedCode = temp.querySelector("pre code");

            if (highlightedCode) {
              element.innerHTML = highlightedCode.innerHTML;
              // 親要素のpreにshikiのクラスを追加
              element.parentElement?.classList.add("shiki");
            }
          } catch (error) {
            console.error(`Failed to highlight ${lang} code:`, error);
          }
        }
      }

      // コピーボタンの追加
      document.querySelectorAll("pre").forEach((element) => {
        const wrapper = document.createElement("div");
        wrapper.style.position = "relative";
        wrapper.classList.add("code-wrapper");

        // `pre`をラッパーに移動
        element.parentElement?.insertBefore(wrapper, element);
        wrapper.appendChild(element);

        const copyButton = document.createElement("button");
        wrapper.appendChild(copyButton);
        copyButton.innerHTML = `
        <svg xmlns="http://www.w3.org/2000/svg" fill="currentColor" class="bi bi-clipboard" viewBox="0 0 16 16">
          <path d="M4 1.5H3a2 2 0 0 0-2 2V14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V3.5a2 2 0 0 0-2-2h-1v1h1a1 1 0 0 1 1 1V14a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V3.5a1 1 0 0 1 1-1h1v-1z"/>
          <path d="M9.5 1a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-3a.5.5 0 0 1-.5-.5v-1a.5.5 0 0 1 .5-.5h3zm-3-1A1.5 1.5 0 0 0 5 1.5v1A1.5 1.5 0 0 0 6.5 4h3A1.5 1.5 0 0 0 11 2.5v-1A1.5 1.5 0 0 0 9.5 0h-3z"/>
        </svg>`;
        copyButton.setAttribute("aria-label", "コピー");
        copyButton.classList.add("copy-button");
        copyButton.addEventListener("click", async () => {
          const code = copyButton.previousElementSibling?.textContent!;

          navigator.clipboard.writeText(code).then(() => {
            copyButton.innerHTML = `
            <svg xmlns="http://www.w3.org/2000/svg" fill="currentColor" class="bi bi-check-lg" viewBox="0 0 16 16">
              <path d="M12.736 3.97a.733.733 0 0 1 1.047 0c.286.289.29.756.01 1.05L7.88 12.01a.733.733 0 0 1-1.065.02L3.217 8.384a.757.757 0 0 1 0-1.06.733.733 0 0 1 1.047 0l3.052 3.093 5.4-6.425a.247.247 0 0 1 .02-.022Z"/>
            </svg>`;
            setTimeout(() => {
              copyButton.innerHTML = `
              <svg xmlns="http://www.w3.org/2000/svg" fill="currentColor" class="bi bi-clipboard" viewBox="0 0 16 16">
                <path d="M4 1.5H3a2 2 0 0 0-2 2V14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V3.5a2 2 0 0 0-2-2h-1v1h1a1 1 0 0 1 1 1V14a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V3.5a1 1 0 0 1 1-1h1v-1z"/>
                <path d="M9.5 1a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-3a.5.5 0 0 1-.5-.5v-1a.5.5 0 0 1 .5-.5h3zm-3-1A1.5 1.5 0 0 0 5 1.5v1A1.5 1.5 0 0 0 6.5 4h3A1.5 1.5 0 0 0 11 2.5v-1A1.5 1.5 0 0 0 9.5 0h-3z"/>
              </svg>`;
            }, 3000);
          });
        });
      });
    } catch (error) {
      console.error("Shikiの処理中にエラーが発生しました:", error);
    } finally {
      isLoading = false;
    }
  });
</script>

<OpenGraph
  url={page.data.image}
  title={page.data.title}
  body={page.data.body}
></OpenGraph>

<Header url={page.data.image} title={page.data.title}></Header>

<main class="container px-3 md:px-10 py-10 mx-auto">
  <div class="flex lg:flex-row items-start justify-between flex-col-reverse">
    <div
      class="flex flex-wrap gap-5 justify-center w-full lg:w-[63%]"
      id="content"
    >
      <article class="prose max-w-full">
        <div class="my-3">
          <Date publishedAt={data.publishedAt} updatedAt={data.updatedAt} />
        </div>
        {@html data.body}
      </article>
    </div>
    <div class="w-full lg:w-[33%] prose lg:top-0 lg:sticky py-5">
      <div class="toc">
        {#if isLoading}
          <div class="flex justify-center">
            <span class="loading loading-spinner loading-lg"></span>
          </div>
        {:else}
          <Toc {headingList} />
        {/if}
      </div>
      <hr />
      <ShareButton />
    </div>
  </div>
</main>

<style>

  :global(pre) {
    /* コピーボタンのために相対配置を全てのpreに適用 */
    position: relative !important;
    padding: 0.5rem;

  }

  :global(.copy-button) {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    height: inherit;
    background: rgba(0, 0, 0, 0.4);
    border: none;
    border-radius: 0.25rem;
    cursor: pointer;
    opacity: 0;
    transition: opacity 0.2s;
  }

  :global(.code-wrapper:hover .copy-button) {
    opacity: 1;
  }

  :global(.copy-button svg) {
    width: 1.3rem;
    height: 1.3rem;
  }
</style>
