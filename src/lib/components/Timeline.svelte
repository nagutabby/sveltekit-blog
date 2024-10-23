<script lang="ts">
  import * as MastodonTimeline from "@idotj/mastodon-embed-timeline";
  import "@idotj/mastodon-embed-timeline/dist/mastodon-timeline.min.css";
  import { onMount } from "svelte";

  onMount(() => {
    new MastodonTimeline.Init({
      instanceUrl: "https://mastodon.social",
      timelineType: "profile",
      userId: "111085130737386025",
      profileName: "@nagutabby",
      maxNbPostFetch: "10",
      maxNbPostShow: "10",
      dateLocale: "ja-JP",
      btnSeeMore: "",
      btnReload: "",
      defaultTheme: "dark",
    });
  });
</script>

<div id="mt-container" class="mt-container">
  <div class="mt-body" role="feed">
    <article aria-busy="true" class="loading-article"></article>
  </div>
</div>

<style lang="scss">
  .mt-container {
    background-color: var(--pico-card-background-color);
    padding: 0.5rem;
    border-radius: 5px;
    scrollbar-color: var(--pico-muted-color) var(--pico-card-background-color);
  }

  .mt-body {
    position: relative;
    padding: 0;
    margin: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-flow: column;
  }
  .loading-article {
    height: 100%;
    width: 100%;
    margin: 0;
  }
</style>
