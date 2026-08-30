---
name: new-content
description: 記事(article)または書評(review)の新規Markdownファイルをbackend/content配下に作成する。タイトル・id・画像(articleはfluentui-emoji、reviewはbooks.or.jpの書影とjp_e_code)を対話的に確定してfrontmatterのみを書き込む。本文はユーザーが書くため生成しない。「記事を書きたい」「書評を追加したい」「新しい記事のテンプレートを作って」などで使う。
---

# 記事/書評テンプレート作成

`backend/content/articles/` または `backend/content/reviews/` に新規Markdownファイルを1つ作成する。**本文(記事の中身)は絶対に書かない。frontmatterと空の本文だけを書き込み、続きはユーザーに書かせる。**

引数(`$ARGUMENTS`)に `article` または `review` があればその種別を使う。なければユーザーに質問する。

前提ツール(macOS想定): `gh`、`curl`、`rsvg-convert`(`brew install librsvg`)、`cwebp`(`brew install webp`)。

## 前提知識(このリポジトリの規約)

- ファイル名は `YYYY-MM-DD.md`(当日の日付、ISO 8601)。同日に複数作る場合は `YYYY-MM-DD-N.md`(Nは2から始まる連番)にする。既存ファイルへの**上書きは行わない**。
- article の frontmatter は3フィールドのみ:
  ```yaml
  ---
  id: kebab-case-english-words
  title: 日本語タイトル
  image: images/Microsoft-Fluentui-Emoji-Flat-<Name>.512.png
  ---
  ```
