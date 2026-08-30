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
    // ルート+layout.tsでprerender = trueを指定して全ページSSG化しているため、
    // サーバーランタイムを持たないadapter-staticで完全な静的サイトとして出力する。
    adapter: adapter(),
    // Header.css/Date.cssのような極小コンポーネントCSSのみを<head>にインライン化する。
    // Tailwind/DaisyUIの本体バンドルやKaTeX CSSはこの閾値を超えるため対象外で、
    // それらはビルド後のscripts/inline-critical-css.mjsでクリティカルCSS抽出する。
    inlineStyleThreshold: 1000,
    paths: {
      // 既定の相対パス(../_app/...)だとbeastiesがCSSファイルをbuild/配下の
      // ものと解決できずスキップしてしまう(root直下のページ以外で無効化される)。
      // ルート絶対パス(/_app/...)に統一して解決できるようにする。サイトは
      // blog.nagutabby.ukのドメインルートで配信されるためサブパス配信の制約はない。
      relative: false
    },
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
