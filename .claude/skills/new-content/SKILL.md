---
name: new-content
description: 記事(article)・書評(review)・プレゼンスライド(slide)の新規ファイルを作成する。article/reviewはbackend/content配下にMarkdownを作成し、タイトル・id・画像(articleはfluentui-emoji、reviewはbooks.or.jpの書影とjp_e_code)を対話的に確定してfrontmatterのみを書き込む(本文はユーザーが書くため生成しない)。slideはweb/static/content/slides配下にHTMLを作成し、スライド本文(見出し・構成・文章)はClaudeが下書きする(article/reviewとは逆の責務)。「記事を書きたい」「書評を追加したい」「スライドを作りたい」などで使う。
---

# 記事/書評/スライド作成

`backend/content/articles/`・`backend/content/reviews/` または `web/static/content/slides/` に新規ファイルを1つ作成する。

**article/review: 本文(記事・書評の中身)は絶対に書かない。frontmatterと空の本文だけを書き込み、続きはユーザーに書かせる。**

**slide: これとは責務が逆。スライド本文(見出し・構成・文章)はClaudeが下書きする。** 「スライド作成時の文章ルール」(後述)が存在するのは、slideに限り書く主体がClaude自身だから。各手順の `(〜のみ)` 表記でどちらの規則が効いているか必ず確認すること。上記の「本文は書かない」原則はslideには適用されない。

引数(`$ARGUMENTS`)に `article`・`review`・`slide` があればその種別を使う。なければユーザーに質問する。

前提ツール(macOS想定): `gh`、`curl`、`rsvg-convert`(`brew install librsvg`)、`cwebp`(`brew install webp`)。slideの本文・図表はClaudeが直接HTMLとして記述する(図表はテンプレート付属のHTML/CSSコンポーネントを使うため追加ツールは不要)。ただしHTML→PDF変換ツールはこのリポジトリに存在しない(既知のギャップ、後述)。

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
- `web/static/content/templates/slide.html` は上記の古いarticle/reviewテンプレートとは異なり、**現行有効なデザイン仕様**(HTML冒頭のコメントに設計方針が明記されている)。slide作成時は必ずこれをコピー元にし、デザイン(CSS・レイアウト・budoux-ja読み込み)そのものは変更しない。
- slideは `backend/content` ではなく `web/static/content/slides/<id>.html` に配置する。article/reviewの `YYYY-MM-DD.md` 命名は適用されない。`<id>` は既存の `how-to-speed-up-local-llm-inference-on-pc.pdf` のような、内容を表すkebab-caseスラッグ(日付を含めない)。slideにはfrontmatterという概念が無く、`web/src/routes/slides/+page.server.ts`/`[name]/+page.server.ts` が `web/static/content/slides/*.pdf` を列挙し、拡張子を除いたファイル名をそのままid/URLスラッグとして扱う。このidはslidesディレクトリ内でのみ一意であればよく、articles/reviewsのidと重複してもよい(別URLルートプレフィックスのため)。
- スライドのレイアウト検証には `web/scripts/validate-slide-layout.mjs` を使う(手順11参照)。`puppeteer-core`でローカルのGoogle Chrome/Chromiumを操作して実際にHTMLを描画し、各`.slide`が1920x1080pxちょうどに収まっているか(`scrollHeight`/`scrollWidth`が`clientHeight`/`clientWidth`を超えていないか)を実測でチェックする。budoux-jaによる実際の改行結果も反映されるため、フォントサイズ変更や画像サイズ変更のたびにこれを実行し、はみ出しが無いことを確認すること。図表はHTML/CSSで直接組むため(後述「図表の使用方針」参照)、非同期描画や外部リクエストの完了待ちは不要。
- 既知のギャップ(このスキル改修では対応しない):
  1. HTML→PDF変換の自動化ツール(Puppeteer/Playwright等)はこのリポジトリに無い。既存の`web/static/content/slides/*.pdf`にHTMLソースは存在せず、"sources"的なサブフォルダ規約も無い。このスキルはHTML作成までを行い、PDF化はユーザーが別途手動で行う。
  2. `scripts/validate-content.sh` はarticles/reviewsのMarkdown frontmatter専用で、slideのHTML構造(必須要素の有無など)を検証するロジックは無い。レイアウト(はみ出し)の検証は上記の`validate-slide-layout.mjs`でカバーするが、それ以外の構造チェックは今回追加しない。
  3. article/reviewにある画像取得スクリプト(手順8/9)に相当するslide用の仕組みは無い。実写真(スクリーンショット等)が必要な場合の自動取得手段は無く、無ければユーザーに画像ファイルの提供を依頼する(実在しない画像パスをそれらしく埋めない)。`slide.html`の`<img>`はサンプルSVG(data:image/svg+xml、外部ファイル無し)がプレースホルダーとして既に入っているため、実データが無い場合はこのサンプルのまま残してもよい(壊れた画像アイコンにはならない)。**関係性・フロー・サイクル・簡単なチャートを表す図は後述「図表の使用方針」のHTML/CSSコンポーネントを使う。単一の概念を表す装飾アイコン(人物・鍵・虫眼鏡など、矢印でつながる関係を持たないもの)は図自体を入れない選択肢を優先する。**
