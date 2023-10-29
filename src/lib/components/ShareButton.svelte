<script lang="ts">
  import { onMount } from "svelte";
  import {
    toggleModal,
    closeWithClickOutside,
    closeWithEscapeKey,
  } from "$lib/modal";

  const isValidDomain = (domain: string): boolean => {
    const regexp = /^[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,6}$/i;
    return regexp.test(domain);
  };

  onMount(() => {
    if ("share" in navigator) {
      const shareButton = document.getElementById("share-button");
      shareButton!.addEventListener("click", async () => {
        history.pushState(
          null,
          document.title,
          location.pathname + location.search
        );
        await navigator.share({
          url: location.href,
          title: document.title,
          text: document.title,
        });
      });
      shareButton!.classList.remove("none");
    }
    const mastodonInstanceNameField = document.getElementById(
      "mastodon-instance-name"
    ) as HTMLInputElement;

    const mastodonShareButton = document.getElementById(
      "mastodon-share-button"
    );
    mastodonShareButton?.addEventListener("click", async (event) => {
      if (localStorage.getItem("mastodon-instance-name")) {
        mastodonInstanceNameField.value = localStorage.getItem(
          "mastodon-instance-name"
        )!;
      }
      toggleModal(event);
    });

    const mastodonShareCloseButtons = Array.from(
      document.querySelectorAll(".mastodon-share-close-button")
    );
    mastodonShareCloseButtons.forEach((mastodonShareCloseButton) => {
      mastodonShareCloseButton.addEventListener("click", async (event) => {
        toggleModal(event);
      });
    });

    const mastodonShareConfirmButton = document.getElementById(
      "mastodon-share-confirm-button"
    );
    const form = document.getElementById("form") as HTMLFormElement;

    mastodonShareConfirmButton?.addEventListener("click", async (event) => {
      if (form.checkValidity()) {
        if (!isValidDomain(mastodonInstanceNameField.value)) {
          mastodonInstanceNameField.focus();
        } else {
          localStorage.setItem(
            "mastodon-instance-name",
            mastodonInstanceNameField.value
          );
          history.pushState(
            null,
            document.title,
            location.pathname + location.search
          );
          window.open(
            `https://${
              mastodonInstanceNameField.value
            }/share?text=${encodeURIComponent(
              document.title
            )}${encodeURIComponent("\n")}${encodeURIComponent(location.href)}`
          );
          toggleModal(event);
        }
      }
    });
    closeWithClickOutside();
    closeWithEscapeKey();
  });
</script>

<button type="button" id="share-button" class="outline none">
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="1em"
    height="1em"
    fill="currentColor"
    class="bi bi-share"
    viewBox="0 0 16 16"
  >
    <path
      d="M13.5 1a1.5 1.5 0 1 0 0 3 1.5 1.5 0 0 0 0-3zM11 2.5a2.5 2.5 0 1 1 .603 1.628l-6.718 3.12a2.499 2.499 0 0 1 0 1.504l6.718 3.12a2.5 2.5 0 1 1-.488.876l-6.718-3.12a2.5 2.5 0 1 1 0-3.256l6.718-3.12A2.5 2.5 0 0 1 11 2.5zm-8.5 4a1.5 1.5 0 1 0 0 3 1.5 1.5 0 0 0 0-3zm11 5.5a1.5 1.5 0 1 0 0 3 1.5 1.5 0 0 0 0-3z"
    />
  </svg>
  アプリでシェア
</button>
<button
  type="button"
  class="outline"
  data-target="mastodon-modal"
  id="mastodon-share-button"
>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="1em"
    height="1em"
    fill="currentColor"
    class="bi bi-mastodon"
    viewBox="0 0 16 16"
  >
    <path
      d="M11.19 12.195c2.016-.24 3.77-1.475 3.99-2.603.348-1.778.32-4.339.32-4.339 0-3.47-2.286-4.488-2.286-4.488C12.062.238 10.083.017 8.027 0h-.05C5.92.017 3.942.238 2.79.765c0 0-2.285 1.017-2.285 4.488l-.002.662c-.004.64-.007 1.35.011 2.091.083 3.394.626 6.74 3.78 7.57 1.454.383 2.703.463 3.709.408 1.823-.1 2.847-.647 2.847-.647l-.06-1.317s-1.303.41-2.767.36c-1.45-.05-2.98-.156-3.215-1.928a3.614 3.614 0 0 1-.033-.496s1.424.346 3.228.428c1.103.05 2.137-.064 3.188-.189zm1.613-2.47H11.13v-4.08c0-.859-.364-1.295-1.091-1.295-.804 0-1.207.517-1.207 1.541v2.233H7.168V5.89c0-1.024-.403-1.541-1.207-1.541-.727 0-1.091.436-1.091 1.296v4.079H3.197V5.522c0-.859.22-1.541.66-2.046.456-.505 1.052-.764 1.793-.764.856 0 1.504.328 1.933.983L8 4.39l.417-.695c.429-.655 1.077-.983 1.934-.983.74 0 1.336.259 1.791.764.442.505.661 1.187.661 2.046v4.203z"
    />
  </svg>
  Mastodonでシェア
</button>

<dialog id="mastodon-modal">
  <article>
    <button
      aria-label="Close"
      class="close mastodon-share-close-button"
      data-target="mastodon-modal"
    />
    <h3>インスタンスを入力🐘</h3>
    <form id="form" on:submit|preventDefault>
      <label for="mastodon-instance-name">
        インスタンス名
        <input
          type="text"
          id="mastodon-instance-name"
          name="mastodon-instance-name"
          placeholder="mastodon.social"
          required
        />
      </label>
      <div class="row">
        <div class="col-6">
          <button
            class="secondary mastodon-share-close-button"
            data-target="mastodon-modal"
          >
            キャンセル
          </button>
        </div>
        <div class="col-6">
          <button
            type="submit"
            data-target="mastodon-modal"
            id="mastodon-share-confirm-button"
          >
            シェア
          </button>
        </div>
      </div>
    </form>
  </article>
</dialog>

<style>
  button {
    display: flex;
    gap: 0 1rem;
    align-items: center;
    justify-content: center;
    padding: var(--form-element-spacing-vertical)
      var(--form-element-spacing-horizontal);
  }
  :global(.none) {
    display: none !important;
  }
  .close {
    background-color: transparent;
    border-color: transparent;
    &:focus {
      box-shadow: none;
    }
  }
</style>
