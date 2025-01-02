<script lang="ts">
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  const portal = (element: HTMLElement) => {
    document.body.appendChild(element);

    return {
      destroy() {
        element.remove();
      },
    };
  };

  const handleSubmit = (event: SubmitEvent) => {
    event.preventDefault();
    const searchField = document.getElementById("search") as HTMLInputElement;
    const searchModal = document.getElementById(
      "search-modal",
    ) as HTMLDialogElement;
    if (searchModal) {
      searchModal.close();
    }
    localStorage.setItem("search", searchField.value);
    goto(`/search?q=${encodeURIComponent(searchField.value)}`);
  };

  onMount(() => {
    const searchField = document.getElementById("search") as HTMLInputElement;
    const searchModal = document.getElementById("search-modal");

    searchModal?.addEventListener("click", async () => {
      if (localStorage.getItem("search")) {
        searchField.value = localStorage.getItem("search")!;
      }
    });
  });
</script>

<div class="flex justify-evenly">
  <a
    href="/"
    class="btn md:btn-lg btn-ghost btn-circle"
    role="button"
    aria-label="ホーム"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="currentColor"
      class="bi bi-house h-1/2"
      viewBox="0 0 16 16"
    >
      <path
        d="M8.707 1.5a1 1 0 0 0-1.414 0L.646 8.146a.5.5 0 0 0 .708.708L2 8.207V13.5A1.5 1.5 0 0 0 3.5 15h9a1.5 1.5 0 0 0 1.5-1.5V8.207l.646.647a.5.5 0 0 0 .708-.708L13 5.793V2.5a.5.5 0 0 0-.5-.5h-1a.5.5 0 0 0-.5.5v1.293L8.707 1.5ZM13 7.207V13.5a.5.5 0 0 1-.5.5h-9a.5.5 0 0 1-.5-.5V7.207l5-5 5 5Z"
      />
    </svg>
  </a>
  <a
    href="/slides"
    class="btn md:btn-lg btn-ghost btn-circle"
    role="button"
    aria-label="スライド"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="currentColor"
      class="bi bi-archive h-1/2"
      viewBox="0 0 16 16"
    >
      <path
        d="M0 2a1 1 0 0 1 1-1h14a1 1 0 0 1 1 1v2a1 1 0 0 1-1 1v7.5a2.5 2.5 0 0 1-2.5 2.5h-9A2.5 2.5 0 0 1 1 12.5V5a1 1 0 0 1-1-1zm2 3v7.5A1.5 1.5 0 0 0 3.5 14h9a1.5 1.5 0 0 0 1.5-1.5V5zm13-3H1v2h14zM5 7.5a.5.5 0 0 1 .5-.5h5a.5.5 0 0 1 0 1h-5a.5.5 0 0 1-.5-.5"
      />
    </svg>
  </a>
  <button
    class="btn md:btn-lg btn-ghost btn-circle"
    aria-label="Search"
    onclick={() => {
      const searchModal = document.getElementById(
        "search-modal",
      ) as HTMLDialogElement;
      searchModal.showModal();
    }}
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="currentColor"
      class="bi bi-search h-1/2"
      viewBox="0 0 16 16"
    >
      <path
        d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001c.03.04.062.078.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1.007 1.007 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0z"
      />
    </svg>
  </button>
</div>

<dialog use:portal class="modal" id="search-modal">
  <div class="modal-box flex flex-col gap-y-5">
    <p class="font-bold text-2xl">検索</p>
    <form class="flex" onsubmit={handleSubmit}>
      <label for="search" class="input input-bordered flex items-center gap-2 w-full">
        <input type="text" class="grow" id="search" required />
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 16 16"
          fill="currentColor"
          class="h-4 w-4 opacity-70"
        >
          <path
            fill-rule="evenodd"
            d="M9.965 11.026a5 5 0 1 1 1.06-1.06l2.755 2.754a.75.75 0 1 1-1.06 1.06l-2.755-2.754ZM10.5 7a3.5 3.5 0 1 1-7 0 3.5 3.5 0 0 1 7 0Z"
            clip-rule="evenodd"
          />
        </svg>
      </label>
    </form>
  </div>
  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>
