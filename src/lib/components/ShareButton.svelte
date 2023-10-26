<script lang="ts">
  import { onMount } from "svelte";

  onMount(() => {
    let hasTouchScreen = false;
    if ("userAgentData" in navigator && navigator.userAgentData.mobile) {
      hasTouchScreen = navigator.maxTouchPoints > 0;
    } else if ("maxTouchPoints" in navigator) {
      hasTouchScreen = navigator.maxTouchPoints > 0;
    } else if ("msMaxTouchPoints" in navigator) {
      hasTouchScreen = (navigator as Navigator).msMaxTouchPoints > 0;
    } else {
      const mQ = matchMedia?.("(pointer:coarse)");
      if (mQ?.media === "(pointer:coarse)") {
        hasTouchScreen = !!mQ.matches;
      } else if ("orientation" in window) {
        hasTouchScreen = true; // deprecated, but good fallback
      } else {
        // Only as a last resort, fall back to user agent sniffing
        const UA = (navigator as Navigator).userAgent;
        hasTouchScreen =
          /\b(BlackBerry|webOS|iPhone|IEMobile)\b/i.test(UA) ||
          /\b(Android|Windows Phone|iPad|iPod)\b/i.test(UA);
      }
    }

    if ("share" in navigator && hasTouchScreen) {
      const shareButton = document.getElementById("share-button");
      shareButton!.addEventListener("click", async () => {
        await navigator.share({ title: document.title, url: location.href });
      });
      shareButton!.classList.remove("none");
    }
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
  シェア
</button>

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
</style>