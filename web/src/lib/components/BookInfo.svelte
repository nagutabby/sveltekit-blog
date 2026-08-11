<script lang="ts">
  type Props = {
    title: string;
    description: string;
    jpECode: string;
  };

  const { title, description, jpECode }: Props = $props();
</script>

<div class="card sm:w-[90%] xl:w-[70%] sm:mx-auto bg-base-100 card-md shadow-sm">
  <div class="card-body">
    <div class="card-title">{title}</div>
    <p>{description}</p>
    <div class="justify-center card-actions">
      <a
        href={`https://www.books.or.jp/book-details/${jpECode}`}
        target="_blank"
        class="btn btn-primary w-full sm:w-[60%]">この本の詳細</a
      >
    </div>
  </div>
</div>
