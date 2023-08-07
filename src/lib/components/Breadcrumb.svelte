<script lang="ts">
  import { page } from "$app/stores";
  const baseURL = "/";
  const partsOfPathname = $page.url.pathname.substring(1).split("/");
  let tmpURL: string;
  export let title: string;

  function clearTmpURL() {
    tmpURL = baseURL;
  }
  function addPathToTmpURL(partOfPathname: string) {
    tmpURL = tmpURL + "/" + partOfPathname;
  }
</script>

<nav aria-label="breadcrumb">
  <ul>
    <li>
      <a href={baseURL}>Home</a>
    </li>
    {#each partsOfPathname as partOfPathname, i}
      {#if i === partsOfPathname.length - 1}
        <li>
          {title}
        </li>
      {:else}
        {clearTmpURL()}
        <li>
          {#each partsOfPathname.slice(0, i) as partOfPathname}
            {addPathToTmpURL(partOfPathname)}
          {/each}
          <a href={tmpURL}>{partOfPathname}</a>
        </li>
      {/if}
    {/each}
  </ul>
</nav>
