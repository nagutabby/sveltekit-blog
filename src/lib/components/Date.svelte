<script lang="ts">
  export let createdAt: string;
  export let revisedAt: string

  function parseDate(date: string | Date) {
    const options: Intl.DateTimeFormatOptions = {
      year: "numeric",
      month: "short",
      day: "numeric"
    }
    date = new Date(date).toLocaleDateString("ja-JP", options);
    return date;
  }
</script>

<div class="date-wrapper">
  <p class="date-item">
    <svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" fill="currentColor" class="bi bi-calendar" viewBox="0 0 16 16">
      <path d="M3.5 0a.5.5 0 0 1 .5.5V1h8V.5a.5.5 0 0 1 1 0V1h1a2 2 0 0 1 2 2v11a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V3a2 2 0 0 1 2-2h1V.5a.5.5 0 0 1 .5-.5zM1 4v10a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V4H1z"/>
    </svg>
    <span>投稿日: {parseDate(createdAt)}</span>
  </p>
  <p class="date-item">
    <svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" fill="currentColor" class="bi bi-arrow-clockwise" viewBox="0 0 16 16">
      <path fill-rule="evenodd" d="M8 3a5 5 0 1 0 4.546 2.914.5.5 0 0 1 .908-.417A6 6 0 1 1 8 2v1z"/>
      <path d="M8 4.466V.534a.25.25 0 0 1 .41-.192l2.36 1.966c.12.1.12.284 0 .384L8.41 4.658A.25.25 0 0 1 8 4.466z"/>
    </svg>
    <span>最終更新日: {parseDate(revisedAt)}</span>
  </p>
</div>

<style lang="scss">
  @use "@picocss/pico/scss/_settings" as *;

  .date-wrapper {
    display: flex;
    flex-flow: column;
  }

  .date-item {
    display: flex;
    align-items: center;
    column-gap: 0.3rem;
  }

  @media (map-get($breakpoints, sm, breakpoint) <= width) {
    .date-wrapper {
      flex-flow: row nowrap;
      column-gap: 1rem;
    }
  }
</style>
