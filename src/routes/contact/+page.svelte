<script lang="ts">
  import { enhance } from "$app/forms";
  import type { ActionData } from "./$types";
  import Header from "$lib/components/Header.svelte";
  import OpenGraph from "$lib/components/OpenGraph.svelte";

  const { form }: { form: ActionData } = $props();

  let isLoading = $state(false);

  const url =
    "https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/ca63de19468e45b2833ebf325dbfd56c/Microsoft-Fluentui-Emoji-3d-Cat-3d.1024.png";
  const title = "お問い合わせ";
  const description = "お気軽にお問い合わせください";
</script>

<OpenGraph {url} {title} body={description}></OpenGraph>

<Header {url} {title}></Header>

<main class="container px-3 lg:px-12 py-10 mx-auto">
  <form
    method="POST"
    id="form"
    use:enhance={() => {
      isLoading = true;
      return async ({ update }) => {
        await update();
        isLoading = false;
      };
    }}
  >
    <div
      class="flex flex-col gap-y-7 lg:flex-row lg:flex-wrap justify-center lg:w-[80%] mx-auto"
    >
      {#if form?.errors && Object.keys(form.errors).length > 0}
        <div role="alert" class="alert alert-error">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-6 w-6 shrink-0 stroke-current"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <ul>
            {#each Object.values(form.errors) as error}
              <li>{error}</li>
            {/each}
          </ul>
        </div>
      {:else if form}
        <div role="alert" class="alert alert-success">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-6 w-6 shrink-0 stroke-current"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span
            >お問い合わせを受け付けました。ご連絡までしばらくお待ちください。</span
          >
        </div>
      {/if}
      <div class="flex w-full flex-col lg:flex-row gap-y-2 lg:gap-y-0">
        <label for="name" class="form-control w-full lg:w-[48%] lg:mr-[2%]">
          <div class="label">
            <span class="label-text text-lg">氏名</span>
          </div>
          <input
            type="text"
            id="name"
            name="name"
            class="input input-bordered w-full"
            required
          />
        </label>
        <label for="email" class="form-control w-full lg:w-[48%] lg:ml-[2%]">
          <div class="label">
            <span class="label-text text-lg">メールアドレス</span>
          </div>
          <input
            type="email"
            id="email"
            name="email"
            class="input input-bordered w-full"
            required
          />
        </label>
      </div>
      <label class="form-control w-full">
        <div class="label">
          <span class="label-text text-lg">本文</span>
        </div>
        <textarea
          name="text"
          class="textarea textarea-bordered w-full"
          rows="5"
          required
        ></textarea>
      </label>
      <div class="form-control hidden" aria-hidden="true">
        <label for="im-robot" class="label cursor-pointer">
          <span class="label-text">私はロボットです</span>
          <input
            name="im-robot"
            id="im-robot"
            type="checkbox"
            class="checkbox"
            value="true"
          />
        </label>
      </div>
      {#if isLoading}
        <button
          class="btn btn-neutral w-full lg:w-[50%] block mx-auto mt-5"
          type="submit"
          disabled
        >
          <div class="flex justify-center items-center gap-x-2">
            <span class="loading loading-spinner"></span>
            送信中…
          </div>
        </button>
      {:else}
        <button
          class="btn btn-neutral w-full lg:w-[80%] block mx-auto"
          type="submit">送信</button
        >
      {/if}
    </div>
  </form>
</main>
