<script lang="ts">
  import "$lib/assets/stylesheets/app.scss";
  import Footer from "$lib/components/Footer.svelte";
  import { onNavigate } from "$app/navigation";

  onNavigate((navigation) => {
    if (!document.startViewTransition) return;

    return new Promise((resolve) => {
      document.startViewTransition(async () => {
        resolve();
        await navigation.complete;
      });
    });
  });

  const { children } = $props();
</script>

{@render children()}

<Footer></Footer>
