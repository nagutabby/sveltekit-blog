#!/usr/bin/env bash
# 書籍のタイトル・出版社から、books.or.jp(日本書籍出版協会 出版書誌データベース)の
# 書影サムネイルを取得し、review用の画像として配置する。
#
# Usage: fetch-review-cover.sh <title> <publisher> <id>
#
# 成功時: web/static/content/reviews/images/<id>.jpg と同<id>.webp を作成し、
#         標準出力に "images/<id>.jpg"(frontmatterのimage用の相対パス)を出力して終了コード0。
# 失敗時(該当書籍が見つからない等): 何も出力せず、標準エラーに理由を出して終了コード1。
#         中間生成物は残さない。呼び出し側はこの場合、画像配置をユーザーに案内する。
#
# 前提: books.or.jp の書影サムネイルはCloudFront経由のリファラーチェックがあるため、
# 同サイトのReferer付きでのみ取得できる。取得できるのは1サイズ(小サムネイル)のみ。

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"

TITLE="${1:-}"
PUBLISHER="${2:-}"
ID="${3:-}"

if [ -z "$TITLE" ] || [ -z "$PUBLISHER" ] || [ -z "$ID" ]; then
  echo "Usage: $0 <title> <publisher> <id>" >&2
  exit 2
fi

RESULT_HTML=$(bash "$SCRIPT_DIR/lib/search-books-or-jp.sh" "$TITLE" "$PUBLISHER") || exit 1

# 検索結果は <h3 class ="result_list_item"> 〜 </h3> が1件分のブロック。
# 出版社が一致する行のコードを集め、13桁のISBN(紙の本)を優先し、
# 無ければ20桁のJP-eコード(電子版)を使う。書影は版が違っても同じ表紙のことが多い。
MATCHES=$(printf '%s\n' "$RESULT_HTML" | awk -v pub="$PUBLISHER" '
  /<h3 class ="result_list_item">/ { code = ""; is_target_publisher = 0 }
  /href = "\/book-details\// {
    line = $0
    sub(/.*book-details\//, "", line)
    sub(/".*/, "", line)
    code = line
  }
  /result_list_discription_publisher/ {
    if (index($0, "出版社：" pub) > 0) { is_target_publisher = 1 }
  }
  /<\/h3>/ {
    if (is_target_publisher && length(code) == 13) { print "13:" code }
    else if (is_target_publisher && length(code) == 20) { print "20:" code }
  }
')

CODE=$(printf '%s\n' "$MATCHES" | grep '^13:' | head -1 | cut -d: -f2)
if [ -z "$CODE" ]; then
  CODE=$(printf '%s\n' "$MATCHES" | grep '^20:' | head -1 | cut -d: -f2)
fi

if [ -z "$CODE" ]; then
  echo "NOT_FOUND: '$TITLE'(出版社: $PUBLISHER)が見つかりませんでした" >&2
  exit 1
fi

THUMB_URL="https://thumbnail-s.images.books.or.jp/$CODE.jpg"
CONTENT_TYPE=$(curl -sI -A "Mozilla/5.0" -e "https://www.books.or.jp/" "$THUMB_URL" | grep -i '^content-type:' | tr -d '\r')

if [[ "$CONTENT_TYPE" != *"image/jpeg"* ]]; then
  echo "NOT_FOUND: 書影画像が取得できませんでした(code=$CODE): $THUMB_URL" >&2
  exit 1
fi

OUT_DIR="$REPO_ROOT/web/static/content/reviews/images"
mkdir -p "$OUT_DIR"
JPG_PATH="$OUT_DIR/$ID.jpg"
WEBP_PATH="$OUT_DIR/$ID.webp"

curl -sL -A "Mozilla/5.0" -e "https://www.books.or.jp/" "$THUMB_URL" -o "$JPG_PATH"

if ! cwebp -quiet -q 80 "$JPG_PATH" -o "$WEBP_PATH" 2>/dev/null; then
  echo "ERROR: webp変換に失敗しました($JPG_PATH)" >&2
  rm -f "$JPG_PATH" "$WEBP_PATH"
  exit 1
fi

echo "images/$ID.jpg"
