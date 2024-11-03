<script lang="ts">
  export let url: string;
  export let image: string;
  export let title: string;
</script>

<article>
  <a href="/{url}">
    <img
      alt=""
      height="250"
      width="250"
      srcset="{image}?fm=webp&fit=clip&w=250&h=250 250w, {image}?fm=webp&fit=clip&w=500&h=500 500w"
      loading="lazy"
    />
    <div class="description">
      <p>{title}</p>
    </div>
  </a>
</article>
