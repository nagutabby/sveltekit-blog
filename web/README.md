# Astro / Svelte frontend

`web/` is a statically generated Astro site. Astro renders pages and shared
metadata, while Svelte hydrates interactive features such as search, contact,
PDF viewing, sharing, and the Mastodon timeline. Article and review Markdown
continues to be read from `../backend/content` during the build.

Run `pnpm dev` for local development, `pnpm check` for Astro and Svelte type
checks, `pnpm test` for unit tests, and `pnpm build` to generate `dist/`.
