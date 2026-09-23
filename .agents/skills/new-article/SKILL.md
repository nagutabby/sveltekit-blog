---
name: new-article
description: このブログの記事を新規作成するときに使う。タイトル・ID・絵文字画像を決め、本文が空のMarkdownを作る。
---

# 記事の新規作成

`backend/content/articles/` に記事を1件作る。**本文は書かない。** frontmatterと空の本文だけを用意し、執筆はユーザーに任せる。

## 規約

- ファイル名は当日の `YYYY-MM-DD.md`。同日分があれば未使用の `YYYY-MM-DD-N.md` (`N` は2から)を使う。既存ファイルを上書きしない。
- `id` は英小文字・数字をハイフンでつないだスラッグで、`backend/content/articles/` と `backend/content/reviews/` の全体で一意にする。両ディレクトリの `id:` を確認する。
- frontmatterは次の4フィールドだけにする。本文や定型見出しは追加しない。

  ```yaml
  ---
  id: kebab-case-english-words
  title: 日本語タイトル
  image: images/Microsoft-Fluentui-Emoji-Flat-<Name>.512.png
  is_draft: true
  ---
  ```

- `is_draft: true` は、空の本文を公開しないために必須。ユーザーが本文を完成させた後に `false` に変更できる。
- `web/src/lib/utils.ts` の `getWebpPath` は `image` の拡張子を `.webp` に変換する。PNGと同名のWebPの両方が必要。
- `web/static/content/templates/article.md` は現行のfrontmatterと食い違うため参照しない。

## 手順

1. 記事のテーマを聞く。すでに分かっていることは聞き直さない。
2. 既存の `id:` を確認し、日本語タイトル、重複しないid、内容に合う [Fluent UI Emoji](https://github.com/microsoft/fluentui-emoji) の英語名を1案ずつ提示してユーザーに確認する。本文は提案しない。
3. 当日の未使用ファイル名を決める。確定した絵文字名で `bash .agents/skills/new-article/scripts/fetch-emoji-image.sh "White Flag"` を実行する（引数は実際の名前に置き換える）。このスクリプトは512pxのPNGとWebPを生成し、`image` に使う相対パスを出力する。`NOT_FOUND` なら別の絵文字を提案する。
4. 確定したfrontmatterのみを新規Markdownファイルに書く。`image` にはスクリプトの出力をそのまま使う。
5. `bash .agents/skills/new-article/scripts/validate-content.sh "<作成したファイル>"` を実行し、`OK` になるまでfrontmatterを修正する。本文には手を加えない。
6. ファイルパスと検証結果を報告し、本文はユーザーが書くこと、公開時に `is_draft` を `false` にすることを伝える。ユーザーが本文を書いた後、AI slopのように見えるか点検を依頼した場合は [check-ai-slop](../check-ai-slop/SKILL.md) を使う。

画像取得には `gh`、`curl`、`rsvg-convert`、`cwebp` を使う（macOS想定）。取得できない画像パスを捏造しない。
