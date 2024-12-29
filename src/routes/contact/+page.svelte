<script lang="ts">
  import { enhance } from "$app/forms";
  import type { ActionData } from "./$types";

  const { form }: { form: ActionData } = $props();

  let isLoading = $state(false);
  let isError = $state(false);
</script>

<form
  method="POST"
  id="form"
  use:enhance={() => {
    isLoading = true;
    return async ({ result, update }) => {
      if (result.type === "failure") {
        isError = true;
      }
      await update();
      isLoading = false;
    };
  }}
>
  <div class="flex flex-col gap-y-5 lg:flex-row flex-wrap justify-center">
    <label
      class="input input-bordered flex items-center gap-2 w-full lg:w-[43%] lg:mr-[2%]"
    >
      氏名
      <input type="text" name="name" class="grow" required />
    </label>
    <label
      class="input input-bordered flex items-center gap-2 w-full lg:w-[43%] lg:ml-[2%]"
    >
      メールアドレス
      <input type="email" name="email" class="grow" required />
    </label>

    <label class="form-control w-full lg:w-[90%]">
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
        class="btn btn-neutral w-full lg:w-[90%] block mx-auto"
        type="submit">送信</button
      >
    {/if}
  </div>
  {#if form}
    <div class="w-full lg:w-[90%] mx-auto mt-5">
      <p>以下の内容でメールを送信しました！✅</p>
      <p>氏名: {form.name}</p>
      <p>メールアドレス: {form.email}</p>
      <p>本文: {form.text}</p>
    </div>
  {/if}
  {#if isError}
    <p>メールを送信できませんでした…😥</p>
  {/if}
</form>