- review の frontmatter:
  ```yaml
  ---
  id: kebab-case-english-words
  title: "書籍タイトル"
  description: "あらすじ・紹介文"
  jp_e_code: "電子版コード(電子版が存在しない書籍は空文字でよい)"
  image: images/<id>.jpg
  rating: 1〜5の整数
  ---
  ```
  `jp_e_code` は [books.or.jp](https://www.books.or.jp/)(日本書籍出版協会 出版書誌データベース)で書籍タイトル・出版社から検索できる電子版(JP-e)コードで、`.claude/skills/new-content/scripts/lookup-jp-e-code.sh` で自動取得する。
- `id` は英単語をハイフンで繋いだkebab-caseで、`articles/`・`reviews/`全体で一意でなければならない(URLスラッグとして使われる)。
- `image` が指す画像は、`web/src/lib/utils.ts` の `getWebpPath` により拡張子を `.webp` に置き換えたパスのみが実際に `<img src>` として配信される(`Card.svelte`/`Header.svelte`)。**元画像(png/jpg)と同名の`.webp`が無いと画像が表示されない。** 画像取得スクリプトは両方を生成する。
- `web/static/content/templates/{article,review}.md` に古いテンプレートが存在するが、`publishedAt`/`updatedAt`が残っていたり`id`が無かったりして**現行実装と食い違っている**。参照せず、上記の実測frontmatterに従うこと。

## 手順

### 1. 種別の確認
引数で `article`/`review` が指定されていなければ AskUserQuestion で聞く。

### 2. 内容のヒアリング
- article: 「何について書きたいか」を自由記述で質問する。
- review: 書籍の**タイトル**と**出版社**(jp_e_code検索に必須)、および感想の方向性(任意)を質問する。

回答を待つ。

### 3. id重複チェック用の既存id一覧を取得
```bash
grep -h '^id:' backend/content/articles/*.md backend/content/reviews/*.md
```
ここで得たid一覧と、これから提案するidが重複しないことを確認する。

### 4. (review のみ) jp_e_codeの検索
手順2で得た書籍タイトル・出版社で検索する。

```bash
bash .claude/skills/new-content/scripts/lookup-jp-e-code.sh "<書籍タイトル>" "<出版社>"
```

標準出力に20桁のコードが出れば、それをjp_e_codeの提案値にする。`NOT_FOUND`(終了コード1)の場合は電子版が存在しないとみなし、jp_e_codeは空文字で提案する(検索を再試行したり、値を捏造したりしない)。

### 5. タイトル・id・(article: 絵文字 / review: jp_e_code)を提案する
- タイトル: ヒアリング内容から(reviewなら書籍タイトルそのまま、articleなら)日本語タイトル案を1つ提案する。
- id: 英単語をハイフンで繋いだkebab-case案を1つ提案する。手順3の一覧と重複しないこと。
- article の場合のみ、内容にふさわしい絵文字を [microsoft/fluentui-emoji](https://github.com/microsoft/fluentui-emoji) から1つ提案する(英語の絵文字名、例: "White Flag")。
- review の場合のみ、手順4で得たjp_e_code(または空文字)を提示する。

AskUserQuestion でこれらをまとめて提示し、承認を求める。修正依頼があれば直るまで再提案する。**この提案はあくまでタイトル・id・画像・jp_e_codeであり、本文には一切触れない。**

### 6. ファイル名の決定(同日複数作成は連番)
```bash
DATE=$(date +%Y-%m-%d)
DIR="backend/content/$TYPE"   # $TYPE は articles または reviews

if [ ! -e "$DIR/$DATE.md" ]; then
  FILENAME="$DATE.md"
else
  N=2
  while [ -e "$DIR/$DATE-$N.md" ]; do
    N=$((N + 1))
  done
  FILENAME="$DATE-$N.md"
fi
echo "$DIR/$FILENAME"
```
当日分がまだ無ければ `$DATE.md`、既にあれば空いている連番(`$DATE-2.md`, `$DATE-3.md`, ...)を使う。既存ファイルへの**上書きは絶対に行わない**。

### 7. (article のみ) 絵文字画像の取得・変換
承認された絵文字名(例 "White Flag")について実行する。

```bash
bash .claude/skills/new-content/scripts/fetch-emoji-image.sh "White Flag"
```

成功すると `web/static/content/articles/images/` に512x512 PNG(fluentui-emojiの2D(Flat)版SVGをラスタライズしたもの)と同名`.webp`を作成し、標準出力に `images/Microsoft-Fluentui-Emoji-Flat-White-Flag.512.png` のような相対パスを返す。これをそのままfrontmatterの `image:` に書く。fluentui-emoji側のフォルダ名は "White flag" のように先頭のみ大文字のsentence caseだが、スクリプト内で大文字小文字を無視して検索するため入力の表記は問わない。

`NOT_FOUND`(終了コード1)の場合はその絵文字がfluentui-emojiに存在しない。手順5に戻って別の絵文字を再提案する。

### 8. (review のみ) 書影画像の取得・変換
手順2で得た書籍タイトル・出版社と、確定したidで実行する。

```bash
bash .claude/skills/new-content/scripts/fetch-review-cover.sh "<書籍タイトル>" "<出版社>" "<id>"
```

成功すると `web/static/content/reviews/images/<id>.jpg` と同名`.webp`を作成し、標準出力に `images/<id>.jpg` を返す。これをそのままfrontmatterの `image:` に書く。

`NOT_FOUND`(終了コード1)の場合は書影を自動取得できなかったということなので、`image:` には `images/<id>.jpg` を書いた上で、「書影画像を `web/static/content/reviews/images/<id>.jpg` および同名`.webp`として配置してください」とユーザーに伝える。

### 9. ファイル作成
手順6で決めた `$DIR/$FILENAME` を新規作成し、確定したfrontmatterのみを書き込む。本文は空にする。review の定型見出し(`## 概要` / `## 感想`)もファイルには書かず、口頭で目安として伝えるだけにする。

### 10. 検証
作成したファイルに対して検証スクリプトを実行し、ファイル名とfrontmatterの構造が規約通りであることを確認する。

```bash
bash .claude/skills/new-content/scripts/validate-content.sh "$DIR/$FILENAME"
```

`NG` が出た場合はファイルを削除せず、指摘された内容(必須フィールドの欠落、id重複、rating範囲外、image実体の不在など)に沿ってfrontmatterを修正し、`OK` になるまで再実行する。**修正はfrontmatterのみに留め、本文には手を加えない。**

### 11. 完了報告
作成したファイルパスと検証結果(`OK`)を伝え、本文はユーザー自身が書くことを伝えて終了する。
