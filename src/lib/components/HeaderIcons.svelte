<script lang="ts">
  import { goto } from "$app/navigation";
  import {
    toggleModal,
    closeWithClickOutside,
    closeWithEscapeKey,
  } from "$lib/modal";
  import { onMount } from "svelte";

  function handleSubmit(event: SubmitEvent) {
    event.preventDefault();
    const searchField = document.getElementById("search") as HTMLInputElement;
    localStorage.setItem("search", searchField.value);
    toggleModal(event);
    const formData = new FormData(event.target as HTMLFormElement);
    const query = formData.get("search") as string;
    goto(`/search?q=${encodeURIComponent(query)}`);
  }

  onMount(() => {
    const searchField = document.getElementById("search") as HTMLInputElement;
    const searchButton = document.getElementById("search-button");

    searchButton?.addEventListener("click", async (event) => {
      if (localStorage.getItem("search")) {
        searchField.value = localStorage.getItem("search")!;
      }
      toggleModal(event);
    });

    const searchCloseButton = document.getElementById("search-close-button");
    searchCloseButton?.addEventListener("click", async (event) => {
      toggleModal(event);
    });

    closeWithClickOutside();
    closeWithEscapeKey();
  });
</script>

<div class="header-icons">
  <a href="/" class="contrast" role="button" aria-label="Home">
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="currentColor"
      class="bi bi-house"
      viewBox="0 0 16 16"
    >
      <path
        d="M8.707 1.5a1 1 0 0 0-1.414 0L.646 8.146a.5.5 0 0 0 .708.708L2 8.207V13.5A1.5 1.5 0 0 0 3.5 15h9a1.5 1.5 0 0 0 1.5-1.5V8.207l.646.647a.5.5 0 0 0 .708-.708L13 5.793V2.5a.5.5 0 0 0-.5-.5h-1a.5.5 0 0 0-.5.5v1.293L8.707 1.5ZM13 7.207V13.5a.5.5 0 0 1-.5.5h-9a.5.5 0 0 1-.5-.5V7.207l5-5 5 5Z"
      />
    </svg>
  </a>
  <button
    class="contrast"
    data-target="search-modal"
    id="search-button"
    aria-label="Search"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="currentColor"
      class="bi bi-search"
      viewBox="0 0 16 16"
    >
      <path
        d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001c.03.04.062.078.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1.007 1.007 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0z"
      />
    </svg>
  </button>
</div>

<dialog id="search-modal">
  <article>
    <button
      aria-label="Close"
      class="close"
      id="search-close-button"
      data-target="search-modal"
    />
    <form data-target="search-modal" on:submit={handleSubmit}>
      <input
        type="search"
        name="search"
        id="search"
        placeholder="Search"
        aria-label="Search"
      />
    </form>
  </article>
</dialog>
