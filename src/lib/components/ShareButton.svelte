<script lang="ts">
  import { onMount } from "svelte";
  import { browser } from "$app/environment";

  let isError = $state(false);
  let mastodonShareText = $state("");
  let blueskyShareText = $state("");

  const handleSubmit = (event: SubmitEvent) => {
    event.preventDefault();
    const mastodonInstanceNameField = document.getElementById(
      "mastodon-instance-name",
    ) as HTMLInputElement;
    const mastodonShareModal = document.getElementById(
      "search-share-modal",
    ) as HTMLDialogElement;
    isError = !isValidDomain(mastodonInstanceNameField.value);
    if (!isError) {
      if (mastodonShareModal) {
        mastodonShareModal.close();
      }
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
        `https://${mastodonInstanceNameField.value}/share?text=${mastodonShareText}`,
      );
    }
  };

  const isValidDomain = (domain: string): boolean => {
    const regexp = /^[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,6}$/i;
    return regexp.test(domain);
  };
  let currentUrl = $state("");

  onMount(() => {
    mastodonShareText = `${encodeURIComponent(document.title)}${encodeURIComponent("\n")}${encodeURIComponent(location.href)}`;
    blueskyShareText = `${encodeURIComponent(document.title)}${encodeURIComponent("<br>")}${encodeURIComponent(location.href)}`;
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
    });
  });
</script>

<p class="text-center mb-7 text-xl">記事をシェア</p>
<div class="button-group">
  <button
    class="btn"
    aria-label="Mastodon"
    onclick={() => {
      const mastodonShareModal = document.getElementById(
        "mastodon-share-modal",
      ) as HTMLDialogElement;
      mastodonShareModal.showModal();
    }}
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="1em"
      height="1em"
      fill="currentColor"
      class="bi bi-mastodon text-secondary"
      viewBox="0 0 16 16"
    >
      <path
        d="M11.19 12.195c2.016-.24 3.77-1.475 3.99-2.603.348-1.778.32-4.339.32-4.339 0-3.47-2.286-4.488-2.286-4.488C12.062.238 10.083.017 8.027 0h-.05C5.92.017 3.942.238 2.79.765c0 0-2.285 1.017-2.285 4.488l-.002.662c-.004.64-.007 1.35.011 2.091.083 3.394.626 6.74 3.78 7.57 1.454.383 2.703.463 3.709.408 1.823-.1 2.847-.647 2.847-.647l-.06-1.317s-1.303.41-2.767.36c-1.45-.05-2.98-.156-3.215-1.928a3.614 3.614 0 0 1-.033-.496s1.424.346 3.228.428c1.103.05 2.137-.064 3.188-.189zm1.613-2.47H11.13v-4.08c0-.859-.364-1.295-1.091-1.295-.804 0-1.207.517-1.207 1.541v2.233H7.168V5.89c0-1.024-.403-1.541-1.207-1.541-.727 0-1.091.436-1.091 1.296v4.079H3.197V5.522c0-.859.22-1.541.66-2.046.456-.505 1.052-.764 1.793-.764.856 0 1.504.328 1.933.983L8 4.39l.417-.695c.429-.655 1.077-.983 1.934-.983.74 0 1.336.259 1.791.764.442.505.661 1.187.661 2.046v4.203z"
      />
    </svg>
  </button>
  <dialog id="mastodon-share-modal" class="modal">
    <div class="modal-box">
      <p class="font-bold text-2xl m-0">Mastodonでシェア</p>
      <p class="my-3">インスタンス名</p>
      <form id="form" onsubmit={handleSubmit}>
        <label
          for="mastodon-instance-name"
          class={`flex-col input input-bordered flex items-center gap-2 w-full ${isError ? "input-error" : ""}`}
        >
          <input
            type="text"
            id="mastodon-instance-name"
            name="mastodon-instance-name"
            placeholder="mastodon.social"
            class="grow w-full"
            required
          />
        </label>
      </form>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button>close</button>
    </form>
  </dialog>
  <a
    role="button"
    class="btn"
    href="https://b.hatena.ne.jp/entry/panel/?url={currentUrl}"
    aria-label="はてなブックマーク"
    target="_blank"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="1em"
      height="1em"
      viewBox="0,0,256,256"
      class="text-secondary"
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
  <a
    role="button"
    class="btn"
    href="https://bsky.app/intent/compose?text={blueskyShareText}"
    aria-label="Bluesky"
    target="_blank"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 512 512"
      fill="currentColor"
      class="text-secondary"
      ><!--!Font Awesome Free 6.7.2 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2025 Fonticons, Inc.--><path
        d="M111.8 62.2C170.2 105.9 233 194.7 256 242.4c23-47.6 85.8-136.4 144.2-180.2c42.1-31.6 110.3-56 110.3 21.8c0 15.5-8.9 130.5-14.1 149.2C478.2 298 412 314.6 353.1 304.5c102.9 17.5 129.1 75.5 72.5 133.5c-107.4 110.2-154.3-27.6-166.3-62.9l0 0c-1.7-4.9-2.6-7.8-3.3-7.8s-1.6 3-3.3 7.8l0 0c-12 35.3-59 173.1-166.3 62.9c-56.5-58-30.4-116 72.5-133.5C100 314.6 33.8 298 15.7 233.1C10.4 214.4 1.5 99.4 1.5 83.9c0-77.8 68.2-53.4 110.3-21.8z"
      /></svg
    >
  </a>
</div>

<style lang="scss">
  .button-group {
    margin-bottom: 2rem;
    display: flex;
    justify-content: space-evenly;
    align-items: center;
    & a,
    button {
      display: flex;
      align-items: center;
      min-width: 50px;
      color: var(--pico-primary);
      background-color: transparent;
      border: transparent;
      padding: 0;
      margin: 0;
      display: block;
      & svg {
        width: 100%;
        height: 100%;
      }
    }
  }
</style>
