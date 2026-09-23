# スライドの構成・検証

## テンプレートの使い方

`web/static/content/templates/slide.html` をHTML全体のコピー元にする。CSS・レイアウト・budoux-ja読み込みを変更しない。

使えるスライド種別は、表紙 (`.slide--cover`)、1カラム (`.cols-1`)、2カラム (`.cols-2`)、比較 (`.compare`)、コード (`.code-block`)、表 (`table.simple-table`)、フロー図 (`.flow-diagram`/`.flow-chain`/`.flow-compare`/`.flow-converge`/`.bar-chart`)、カード代替 (`.minimal-list`/`.statement-compare`/`.chat-thread`/`.sticky-board`)、高橋メソッド風 (`.slide--impact`)、参考文献 (`ol.references`)。1枚目は表紙、最終ページは参考文献とし、他は内容に必要なものだけを選ぶ。`.slide--impact` は1デッキに1〜2回まで。

見出し直後の本文配置には次のクラスを使う。これらの `gap` は内容ごとに異なるため、必要なら `style="gap:Npx"` で指定する。`.slide-inner`/`.cols-1` の `gap` と別階層なので二重加算にはならない。

- `.content-center`: 図解や短い強調ステートメントを上下左右中央に置く
- `.content-middle`: 本文を左寄せのまま上下中央に置く
- `.content-start`: リード文と箇条書きなどを上から順に読むため、上詰め・左寄せにする

自己紹介の名前や締めの一言など、見出し以外の強調には `.emphasis-text` を使う。フォントの太さや行高を都度インライン指定しない。ページ番号はテンプレートの既存スクリプトに任せる。

実写真やスクリーンショットを使う場合、画像が無ければユーザーに提供を依頼する。実在しない画像パスを埋めない。画像内に第三者の実名があれば公開前に加工を確認する。抽象的な図や実画面の代替表現は [diagrams.md](diagrams.md) を参照する。テンプレートのdata URLによるサンプルSVGは、実画像が無ければプレースホルダーとして残せる。

## レイアウト検証

`node web/scripts/validate-slide-layout.mjs "web/static/content/slides/<id>.html"` は `puppeteer-core` でローカルのChrome/Chromiumを使い、各 `.slide` が1920×1080pxに収まるかを実測する。budoux-jaとWebフォント適用後の改行も反映される。文章量、フォントサイズ、画像サイズを変更したら再実行する。

`NG` なら該当スライドの文章量・画像サイズを調整して再実行する。見出しと本文の間隔が不統一なら、`.slide-inner`/`.cols-1`/`.cols-2` 直下の子要素のインライン `margin` と `gap` の二重加算を確認する。Chrome/Chromiumが無く実行できない場合は理由を完了報告に明示する。記事・書評用の `validate-content.sh` はHTML構造を検証しない。

## PDF化

レイアウト検証が `OK` になったら、`node web/scripts/export-slide-pdf.mjs "web/static/content/slides/<id>.html"` で同名のPDFを出力する。テンプレートの `@page { size: 1920px 1080px; margin: 0 }` と `preferCSSPageSize: true` を使い、各 `.slide` を1ページにする。

ChromiumはCSSのpxをPDFのptへ換算するため、Ghostscriptが無い場合、PDFの `/MediaBox` は `[0 0 1440 810]` になる。スクリプトはGhostscript (`gs`) があれば `-dPDFFitPage` で内容をベクターのまま拡大し、`/MediaBox` を `[0 0 1920 1080]` に補正する。Ghostscriptが無くてもPDF化自体は可能。`.slide` が無い、ファイルが無い、補正に失敗したなどの `NG` は修正して再実行する。Chrome/Chromiumが無くPDF化できない場合は理由と実行コマンドをユーザーに伝える。
