<script lang="ts">
  type Props = {
    publishedAt: Date;
  }
  const { publishedAt }: Props = $props();

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

<div class="flex gap-x-4 items-center">
  <p class="flex gap-x-1 justify-center items-center !m-0">
    <svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" fill="currentColor" class="bi bi-calendar" viewBox="0 0 16 16">
      <path d="M3.5 0a.5.5 0 0 1 .5.5V1h8V.5a.5.5 0 0 1 1 0V1h1a2 2 0 0 1 2 2v11a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V3a2 2 0 0 1 2-2h1V.5a.5.5 0 0 1 .5-.5zM1 4v10a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V4H1z"/>
    </svg>
    <span>{parseDate(publishedAt)}</span>
  </p>
</div>