- 表(`table.simple-table`)のキャプションは`<caption>`で上(`caption-side: top`)、図(`figure`)のキャプションは`<figcaption>`で下(`<img>`直後に書くだけでよい)に配置する。この位置関係は固定で、逆にしない。
- `figure img`は`object-fit: contain`(`cover`にしない)。配置先(1カラムの横長figureと2カラムの縦長figureなど)でコンテナの縦横比がまちまちなため、`cover`だと画像側の縦横比次第で内容の一部が見切れることがある。`contain`なら常に画像全体が収まる。既存の`<img>`のCSSは変更せずそのまま使い、個別に`style="object-fit: cover"`等で上書きしない。自作のプレースホルダーSVG(data:image/svg+xml)を書くときは、画像全体を覆う背景矩形(`<rect>`でキャンバス全面を塗るなど)を入れない。`contain`で余白ができても透過なのでスライド本体の背景と自然に馴染む(枠線・背景色を持つ「カード」を避ける方針とも一致する)。

## 手順

### 1. 種別の確認
引数で `article`/`review`/`slide` が指定されていなければ AskUserQuestion で聞く(選択肢: 記事/書評/スライド)。

### 2. 内容のヒアリング
- article: 「何について書きたいか」を自由記述で質問する。
- review: 書籍の**タイトル**と**出版社**(jp_e_code検索に必須)、および感想の方向性(任意)を質問する。
- slide: **テーマ・発表タイトル案**、**想定聴衆**(社内LT/社外、技術レベル)、**伝えたい要点**(自由記述、複数可)、スライドに使いたい**具体的な数値・コード例・比較対象**、実在する**参考文献**(あれば)を質問する。articleと異なりClaudeが本文を書くため、事実関係(数値・固有名詞・コードの前提)はここで確認し、不明な点を後で推測で埋めない。

回答を待つ。

### 3. id重複チェック用の既存id一覧を取得
(article/review のみ)
```bash
grep -h '^id:' backend/content/articles/*.md backend/content/reviews/*.md
```
ここで得たid一覧と、これから提案するidが重複しないことを確認する。

(slide のみ)
```bash
ls web/static/content/slides/ | sed -E 's/\.[^.]+$//' | sort -u
```
`web/static/content/slides/` 直下の既存ファイル(現状は`.pdf`のみだが今後`.html`も混在しうるため拡張子を問わず列挙する)からidを列挙し、重複しないことを確認する。articles/reviewsのidとは別名前空間のため、そちらとの重複は問わない。

### 4. (review のみ) jp_e_codeの検索
手順2で得た書籍タイトル・出版社で検索する。

