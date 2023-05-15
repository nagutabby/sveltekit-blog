<script lang="ts">
  export let url: string;
  export let image: string;
  export let title: string;
  export let date: string | Date;

  date = new Date(date);
  const year = date.getFullYear();
  const month = ("00" + (date.getMonth() + 1)).slice(-2);
  const day = ("00" + date.getDate()).slice(-2);
  date = year + "-" + month + "-" + day;
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
      text-decoration: none;
      padding: var(--block-spacing-vertical) var(--block-spacing-horizontal);
      display: flex;
      flex-flow: column;
      & .description {
        text-align: center;
        margin: 1rem 0 0 0;
        display: flex;
        flex-flow: column;
        gap: 1rem;
        & p {
          margin: 0;
        }
      }
    }
    & img {
      margin: 0 auto;
      min-height: 250px;
      min-width: 250px;
      aspect-ratio: 1 / 1;
      object-fit: cover;
      border-radius: 4%;
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
