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

<style>
  article {
    padding: 0;
    &:hover {
      filter: brightness(110%);
      transition: all 0.3s ease;
    }

    & a {
      --pico-color: #c2c7d0;
      text-decoration: none;
      height: 100%;
      border-radius: 5px;
      display: flex;
      flex-flow: column;
      &:hover, &:active, &:focus {
        --pico-color: #c2c7d0;
      }
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
      aspect-ratio: 1;
      object-fit: cover;
      border-radius: 5px 5px 0px 0px;
      background-color: hsl(205, 20%, 90%);
    }
  }
</style>
