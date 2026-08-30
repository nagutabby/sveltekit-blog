<script lang="ts">
  import * as MastodonTimeline from "@idotj/mastodon-embed-timeline";
  import "@idotj/mastodon-embed-timeline/dist/mastodon-timeline.min.css";
  import { onMount, onDestroy } from "svelte";

  const INSTANCE_URL = "https://mastodon.social";
  // アバターは最大でも.mt-post-avatar-image-big(2.25rem=36px)でしか
  // 表示されないため、Retina考慮の2倍(72px)あれば十分。Mastodon側は
  // avatar/avatar_staticとも常に400x400で返してくるため、そのまま使うと
  // 表示サイズの100倍近いバイト数を配信することになる。
  const AVATAR_SIZE = 72;

  function toOptimizedAvatarUrl(originalUrl: string): string {
    const params = new URLSearchParams({
      url: originalUrl.replace(/^https?:\/\//, ""),
      w: String(AVATAR_SIZE),
      h: String(AVATAR_SIZE),
      fit: "cover",
      output: "webp",
    });
    return `https://images.weserv.nl/?${params}`;
  }

  function rewriteAvatarUrls(value: unknown): void {
    if (Array.isArray(value)) {
      value.forEach(rewriteAvatarUrls);
      return;
    }
    if (value && typeof value === "object") {
      const record = value as Record<string, unknown>;
      for (const key of ["avatar", "avatar_static"]) {
        if (typeof record[key] === "string") {
          record[key] = toOptimizedAvatarUrl(record[key] as string);
        }
      }
      Object.values(record).forEach(rewriteAvatarUrls);
    }
  }

  let originalFetch: typeof window.fetch;

  onMount(() => {
    // @idotj/mastodon-embed-timelineはavatar/avatar_staticのURLを
    // Mastodon APIのレスポンスからそのまま<img src>へ渡してしまい、URLを
    // 差し替えるオプションが無い(内部メソッドは全てprivateで上書き不可)。
    // このためAPIレスポンスの時点でURLを最適化版に書き換える。
    originalFetch = window.fetch.bind(window);
    window.fetch = async (...args) => {
      const response = await originalFetch(...args);
      const url = typeof args[0] === "string" ? args[0] : args[0].toString();

      if (!url.startsWith(INSTANCE_URL)) {
        return response;
      }

      const data = await response.clone().json().catch(() => null);
      if (data === null) {
        return response;
      }

      rewriteAvatarUrls(data);

      return new Response(JSON.stringify(data), {
        status: response.status,
        statusText: response.statusText,
        headers: response.headers,
      });
    };

    new MastodonTimeline.Init({
      instanceUrl: INSTANCE_URL,
      timelineType: "profile",
      userId: "111085130737386025",
      profileName: "@nagutabby",
      maxNbPostFetch: "10",
      maxNbPostShow: "10",
      dateFormatLocale: "ja-JP",
      btnSeeMore: "",
      btnReload: "",
      defaultTheme: "dark",
    });
  });

  onDestroy(() => {
    if (originalFetch) {
      window.fetch = originalFetch;
    }
  });
</script>

<div id="mt-container" class="mt-container w-full !bg-base-100 rounded-md">
  <div class="mt-body !p-0 m-0" role="feed">
    <div class="flex justify-center">
      <span class="loading loading-spinner loading-lg mt-5"></span>
    </div>
  </div>
</div>
