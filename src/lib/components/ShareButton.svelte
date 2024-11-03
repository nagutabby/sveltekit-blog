<script lang="ts">
  import { onMount } from "svelte";
  import {
    toggleModal,
    closeWithClickOutside,
    closeWithEscapeKey,
  } from "$lib/modal";
  import { browser } from "$app/environment";

  const isValidDomain = (domain: string): boolean => {
    const regexp = /^[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,6}$/i;
    return regexp.test(domain);
  };
  let currentUrl = "";

  onMount(() => {
    if (browser) {
      currentUrl = window.location.href;
    }

    const mastodonInstanceNameField = document.getElementById(
      "mastodon-instance-name",
    ) as HTMLInputElement;

    const mastodonShareButton = document.getElementById(
      "mastodon-share-button",
    );
    mastodonShareButton?.addEventListener("click", async (event) => {
      if (localStorage.getItem("mastodon-instance-name")) {
        mastodonInstanceNameField.value = localStorage.getItem(
          "mastodon-instance-name",
        )!;
      }
      toggleModal(event);
    });

    const mastodonShareCloseButtons = Array.from(
      document.querySelectorAll(".mastodon-share-close-button"),
    );
    mastodonShareCloseButtons.forEach((mastodonShareCloseButton) => {
      mastodonShareCloseButton.addEventListener("click", async (event) => {
        toggleModal(event);
      });
    });

    const mastodonShareConfirmButton = document.getElementById(
      "mastodon-share-confirm-button",
    );
    const form = document.getElementById("form") as HTMLFormElement;

    mastodonShareConfirmButton?.addEventListener("click", async (event) => {
      if (form.checkValidity()) {
        if (!isValidDomain(mastodonInstanceNameField.value)) {
          mastodonInstanceNameField.focus();
        } else {
          localStorage.setItem(
            "mastodon-instance-name",
            mastodonInstanceNameField.value,
          );
          history.pushState(
            null,
            document.title,
            location.pathname + location.search,
          );
          window.open(
            `https://${
              mastodonInstanceNameField.value
            }/share?text=${encodeURIComponent(
              document.title,
            )}${encodeURIComponent("\n")}${encodeURIComponent(location.href)}`,
          );
          toggleModal(event);
        }
      }
    });
    closeWithClickOutside();
    closeWithEscapeKey();
  });
</script>

<p class="share-header">記事をシェア</p>
<div class="button-group">
  <button
    data-target="mastodon-modal"
    id="mastodon-share-button"
    aria-label="Mastodon"
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
  </button>
  <a
    role="button"
    href="https://b.hatena.ne.jp/entry/panel/?url={currentUrl}"
    aria-label="はてなブックマーク"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="1em"
      height="1em"
      viewBox="0,0,256,256"
    >
      <g transform="translate(-40.96,-40.96) scale(1.32,1.32)"
        ><g fill="currentColor"
          ><g transform="scale(5.33333,5.33333)"
            ><path
              d="M11.48438,6c-3.032,0 -5.5,2.467 -5.5,5.5v25c0,3.033 2.468,5.5 5.5,5.5h25c3.032,0 5.5,-2.467 5.5,-5.5v-25c0,-3.033 -2.467,-5.5 -5.5,-5.5zM18.0332,15.96289c4.51575,0.00655 7.94922,0.46162 7.94922,3.76562c0,2.616 -1.48228,3.41925 -2.73828,3.65625c1.513,0.07 3.73828,0.891 3.73828,4c0.001,4.754 -5.21102,4.61914 -10.91602,4.61914c-1.148,0 -2.08203,-0.898 -2.08203,-2v-12.02344c0,-1.094 0.91373,-1.983 2.05273,-2c0.68238,-0.01037 1.35099,-0.01851 1.99609,-0.01758zM30.98438,16h2c0.552,0 1,0.448 1,1v9c0,0.552 -0.448,1 -1,1h-2c-0.552,0 -1,-0.448 -1,-1v-9c0,-0.552 0.448,-1 1,-1zM19.35547,19.51758c-0.66639,-0.01888 -1.19922,0.05469 -1.19922,0.05469v3.15234c-0.001,0 3.79102,0.42083 3.79102,-1.57617c0,-1.31125 -1.48115,-1.59939 -2.5918,-1.63086zM19.47266,25.47852c-0.48959,0.01206 -0.95411,0.04297 -1.31836,0.04297v3.45508c1.772,0 4.49414,0.37089 4.49414,-1.78711c0.00075,-1.617 -1.707,-1.74713 -3.17578,-1.71094zM31.98438,28c1.105,0 2,0.895 2,2c0,1.105 -0.895,2 -2,2c-1.105,0 -2,-0.895 -2,-2c0,-1.105 0.895,-2 2,-2z"
            ></path></g
          ></g
        ></g
      >
    </svg>
  </a>
</div>
<dialog id="mastodon-modal">
  <article>
    <button
      aria-label="Close"
      class="close mastodon-share-close-button"
      data-target="mastodon-modal"
    />
    <h3>Mastodonでシェア</h3>
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
      <div class="cancel-and-share-buttons">
        <button
          class="secondary outline mastodon-share-close-button"
          data-target="mastodon-modal"
        >
          キャンセル
        </button>
        <button
          type="submit"
          data-target="mastodon-modal"
          id="mastodon-share-confirm-button"
        >
          シェア
        </button>
      </div>
    </form>
  </article>
</dialog>
