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
