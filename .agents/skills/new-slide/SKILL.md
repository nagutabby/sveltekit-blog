---
name: new-slide
description: このブログのプレゼンスライドを新規作成するときに使う。HTML本文を下書きし、レイアウト検証とPDF化を行う。
---

# プレゼンスライドの新規作成

`web/static/content/slides/<id>.html` を作り、可能なら同名のPDFを書き出す。スライドの見出し・構成・本文はCodexが下書きする。`/slides` への掲載にはPDFが必要。

## 必要な資料

- 作文時は [writing.md](references/writing.md) と [layout.md](references/layout.md) を読む。
- 図表を組むとき、画像・スクリーンショットを扱うときは [diagrams.md](references/diagrams.md) を読む。
- 日本語フォントを同梱するときは [fonts.md](references/fonts.md) を読む。
- `web/static/content/templates/slide.html` は現行のデザイン仕様。HTMLのコピー元として読み、CSS・レイアウト・budoux-ja読み込みを維持する。

## 規約

- `id` は内容を表す日付なしのkebab-caseスラッグ。`web/static/content/slides/` 直下のHTML・PDFと重複させない。記事・書評のidとの重複は問わない。既存ファイルを上書きしない。
- 1枚目は表紙、最終ページは参考文献。間の種別は内容に合わせて選び、全種類を使おうとしない。
- 表のキャプションは `<caption>` で上、図のキャプションは `<figcaption>` で下に置く。`figure img` は `object-fit: contain` を維持する。
- 実在しない数値・固有名詞・参考文献・画像パスを作らない。実写真が必要で手元に無ければユーザーに提供を依頼する。装飾アイコンが目的なら図を省く選択肢を優先する。

## 手順

1. テーマと発表タイトル案、想定聴衆、伝えたい要点、具体的な数値・コード例・比較対象、参考文献を確認する。すでに分かっていることは聞き直さない。事実関係が曖昧なら推測で埋めない。
2. `web/static/content/slides/` のファイル名を調べ、発表タイトルと重複しないidを提案してユーザーに確認する。
3. テンプレートと必要な参照資料を読み、構成と本文を下書きする。表紙、1・2カラム、比較、コード、表、フロー図、高橋メソッド風、参考文献などから内容に合う形式を選ぶ。テンプレートのコンポーネントを使い、`<title>` と表紙の `<h1>` を確定したタイトルにする。ページ番号の既存スクリプトには手を加えない。
4. 完成したHTML全体を `<id>.html` として新規保存する。
5. `node web/scripts/validate-slide-layout.mjs "web/static/content/slides/<id>.html"` を実行する。1920×1080pxからのはみ出しと見出し・本文間隔の不統一があれば、文章量・画像サイズ・配置を直して再実行する。Chrome/Chromiumが無く実行できなければ、その理由を完了報告に記す。
6. レイアウトが `OK` なら `node web/scripts/export-slide-pdf.mjs "web/static/content/slides/<id>.html"` を実行する。ローカルにGhostscriptがあれば、スクリプトがPDFの `/MediaBox` を `[0 0 1920 1080]` に補正する。Chrome/Chromiumが無くPDF化できなければ、手元での実行コマンドを伝える。
7. HTMLとPDFのパス、レイアウト検証とPDF化の結果を報告する。本文はCodexの下書きなので、数値・固有名詞・参考文献の確認をユーザーに依頼する。

スライドHTMLの構造を検証する専用スクリプトはない。`validate-content.sh` は記事・書評のfrontmatter専用であり、スライドには使わない。
