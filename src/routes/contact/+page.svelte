<script lang="ts">
  import { enhance } from "$app/forms";
  import Breadcrumb from "$lib/components/Breadcrumb.svelte";
  import type { ActionData } from "./$types";

  export let form: ActionData;

  let isLoading = false;
  let hasError = false;
</script>

<Breadcrumb title="お問い合わせ" />

<form
  method="POST"
  id="form"
  use:enhance={() => {
    isLoading = true;
    return async ({ result, update }) => {
      if (result.type === "failure") {
        hasError = true;
      }
      await update();
      isLoading = false;
    };
  }}
>
  <div class="grid">
    <label for="name">
      お名前
      <input
        type="text"
        id="name"
        name="name"
        placeholder="山田 太郎"
        required
      />
    </label>
    <label for="email">
      メールアドレス
      <input
        type="email"
        id="email"
        name="email"
        placeholder="user@example.com"
        required
      />
    </label>
  </div>
  <label for="text">
    お問い合わせ内容
    <textarea
      id="text"
      name="text"
      placeholder="ご要件をご記入ください"
      rows="5"
      required
    />
  </label>
  {#if isLoading}
    <button type="submit" aria-busy="true" disabled>送信中…</button>
  {:else}
    <button type="submit">送信</button>
  {/if}
  {#if form}
    <p>以下の内容でメールを送信しました！✅</p>
    <p>氏名: {form.name}</p>
    <p>メールアドレス: {form.email}</p>
    <p>本文: {form.text}</p>
  {/if}
  {#if hasError}
    <p>メールを送信できませんでした…😥</p>
  {/if}
</form>
