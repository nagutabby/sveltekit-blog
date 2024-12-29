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

<div id="mt-container" class="mt-container w-full !bg-base-100 rounded-md">
  <div class="mt-body p-0" role="feed">
    <div class="flex justify-center">
      <span class="loading loading-spinner loading-lg mt-5"></span>
    </div>
  </div>
</div>
