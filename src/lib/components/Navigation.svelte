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

<div class="flex sm:gap-x-4 overflow-x-auto nav-scrollbar">
  <a
    href="/"
    class="btn md:btn-lg btn-ghost flex gap-x-2"
    role="button"
    aria-label="ホーム"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="currentColor"
      class="bi bi-house h-1/3 h-[1em]"
      viewBox="0 0 16 16"
    >
      <path
        d="M8.707 1.5a1 1 0 0 0-1.414 0L.646 8.146a.5.5 0 0 0 .708.708L2 8.207V13.5A1.5 1.5 0 0 0 3.5 15h9a1.5 1.5 0 0 0 1.5-1.5V8.207l.646.647a.5.5 0 0 0 .708-.708L13 5.793V2.5a.5.5 0 0 0-.5-.5h-1a.5.5 0 0 0-.5.5v1.293L8.707 1.5ZM13 7.207V13.5a.5.5 0 0 1-.5.5h-9a.5.5 0 0 1-.5-.5V7.207l5-5 5 5Z"
      />
    </svg>
    <span>ホーム</span>
  </a>
  <a
    href="/slides"
    class="btn md:btn-lg btn-ghost flex gap-x-2"
    role="button"
    aria-label="スライド"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="currentColor"
      class="bi bi-archive h-[1em]"
      viewBox="0 0 16 16"
    >
      <path
        d="M0 2a1 1 0 0 1 1-1h14a1 1 0 0 1 1 1v2a1 1 0 0 1-1 1v7.5a2.5 2.5 0 0 1-2.5 2.5h-9A2.5 2.5 0 0 1 1 12.5V5a1 1 0 0 1-1-1zm2 3v7.5A1.5 1.5 0 0 0 3.5 14h9a1.5 1.5 0 0 0 1.5-1.5V5zm13-3H1v2h14zM5 7.5a.5.5 0 0 1 .5-.5h5a.5.5 0 0 1 0 1h-5a.5.5 0 0 1-.5-.5"
      />
    </svg>
    <span>スライド</span>
  </a>
  <a
    href="/reviews"
    class="btn md:btn-lg btn-ghost flex gap-x-2"
    role="button"
    aria-label="レビュー"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="currentColor"
      class="bi bi-book h-[1em]"
      viewBox="0 0 16 16"
    >
      <path
        d="M1 2.828c.885-.37 2.154-.769 3.388-.893 1.33-.134 2.458.063 3.112.752v9.746c-.935-.53-2.12-.603-3.213-.493-1.18.12-2.37.461-3.287.811zm7.5-.141c.654-.689 1.782-.886 3.112-.752 1.234.124 2.503.523 3.388.893v9.923c-.918-.35-2.107-.692-3.287-.81-1.094-.111-2.278-.039-3.213.492zM8 1.783C7.015.936 5.587.81 4.287.94c-1.514.153-3.042.672-3.994 1.105A.5.5 0 0 0 0 2.5v11a.5.5 0 0 0 .707.455c.882-.4 2.303-.881 3.68-1.02 1.409-.142 2.59.087 3.223.877a.5.5 0 0 0 .78 0c.633-.79 1.814-1.019 3.222-.877 1.378.139 2.8.62 3.681 1.02A.5.5 0 0 0 16 13.5v-11a.5.5 0 0 0-.293-.455c-.952-.433-2.48-.952-3.994-1.105C10.413.809 8.985.936 8 1.783"
      />
    </svg>
    <span>レビュー</span>
  </a>
  <button
    class="btn md:btn-lg btn-ghost flex gap-x-2"
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
      class="bi bi-search h-[1em]"
      viewBox="0 0 16 16"
    >
      <path
        d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001c.03.04.062.078.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1.007 1.007 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0z"
      />
    </svg>
    <span>検索</span>
  </button>
</div>

<dialog use:portal class="modal" id="search-modal">
  <div class="modal-box flex flex-col gap-y-5">
    <p class="font-bold text-2xl">検索</p>
    <form class="flex" onsubmit={handleSubmit}>
      <label
        for="search"
        class="input input-bordered flex items-center gap-2 w-full"
      >
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
