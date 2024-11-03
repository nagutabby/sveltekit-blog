<script lang="ts">
  import "../app.scss";
  import { setupIonicBase } from "ionic-svelte";
  import { IonPage } from "ionic-svelte";
  import "ionic-svelte/components/all";
  import "../variables.scss";
  import { page } from "$app/stores";
  import Header from "$lib/components/Header.svelte";
  import Footer from "$lib/components/Footer.svelte";
  import OpenGraph from "$lib/components/OpenGraph.svelte";

  setupIonicBase();
</script>

<ion-app>
  <IonPage>
    <ion-content>
      {#if $page.url.pathname === "/" || $page.url.pathname === "/contact" || $page.url.pathname === "/search"}
        <OpenGraph
          url={$page.data.image.url}
          title={$page.data.title}
          description={$page.data.description}
        ></OpenGraph>

        <Header
          url={$page.data.image.url}
          title={$page.data.title}
          description={$page.data.description}
        ></Header>

        <main class="container">
          <slot />
        </main>

        <Footer></Footer>
      {:else}
        <slot />
      {/if}
    </ion-content>
  </IonPage>
</ion-app>
