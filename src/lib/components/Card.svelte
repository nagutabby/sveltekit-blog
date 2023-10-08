<script lang="ts">
  export let url: string;
  export let image: string;
  export let title: string;
  export let date: string;
</script>

<article>
  <a href={url}>
    <img
      alt=""
      height="250"
      width="250"
      srcset="{image}?fm=webp&fit=clip&w=250&h=250 250w, {image}?fm=webp&fit=clip&w=500&h=500 500w"
      loading="lazy"
    />
    <div class="description">
      <p>{title}</p>
      <p>{date}</p>
    </div>
  </a>
</article>

<style>
  article {
    transition: all 0.2s ease;
    padding: 0;

    & a {
      height: 100%;
      text-decoration: none;
      border-radius: 5px;
      display: flex;
      flex-flow: column;
      & .description {
        text-align: center;
        display: flex;
        flex-flow: column;
        width: 90%;
        margin: 1rem auto;
      }
    }
    & img {
      width: 100%;
      aspect-ratio: 3 / 2;
      object-fit: cover;
      border-radius: 5px 5px 0px 0px;
      background-color: hsl(205, 20%, 94%);
    }
  }

  @media (prefers-color-scheme: light) {
    article:hover {
      filter: brightness(90%);
    }
    p {
      color: hsl(205, 20%, 32%);
    }
  }
  @media (prefers-color-scheme: dark) {
    article:hover {
      filter: brightness(110%);
    }
    p {
      color: hsl(205, 16%, 77%);
    }
  }
</style>
