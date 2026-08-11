import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";

/** @type {import('@sveltejs/kit').Config} */
const config = {
  compilerOptions: {
    runes: true,
  },
  // Consult https://kit.svelte.dev/docs/integrations#preprocessors
  // for more information about preprocessors
  preprocess: [vitePreprocess({})],

  kit: {
    // 全ページprerender=true(SSG)なのでfallbackは不要。strict(既定true)のまま
    // 全ルートがprerenderされていることをビルド時に強制する。
    adapter: adapter({
      pages: 'build',
      assets: 'build',
    }),
    prerender: {
      // 記事本文中の壊れたリンク/画像参照が1件あるだけで全体のビルドが
      // 失敗しないようにする(警告は出るがビルドは継続する)。
      handleHttpError: 'warn',
      // /page/[page]・/reviews/page/[page]はentries()が記事/レビュー件数
      // から動的にページ番号を計算する。件数が少なく2ページ目が存在しない
      // 場合、entries()は空配列を返しクロールでも見つからないが、それは
      // 想定内の状態なのでビルドを失敗させない。
      handleUnseenRoutes: 'warn'
    }
  },
};

export default config;
