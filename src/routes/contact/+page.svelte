<script lang="ts">
  import { onMount } from "svelte";
  import type { ActionData } from "./$types";
  import { enhance } from "$app/forms";
  import Breadcrumb from "$lib/components/Breadcrumb.svelte";

  export let form: ActionData;

  onMount(async () => {
    const form = document.getElementById("form") as HTMLFormElement;
    const submitButtonField = document.getElementById(
      "submit-button"
    ) as HTMLInputElement;
    submitButtonField?.addEventListener("click", async (event) => {
      if (form.checkValidity()) {
        submitButtonField.disabled = true;
        submitButtonField.innerHTML = "送信中…";
        submitButtonField.setAttribute("aria-busy", "true");
        form.submit();
      }
    });
  });
</script>

<Breadcrumb title="お問い合わせ" />

<form method="POST" id="form" use:enhance>
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
  <button type="submit" id="submit-button">送信</button>
  {#if form}
    {#if form.isFailed}
      <p>メールを送信できませんでした…😥</p>
    {:else}
      <p>メールを送信しました！✅</p>
    {/if}
  {/if}
</form>
