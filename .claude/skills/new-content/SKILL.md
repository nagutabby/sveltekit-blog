---
name: new-content
description: 記事(article)または書評(review)の新規Markdownファイルをbackend/content配下に作成する。タイトル・id・(articleなら)絵文字画像を対話的に確定してfrontmatterのみを書き込む。本文はユーザーが書くため生成しない。「記事を書きたい」「書評を追加したい」「新しい記事のテンプレートを作って」などで使う。
---

# 記事/書評テンプレート作成

`backend/content/articles/` または `backend/content/reviews/` に新規Markdownファイルを1つ作成する。**本文(記事の中身)は絶対に書かない。frontmatterと空の本文だけを書き込み、続きはユーザーに書かせる。**

引数(`$ARGUMENTS`)に `article` または `review` があればその種別を使う。なければユーザーに質問する。

## 前提知識(このリポジトリの規約)

- ファイル名は `YYYY-MM-DD.md`(当日の日付、ISO 8601)。同日に複数作る場合の連番(`-2`など)は使わない。当日分がすでに存在する場合は**上書きせずエラーにして中断する**。
- article の frontmatter は3フィールドのみ:
  ```yaml
  ---
  id: kebab-case-english-words
  title: 日本語タイトル
  image: images/Microsoft-Fluentui-Emoji-3d-<Name>-3d.1024.png
  ---
  ```
- review の frontmatter:
  ```yaml
  ---
  id: kebab-case-english-words
  title: "書籍タイトル"
  description: "あらすじ・紹介文"
  jp_e_code: "書誌コード(未取得なら空文字でよい)"
  image: images/<id>.jpg
  rating: 1〜5の整数
  ---
  ```
- `id` は英単語をハイフンで繋いだkebab-caseで、`articles/`・`reviews/`全体で一意でなければならない(URLスラッグとして使われる)。
- `web/static/content/templates/{article,review}.md` に古いテンプレートが存在するが、`publishedAt`/`updatedAt`が残っていたり`id`が無かったりして**現行実装と食い違っている**。参照せず、上記の実測frontmatterに従うこと。

## 手順

### 1. 種別の確認
引数で `article`/`review` が指定されていなければ AskUserQuestion で聞く。

### 2. 内容のヒアリング
「何について書きたいか」を自由記述で質問する。回答を待つ。

### 3. id重複チェック用の既存id一覧を取得
```bash
grep -h '^id:' backend/content/articles/*.md backend/content/reviews/*.md
```
ここで得たid一覧と、これから提案するidが重複しないことを確認する。

### 4. タイトル・id(・article なら絵文字)を提案する
- タイトル: ヒアリング内容から日本語タイトル案を1つ提案する。
- id: 英単語をハイフンで繋いだkebab-case案を1つ提案する。手順3の一覧と重複しないこと。
- article の場合のみ、内容にふさわしい絵文字を [microsoft/fluentui-emoji](https://github.com/microsoft/fluentui-emoji) から1つ提案する(英語の絵文字名、例: "White Flag")。

AskUserQuestion でこれらをまとめて提示し、承認を求める。修正依頼があれば直るまで再提案する。**この提案はあくまでタイトル・id・画像であり、本文には一切触れない。**

### 5. 当日ファイルの重複チェック
```bash
DATE=$(date +%Y-%m-%d)
test -f "backend/content/articles/$DATE.md" && echo exists
test -f "backend/content/reviews/$DATE.md" && echo exists
```
対象種別のファイルが既に存在する場合は、作成を中断してユーザーに報告する(上書きしない。別日付を使う・既存ファイルを手動でリネームするなどの対応はユーザーに委ねる)。

### 6. (article のみ) 絵文字画像の取得・変換
承認された絵文字名(例 "White Flag")について:

```bash
EMOJI_NAME="White Flag"                     # 承認された絵文字の英語名
SNAKE=$(echo "$EMOJI_NAME" | tr '[:upper:] ' '[:lower:]_')   # white_flag
HYPHEN=$(echo "$EMOJI_NAME" | tr ' ' '-')                     # White-Flag

# 実ファイル名を確認(スペースを含むディレクトリ名に注意)
gh api "repos/microsoft/fluentui-emoji/contents/assets/$EMOJI_NAME/3D" --jq '.[].name'

# ダウンロード
curl -sL "https://raw.githubusercontent.com/microsoft/fluentui-emoji/main/assets/$EMOJI_NAME/3D/${SNAKE}_3d.png" \
  -o /tmp/emoji_source.png

# 1024x1024 PNGにリサイズしてリポジトリの命名規則に合わせて配置
sips -s format png -z 1024 1024 /tmp/emoji_source.png \
  --out "web/static/content/articles/images/Microsoft-Fluentui-Emoji-3d-${HYPHEN}-3d.1024.png"

rm -f /tmp/emoji_source.png
```

`gh api` のディレクトリ一覧で得た実ファイル名と `${SNAKE}_3d.png` が一致しない場合は、実際のファイル名に合わせて `curl` のURLを修正する。

frontmatterの `image:` には `images/Microsoft-Fluentui-Emoji-3d-${HYPHEN}-3d.1024.png` を書く。

### 7. (review のみ) 画像の案内
review では画像の自動取得は行わない。`image:` には `images/<id>.jpg` を書き、「書影画像を `web/static/content/reviews/images/<id>.jpg` に配置してください」とユーザーに伝える。

### 8. ファイル作成
`backend/content/{articles,reviews}/YYYY-MM-DD.md` を新規作成し、確定したfrontmatterのみを書き込む。本文は空にする。review の定型見出し(`## 概要` / `## 感想`)もファイルには書かず、口頭で目安として伝えるだけにする。

### 9. 検証
作成したファイルに対して検証スクリプトを実行し、ファイル名とfrontmatterの構造が規約通りであることを確認する。

```bash
bash .claude/skills/new-content/scripts/validate-content.sh backend/content/{articles,reviews}/YYYY-MM-DD.md
```

`NG` が出た場合はファイルを削除せず、指摘された内容(必須フィールドの欠落、id重複、rating範囲外、image実体の不在など)に沿ってfrontmatterを修正し、`OK` になるまで再実行する。**修正はfrontmatterのみに留め、本文には手を加えない。**

### 10. 完了報告
作成したファイルパスと検証結果(`OK`)を伝え、本文はユーザー自身が書くことを伝えて終了する。
