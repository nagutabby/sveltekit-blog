<script lang="ts">
  type Props = {
    title: string;
    description: string;
    infoLink: string;
  };

  const { title, description, infoLink }: Props = $props();
</script>

<div class="card sm:w-[90%] xl:w-[70%] sm:mx-auto bg-base-100 card-md shadow-sm">
  <div class="card-body">
    <div class="card-title">{title}</div>
    <p>{description}</p>
    <div class="justify-center card-actions">
      <a
        href={infoLink}
        target="_blank"
        class="btn btn-primary w-full sm:w-[60%]">Google Booksで見る</a
      >
    </div>
  </div>
</div>
