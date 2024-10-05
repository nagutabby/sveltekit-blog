<script lang="ts">
  import { goto } from "$app/navigation";

  function handleSubmit(event: SubmitEvent) {
    event.preventDefault();
    const formData = new FormData(event.target as HTMLFormElement);
    const query = formData.get("search") as string;
    goto(`/search?q=${encodeURIComponent(query)}`);
  }
</script>

<form on:submit={handleSubmit}>
  <input type="search" name="search" placeholder="Search" aria-label="Search" />
</form>

<style lang="scss">
  @use "@picocss/pico/scss/_settings" as *;

  form {
    display: flex;
    flex-flow: column;
    position: absolute;
    padding-left: inherit;
    padding-right: inherit;
    top: -10vh;
    left: 0;
    width: 100%;
  }

  @media (map-get($breakpoints, md, breakpoint) <=width) {
    form {
      display: flex;
      flex-flow: row;
      column-gap: 3%;
    }
  }
</style>
