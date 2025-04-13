<script lang="ts">
  import "$lib/assets/stylesheets/app.scss";
  import Footer from "$lib/components/Footer.svelte";

  const { children } = $props();
</script>

{@render children()}

<Footer></Footer>
