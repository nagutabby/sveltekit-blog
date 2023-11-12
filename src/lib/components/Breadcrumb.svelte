<script lang="ts">
  import { page } from "$app/stores";
  export let titles: string[];
  const partsOfPathname = $page.url.pathname.substring(0).split("/");
  let tmpURL: string;

  let rawBreadcrumbElements: string[] = [];

  partsOfPathname.forEach((_, i) => {
    if (i === partsOfPathname.length - 1) {
      rawBreadcrumbElements.push(`<li>${titles[i]}</li>`);
    } else {
      tmpURL = "";
      partsOfPathname.slice(0, i + 1).forEach((partOfPathname) => {
        tmpURL = `${tmpURL}${partOfPathname}/`;
      });
      rawBreadcrumbElements.push(`<li><a href=${tmpURL}>${titles[i]}</a></li>`);
    }
  });
</script>

<nav aria-label="breadcrumb">
  <ul>
    {#each rawBreadcrumbElements as rawBreadcrumbElement}
      {@html rawBreadcrumbElement}
    {/each}
  </ul>
</nav>
