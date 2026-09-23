---
name: new-review
description: このブログの書評を新規作成するときに使う。書誌情報と書影を調べ、本文が空のMarkdownを作る。
---

# 書評の新規作成

`backend/content/reviews/` に書評を1件作る。**本文は書かない。** frontmatterと空の本文だけを用意し、執筆はユーザーに任せる。

## 規約

- ファイル名は当日の `YYYY-MM-DD.md`。同日分があれば未使用の `YYYY-MM-DD-N.md` (`N` は2から)を使う。既存ファイルを上書きしない。
- `id` は英小文字・数字をハイフンでつないだスラッグで、`backend/content/articles/` と `backend/content/reviews/` の全体で一意にする。両ディレクトリの `id:` を確認する。
- frontmatterは次のフィールドを使う。`description` は短い紹介文、`rating` は1〜5の整数。本文や `## 概要` / `## 感想` の定型見出しは追加しない。

  ```yaml
  ---
  id: kebab-case-english-words
  title: "書籍タイトル"
  description: "あらすじ・紹介文"
  jp_e_code: "電子版コード。無ければ空文字"
  image: images/<id>.jpg
  rating: 1
  is_draft: true
  ---
  ```

- `jp_e_code` は [books.or.jp](https://www.books.or.jp/) で確認できる20桁の電子版コード。見つからなければ空文字とし、捏造しない。
- `is_draft: true` は、空の本文を公開しないために必須。ユーザーが本文を完成させた後に `false` に変更できる。
- `web/src/lib/utils.ts` の `getWebpPath` は `image` の拡張子を `.webp` に変換する。JPGと同名のWebPの両方が必要。
- `web/static/content/templates/review.md` は現行のfrontmatterと食い違うため参照しない。

## 手順

1. 書籍タイトル、出版社、感想の方向性、評価（1〜5）と紹介文に入れたい内容を確認する。すでに分かっていることは聞き直さない。評価や紹介文を推測で埋めない。
2. `bash .agents/skills/new-review/scripts/lookup-jp-e-code.sh "<書籍タイトル>" "<出版社>"` を実行する。`NOT_FOUND`（終了コード1）なら `jp_e_code` は空文字にする。既存の `id:` を確認する。
3. 書籍タイトル、重複しないid、得られた `jp_e_code` をまとめて提示してユーザーに確認する。本文は提案しない。
4. 当日の未使用ファイル名を決め、`bash .agents/skills/new-review/scripts/fetch-review-cover.sh "<書籍タイトル>" "<出版社>" "<id>"` を実行する。成功時はJPGとWebPが作られ、`image` に使う相対パスが出力される。
5. frontmatterのみを新規Markdownファイルに書く。書影が `NOT_FOUND` なら `image: images/<id>.jpg` とし、ユーザーに `web/static/content/reviews/images/<id>.jpg` と同名の `.webp` の配置を依頼する。
6. `bash .agents/skills/new-review/scripts/validate-content.sh "<作成したファイル>"` を実行する。書影以外の指摘はfrontmatterを修正して再検証する。書影が未配置なら検証が `NG` になることを正直に報告する。本文には手を加えない。
7. ファイルパスと検証結果を報告し、本文はユーザーが書くこと、公開時に `is_draft` を `false` にすることを伝える。ユーザーが本文を書いた後、AI slopのように見えるか点検を依頼した場合は [check-ai-slop](../check-ai-slop/SKILL.md) を使う。

書誌検索・画像取得には `curl` と `cwebp` を使う（macOS想定）。取得できない書誌情報や画像を捏造しない。