```bash
bash .claude/skills/new-content/scripts/lookup-jp-e-code.sh "<書籍タイトル>" "<出版社>"
```

標準出力に20桁のコードが出れば、それをjp_e_codeの提案値にする。`NOT_FOUND`(終了コード1)の場合は電子版が存在しないとみなし、jp_e_codeは空文字で提案する(検索を再試行したり、値を捏造したりしない)。

### 5. タイトル・id・(article: 絵文字 / review: jp_e_code)を提案する
- タイトル: ヒアリング内容から(reviewなら書籍タイトルそのまま、articleなら)日本語タイトル案を1つ提案する。slideの場合は表紙に使う発表タイトル案を1つ提案する。
- id: 英単語をハイフンで繋いだkebab-case案を1つ提案する。手順3の一覧と重複しないこと。
- article の場合のみ、内容にふさわしい絵文字を [microsoft/fluentui-emoji](https://github.com/microsoft/fluentui-emoji) から1つ提案する(英語の絵文字名、例: "White Flag")。
- review の場合のみ、手順4で得たjp_e_code(または空文字)を提示する。
- slide の場合、絵文字・jp_e_codeの提案は不要(slideにfrontmatterという概念は無い)。

AskUserQuestion でこれらをまとめて提示し、承認を求める。修正依頼があれば直るまで再提案する。
- (article/review) この提案はあくまでタイトル・id・画像・jp_e_codeであり、本文には一切触れない。
- (slide) この提案はタイトル・idのみであり、スライドの構成・本文は次の手順6で別途下書きする。

### 6. (slide のみ) スライドの構成決定・本文ドラフト作成
article/reviewと異なり、**slideはスライド本文(見出し・構成・文章)をClaudeが下書きする**(冒頭の「本文は書かない」原則はslideには適用されない)。

1. `web/static/content/templates/slide.html` を読み、コピー元にする(デザイン・CSS・budoux-ja読み込みは変更しない)。
2. 手順2のヒアリング内容から、スライド全体の構成を組み立てる。使えるスライド種別:
   - 表紙(`.slide--cover`、1枚目固定)
   - 1カラム(`.cols-1`): 見出し+リード文+画像+箇条書き
   - 2カラム(`.cols-2`): 画像とテキストの左右配置。画像から見せて興味を引きたい流れなら画像を左、結論から入りたい流れなら画像を右にする
   - 比較(`.compare`): 本当に二択の対比があるときのみ使う。3つ以上を無理に押し込まない
   - コード例(`.code-block`): 説明すべき具体的なコードがあるときのみ使う
   - 表(`table.simple-table`): 比較すべき具体的な数値があるときのみ使う
   - フロー図(`.flow-diagram`/`.flow-chain`/`.flow-compare`/`.flow-converge`/`.bar-chart`): 関係性・手順・サイクル・簡単なチャートなど、ノードと矢印(または軸と値)で表せる内容にのみ使う。詳細は後述「図表の使用方針」を参照
   - カード代替(`.minimal-list`/`.statement-compare`/`.chat-thread`/`.sticky-board`): `.card-grid`/`.compare`が単調・抽象的になりすぎる場合、または実データのスクリーンショットを避けたい場合に使う。詳細は後述「図表の使用方針」内の「カード代替コンポーネント」「実データのスクリーンショットより抽象化した再現を優先する」を参照
   - 高橋メソッド風(`.slide--impact`): 聴衆の注意を引く一言・話の転換点に**1〜2回まで**
   - 参考文献(`ol.references`、最終ページ固定)

   1枚目=表紙、最終ページ=参考文献は必ず守り、間の種別は内容に応じて選ぶ(全種別を使う必要はない)。
3. 後述の「スライド作成時の文章ルール」に従って各スライドの見出し・本文・コード・表・参考文献を作成する。
4. 完成したHTML全体(`<html>`〜`</html>`)を組み立てる。`<title>`と表紙の`<h1>`をタイトルにする。ページ番号は既存の`<script>`が自動計算するため手を加えない。
5. 実写真(スクリーンショット等)の`<img>`を使いたい場合、画像ファイルを自動取得する手段はこのスキルには無い(手順8/9の画像取得スクリプトはarticle/review専用)。実データが無ければユーザーに画像ファイルの提供を依頼する。実在しない画像パスをそれらしく埋めない。**単一概念を表す装飾アイコン(人物・鍵など)が欲しい場合、幾何学図形を自作しない。** HTML/CSSのフロー図コンポーネントで表現できる関係性が無いなら、図自体を入れない選択肢を優先する(分量・改行の目安の項を参照)。**Slack・Notion等の実スクリーンショットを使いたい場合、先に「図表の使用方針」内の「実データのスクリーンショットより抽象化した再現を優先する」を検討する。**実データを使う場合は、第三者の実名等が写り込んでいないか確認し、写っていれば加工(モザイク等)をユーザーに確認した上で行う。

### 7. ファイル名の決定
(article/review・同日複数作成は連番)
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

(slide のみ)
```bash
DIR="web/static/content/slides"
FILENAME="$ID.html"
if [ -e "$DIR/$FILENAME" ]; then
  echo "既に存在します。手順5のidを変更してください: $DIR/$FILENAME" >&2
fi
```
slideは日付ではなく手順5で確定したid(kebab-caseスラッグ)をそのままファイル名にする(既存`*.pdf`と同じ命名規則)。既存ファイルへの上書きは行わない。

### 8. (article のみ) 絵文字画像の取得・変換
承認された絵文字名(例 "White Flag")について実行する。

```bash
bash .claude/skills/new-content/scripts/fetch-emoji-image.sh "White Flag"
```

成功すると `web/static/content/articles/images/` に512x512 PNG(fluentui-emojiの2D(Flat)版SVGをラスタライズしたもの)と同名`.webp`を作成し、標準出力に `images/Microsoft-Fluentui-Emoji-Flat-White-Flag.512.png` のような相対パスを返す。これをそのままfrontmatterの `image:` に書く。fluentui-emoji側のフォルダ名は "White flag" のように先頭のみ大文字のsentence caseだが、スクリプト内で大文字小文字を無視して検索するため入力の表記は問わない。

`NOT_FOUND`(終了コード1)の場合はその絵文字がfluentui-emojiに存在しない。手順5に戻って別の絵文字を再提案する。

### 9. (review のみ) 書影画像の取得・変換
手順2で得た書籍タイトル・出版社と、確定したidで実行する。

```bash
bash .claude/skills/new-content/scripts/fetch-review-cover.sh "<書籍タイトル>" "<出版社>" "<id>"
```

成功すると `web/static/content/reviews/images/<id>.jpg` と同名`.webp`を作成し、標準出力に `images/<id>.jpg` を返す。これをそのままfrontmatterの `image:` に書く。

`NOT_FOUND`(終了コード1)の場合は書影を自動取得できなかったということなので、`image:` には `images/<id>.jpg` を書いた上で、「書影画像を `web/static/content/reviews/images/<id>.jpg` および同名`.webp`として配置してください」とユーザーに伝える。

### 10. ファイル作成
手順7で決めたパス(`$DIR/$FILENAME`、slideは `web/static/content/slides/$FILENAME`)を新規作成する。
- article/review: 確定したfrontmatterのみを書き込む。本文は空にする。review の定型見出し(`## 概要` / `## 感想`)もファイルには書かず、口頭で目安として伝えるだけにする。
- slide: 手順6で作成したHTML全体を書き込む。

### 11. 検証
(article/review)
作成したファイルに対して検証スクリプトを実行し、ファイル名とfrontmatterの構造が規約通りであることを確認する。

```bash
bash .claude/skills/new-content/scripts/validate-content.sh "$DIR/$FILENAME"
```

`NG` が出た場合はファイルを削除せず、指摘された内容(必須フィールドの欠落、id重複、rating範囲外、image実体の不在など)に沿ってfrontmatterを修正し、`OK` になるまで再実行する。**修正はfrontmatterのみに留め、本文には手を加えない。**

(slide のみ) レイアウト(1920x1080pxからのはみ出し、および見出し(h2)と本文の間隔の統一性)を実ブラウザで検証する。frontmatterの構造チェックに相当する自動検証(`validate-content.sh`)は無いため、それは対象外。

```bash
node web/scripts/validate-slide-layout.mjs "web/static/content/slides/$FILENAME"
```

`NG`(いずれかのスライドではみ出しあり、1920x1080ちょうどでない、または見出しと本文の間隔がスライドごとに異なる)が出た場合はファイルを削除せず、該当スライドの文章量・画像サイズを調整し(「分量・改行の目安」参照)、`OK` になるまで再実行する。間隔の不統一が出た場合は、h2やその直後の要素に個別の`margin`インラインスタイルを足していないか確認する(`.slide-inner`/`.cols-1`の`gap`と二重に加算されるのが典型的な原因)。ローカルにGoogle Chrome/Chromiumが無い等の理由でスクリプト自体が実行できない場合は、その旨を完了報告で明示した上でこの手順を省略してよい(中途半端な簡易チェックをその場で追加しない)。

### 12. 完了報告
(article/review) 作成したファイルパスと検証結果(`OK`)を伝え、本文はユーザー自身が書くことを伝えて終了する。

(slide のみ) 作成したファイルパス(`web/static/content/slides/<id>.html`)を伝えたうえで、必ず次を明示する:
1. 本文(見出し・文章・コード例・参考文献)はClaudeが下書きしたものであり、article/reviewと異なりユーザーの手直し前提であること。特に数値・固有名詞・参考文献の実在性の確認を依頼する。
2. HTML→PDF変換はこのスキルでは行わない(変換ツール未導入のため)。`/slides`に掲載するには、別途手動でPDF化し同じファイル名(`<id>.pdf`)で`web/static/content/slides/`に配置する必要がある。
3. 手順11のレイアウト検証結果(`OK`、または未実施ならその理由)。frontmatterの構造チェックに相当する自動検証は無いこと。

## スライド作成時の文章ルール

slideの本文(見出し・リード文・箇条書き・コード・表・参考文献)を下書きする際は、以下に従う。article/reviewでは本文を書かないが、slideに限り書く方針であることに注意する(冒頭の原則の例外)。

### 文体
- 文末に「。」(句点)を使わない
- 「です」「ます」などの敬体は使わない。体言止め・言い切りを基本にする
- 1文=1主張。1つの見出し・1つの箇条書き項目に複数の主張を詰め込まない
- 見出しと本文(リード文・箇条書き)を同じ内容の言い換えにしない。見出しは結論・トピック、本文は見出しにない情報(根拠・数値・詳細)を足す
  - 悪い例: 見出し「レスポンス速度が改善」→本文「レスポンス速度が良くなった」
  - 良い例: 見出し「レスポンス速度が改善」→本文「p95レイテンシ 480ms→120ms」

### 分量・改行の目安
budoux-jaが文節単位で自動改行するため、改行位置そのものを指定する必要は無い。以下は1カラム(実効幅1680px = 1920px − 左右padding 120px×2)を基準にした、詰め込みすぎを防ぐための目安(全角文字数、上限であり厳密な文字数管理ではない)。

| 要素 | フォントサイズ | 1行あたりの目安(全角) |
|---|---|---|
| h1(表紙タイトル・高橋メソッド風) | 72px | 約23文字。表紙は1行に収める。`.slide--impact`は上限まで書かず一言程度の短さを優先する |
| h2(セクション見出し) | 56px | 約30文字。1行に収める |
| .lead(リード文) | 36px | 約46文字/行、最大2行 |
| bullets(箇条書き1項目) | 34px | 約49文字。1項目1行が基本 |

`.compare`/`.stat-list`のように2〜3カラムに文章を置く場合、実効幅はカラム数に応じて狭くなるため、上表よりさらに短くする。コード例・表・参考文献は上表の対象外(コードは実際のコードの区切りに従い、表・参考文献は縦方向の高さ(スライドは`overflow:hidden`で見切れる)に収まるよう数値・固有名詞中心の短い記述にする)。情報が多い場合は文章を削るか複数スライドに分割する。

### AI生成コンテンツ特有の表現を避ける
- 「革新的」「画期的」「鍵となる」「〜が重要です」のような具体性の無い誇張・紋切り型の形容語を使わない。代わりに具体的な事実・数値・固有名詞で語る
- 箇条書きを機械的に3つに揃えない(「〜な3つのポイント」という定型を避ける)。実際に伝えたい項目数にする
- 絵文字は使わない
- 抽象的な一般論より検証可能な具体的事実を優先する(例: 「開発体験が大きく向上」ではなく「ビルド時間 12分→2分」)

### 技術系スライドとしての配慮
- コード例は`foo`/`bar`/`hoge`のような汎用プレースホルダーではなく、内容の文脈に即した具体的な変数名・関数名にする
- 想定聴衆は技術者のため、専門用語は説明なしで使ってよい
- 表で比較する際は「◎/○/△」のような定性的な記号評価ではなく、具体的な数値(実測値・件数・時間など)を入れる。数値が無い項目は無理に埋めず、その旨を明示するか項目自体を削る

### スライド種別ごとの使用方針
- 2カラム(画像+テキスト)は、画像から見せて興味を引きたい流れなら画像を左、結論・主張から入りたい流れなら画像を右にする。片方に決め打ちしない
- `.compare`は本当に二択の対比があるときだけ使う。3つ以上を無理に2つに削らない(3つ以上並べたいときは`.stat-list`を検討する)
- コード例・表は、実際に見せるべき具体的なコード・数値がある場合のみ使う。中身が薄いのに体裁だけ整えるために挿入しない
- `.slide--impact`(高橋メソッド風)は1デッキ中1〜2回までにする。聴衆の注意を引きたい一言・話の転換点にのみ使う

### 参考文献のルール
- 実在し検証可能な文献だけを載せる。著者名・タイトル・出版社・URLを捏造しない
- ヒアリングで得られた情報の範囲でしか特定できない場合、それらしい文献情報をでっち上げず、「(要確認)」など未確認である旨を明記するか、ユーザーに実際の情報源を確認する

## 図表の使用方針(HTML/CSSが第一選択)

「それっぽいが意味の伝わらない自作SVG図形」を避けつつ、外部レンダリングサービスにも依存しないため、概念図・フロー図・簡単なチャートは**HTML/CSS(flexbox)で直接組む**のが第一選択。`web/static/content/templates/slide.html`にすでに用意されている`.flow-diagram`/`.flow-chain`/`.flow-compare`/`.flow-converge`/`.bar-chart`をそのまま使う(コピー元テンプレートに実例スライドあり)。

**経緯**: 最初は自作の幾何学図形SVG(意味が伝わらないと指摘された)→ 外部レンダリングサービスによるSVG画像生成(サーバーサイド描画、文字が本文より極端に小さく配色もテーマと合わずスライドとして見劣りする)と試行錯誤し、最終的に**本文と同じフォント・配色・文字サイズを完全に共有できるHTML/CSSでの直接構築**に落ち着いた。外部レンダリングサービスやクライアントサイド描画ライブラリを新規に第一選択として使わない。

### 使い方
1. どのコンポーネントが合うか選ぶ(下記「いつ使う/使わないか」)
2. テンプレートの該当コンポーネントのHTMLをコピーし、ノードのテキストと色クラス(`.muted`/`.accent`/`.warn`)を内容に合わせて書き換える
3. 2カラムの片側に置く場合は`<div class="col diagram">`、1カラムで見出し+本文と同居させ高さを抑えたい場合は`.flow-diagram.inline`(固定220px)、スライド全体に大きく表示する場合は`.flow-diagram`(`flex:1`で自動拡張)を使う

### コンポーネント一覧
- `.flow-chain`: 「A → B → C」のような直線的な手順・サイクル。ノードは`<span class="node muted|accent|warn">`、矢印は`<span class="arrow">`
- `.flow-relation`(`.flow-chain`を`.flow-relation`で囲む): 2要素の関係を1行で示し、`.flow-caption`で短い注釈(「当てはまらない」等)を添える
- `.flow-compare`: 「誤解していた向き」vs「気づいた向き」のような2〜3行の対比。`.flow-compare-row.wrong`/`.right`で行ごとに色分けし、罫線区切りのみでカード化しない
- `.flow-converge`: 複数要素が1つの結論に収束する図(例: 3つの壁→型への依存)。`.flow-converge-sources`に縦積みのリスト、`.flow-converge-target`に収束先
- `.bar-chart`: 棒グラフ。バーの高さは`style="height:Npx;"`で個別に指定し、値・週などのラベルは本文と同じ文字サイズ(30px前後)で表示する(小さすぎて読めない状態を避ける)

**矢印(`.arrow`/`.step-arrow`/`.statement-arrow`)には意味的な色を付けない。** 常にmutedグレーで統一し、色分けが必要な場合はnode/テキスト側の`.accent`/`.warn`で表現する。矢印まで色分けすると煩雑になり、かえって読み取りにくくなった経緯がある。

### カード代替コンポーネント(`.card-grid`/`.compare`が単調・抽象的すぎる場合)
「矢印とカードばかりで抽象的」という指摘を受けて追加した、より具体的・直感的な代替表現。**これらも「カードを避ける」という基本方針の範囲内の代替案であり、`.card-grid`と並ぶ新たな例外を無制限に増やすものではない。**
- `.minimal-list`/`.minimal-item`(+`.minimal-num`/`.minimal-content`/`.minimal-title`/`.minimal-desc`): 罫線区切りの番号リスト。3〜4項目の一覧を`.card-grid`より軽い見た目で並べたいときに使う
- `.statement-compare`/`.statement-box`(+`.statement-tag`/`.statement-text`/`.statement-arrow`): ビッグ・ステートメント対比。2つの短い主張を罫線+太字で強く対比させたいときに使う。`.compare`より「主張」そのものを主役にしたい場合に選ぶ(`.statement-label`は下記「カードにはbodyを書かない」により使わない)
- `.chat-thread`/`.chat-item`(+`.chat-label`/`.chat-plain`/`.chat-bubble.said`): 「心の中で思っていたこと」と「実際の発言・行動」のギャップを見せたいときに使う。`.chat-item.thought`側は`.chat-plain`(地の文、箱なし)、`.chat-item.said`側は`.chat-bubble.said`(塗りつぶしの吹き出し)で役割を分ける(両方を吹き出し化すると尖った形状同士がぶつかり合い読みにくくなった経緯がある)
- `.sticky-board`/`.sticky-note`(+`.sticky-note-tag`、`.warn`/`.accent`で色を変える): 付箋UI。レトロスペクティブ(KPT等)やブレインストーミングの内容を再現したいときに使う

### カードにはbodyを書かない
`.card`/`.comp-card`/`.statement-box`のような、罫線・背景・左アクセントボーダーで囲われた「カード」型の箱には、タイトル・タグなど短い1行のみを書く。**タイトルの言い換えに過ぎない説明文(body/label)を追加しない。** 例: タイトル「他者に委ねる」に対しさらに「自分の現在地と弱さを開示し、協力を仰ぐ」のような一文を重ねて書かない。これはテキストとの重複であり、本当に補足したい情報があるならカードの外(リード文・箇条書き)で表現する。過去に`.comp-card`の`.card-body`や`.statement-box`の`.statement-label`がタイトルの言い換えになっていて削除した経緯がある。

この制約は「カードを避ける」方針の一部であり、カード自体を持たない`.minimal-item`(罫線区切りの番号リスト)や`.timeline-item`(タイムライン)の説明文(`.minimal-desc`等)には適用されない。これらは背景・枠を持たない列挙で、説明文自体がその項目の主内容であるため。

### いつ使う/使わないか
- 使う: 関係性(A→B)、手順の連なり、堂々巡り・循環、分岐、複数要素が1つの原因に収束する構造、簡単な棒グラフなど、**ノードと矢印(または軸と値)で表現できる内容**
- 使わない: 人物・鍵・虫眼鏡・チェックリストのような、関係性を持たない単一概念の装飾アイコン。この場合は図自体を入れない(分量・改行の目安の項、および「抽象的になりすぎる場合は図表を入れる必要はない」という既存方針に従う)
- 箇条書きで既に列挙している内容をそのままフローチャート化するような、テキストとの単純な重複は避ける
- **`.flow-caption`は、ノード・矢印だけでは伝わらない追加情報がある場合にのみ使う。** ノードの並びをそのまま文章化しただけの一文(例: 「チーム→個人」の図に「チームから個人へ仕事が降ってくる感覚」という説明文を添える)は書かない。図が既に伝えている内容を言い換えただけの注釈はテキストとの重複であり、`.flow-compare-tag`(「これまでの認識」等)と箇条書きだけで十分伝わる。`.flow-caption`を使ってよいのは、堂々巡りの帰結(例: 「適応(アクション)が決まらない堂々巡り」)のように、図だけでは読み取れない新しい解釈・結論を足す場合に限る
- 矢印・カードを多用したスライドが続くと「抽象的」という印象を与えやすい。1デッキの中で`.flow-chain`系と上記「カード代替コンポーネント」を適度に混在させ、同じ構図・同じ見た目のコンポーネントばかりが連続しないようにする

### 実データのスクリーンショットより抽象化した再現を優先する
Slack・Notion・チャットツール等の実スクリーンショットを`<img>`でそのまま貼ると、以下の問題が起きやすい(社内レビューで実際に指摘された経緯がある)。
- 画面全体の詳細情報が写り込み、スライドで伝えたい要点に対してノイズになる
- 関係の無い社内の固有名詞・プロジェクト名・実データが写り込み、社外公開前レビューで指摘され差し替えになる
- 第三者(同僚等)の実名が写り込み、加工(モザイク等)が必要になる

実データでなければ成立しない場合を除き、まず上記の「カード代替コンポーネント」(`.chat-thread`で発言と本音のギャップを、`.sticky-board`でふりかえりボードの要点を再現する等)で置き換えられないか検討する。実スクリーンショットを使う場合も、第三者の実名は必ず加工し、画像に写っている情報が実データか研修等の模擬データかが伝わるようにする。

### 配色(3トーンで統一する)
- `accent`(ティール `#2dd4bf`): 経験主義・適応・成功・カイゼンの方向性
- `warn`(アンバー `#fbbf24`): 壁・警告・違和感・過渡期
- `muted`(グレー `#aab4bd`): 過去の型・停滞・前提

### レイアウトのリズム
`.flow-chain`/`.flow-relation`を2カラム(`cols-2`)の片側に固定で置き続けると、前回のレビューで指摘された通りマンネリ化する。図と本文の上下配置を入れ替える(diagramを上にするスライドと下にするスライドを混在させる)、`.flow-converge`/`.flow-compare`を使う全幅の1カラムスライドを間に挟むなど、同じ構図が3枚以上連続しないように配置する。

HTML/CSSで表現しきれないほど複雑な関係性が必要な場合は、外部レンダリングサービスに頼らず、図を簡略化するか内容を複数スライドに分割する。

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
