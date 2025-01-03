<script lang="ts">
  import { enhance } from "$app/forms";
  import type { ActionData } from "./$types";

  const { form }: { form: ActionData } = $props();

  let isLoading = $state(false);
</script>

<form
  method="POST"
  id="form"
  use:enhance={() => {
    isLoading = true;
    return async ({ result, update }) => {
      await update();
      isLoading = false;
    };
  }}
>
  <div
    class="flex flex-col gap-y-5 lg:flex-row lg:flex-wrap justify-center lg:w-[80%] mx-auto"
  >
    {#if form?.errors}
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
        <span>{form.errors.imRobot}</span>
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
    <label
      class="input input-bordered flex items-center gap-2 w-full lg:w-[48%] lg:mr-[2%]"
    >
      氏名
      <input type="text" name="name" class="grow" required />
    </label>
    <label
      class="input input-bordered flex items-center gap-2 w-full lg:w-[48%] lg:ml-[2%]"
    >
      メールアドレス
      <input type="email" name="email" class="grow" required />
    </label>

    <label class="form-control w-full">
      <div class="label">
        <span class="label-text text-lg">内容</span>
      </div>
      <textarea
        name="text"
        class="textarea textarea-bordered w-full"
        rows="5"
        required
      ></textarea>
      <div class="label"></div>
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
