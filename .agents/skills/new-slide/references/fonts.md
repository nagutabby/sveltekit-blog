## 日本語フォントの同梱(推奨)

テンプレートの`--font-sans`は`"Noto Sans JP", "Hiragino Sans", "Hiragino Kaku Gothic ProN", "Yu Gothic", ...`の順で、`vendor/NotoSansJP-Regular.woff2`/`-Bold.woff2`があれば最優先で使う`@font-face`が既に用意されている。Hiragino(macOS専用)・Yu Gothic(Windows同梱)はOS依存のため、ヘッドレスChromeでPDF化する環境(Linux CI等)によっては日本語フォントが1つも無く、文字化け・トウフ文字になるリスクがある。スライド本文が完成したら、以下の手順で実際に使われている文字だけをサブセット化し、`vendor/`に配置する(この2ファイルが無くても`@font-face`は静かにフォールバックするだけなので、省略しても壊れたスライドにはならないが、社外公開前には行うことを推奨する)。

前提: `fontTools`(`pip install fonttools`)と、woff2出力用の`brotli`(`pip install brotli`)。

1. Noto Sans JPのRegular/Bold(完全版のOTF/TTF)を入手する。[Google Fonts](https://fonts.google.com/noto/specimen/Noto+Sans+JP)からダウンロードするか、[notofonts/noto-cjk](https://github.com/notofonts/noto-cjk)リポジトリのJapanese向けStatic OTFを取得する(配布パスはリポジトリの更新で変わることがあるため、都度確認する)。
2. 完成したスライドHTMLから、実際に使われている文字を抽出する(タグを除去し、重複を除いた文字集合を1ファイルに書き出す)。
3. `pyftsubset`でサブセット化する(`--unicodes-file`ではなく、リテラル文字を渡す`--text-file`を使う点に注意):
   ```bash
   pyftsubset NotoSansJP-Regular.otf --text-file=chars.txt --unicodes="U+0020-007E" \
     --output-file=vendor/NotoSansJP-Regular.woff2 --flavor=woff2
   pyftsubset NotoSansJP-Bold.otf --text-file=chars.txt --unicodes="U+0020-007E" \
     --output-file=vendor/NotoSansJP-Bold.woff2 --flavor=woff2
   ```
4. `node web/scripts/validate-slide-layout.mjs "<slideのパス>"`を再実行し、フォント差し替え後もはみ出しが無いことを確認する(このスクリプトは`document.fonts.ready`を待ってから計測するため、Webフォント読み込み後のレイアウトを正しく検証できる)。
