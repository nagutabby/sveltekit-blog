#!/usr/bin/env bash
# 書籍のタイトルと出版社から、日本書籍出版協会の出版書誌データベース(books.or.jp)を検索し、
# 電子版(JP-e)コードをfrontmatterのjp_e_code用に取得する。
#
# Usage: lookup-jp-e-code.sh <title> <publisher>
#
# 見つかった場合は標準出力にコード(20桁の数字)を1行出力して終了コード0。
# 見つからない場合(電子版が存在しない、該当書籍がヒットしないなど)は何も出力せず、
# 標準エラーに理由を出して終了コード1を返す。呼び出し側はこの場合jp_e_codeを空文字にしてよい。

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

TITLE="${1:-}"
PUBLISHER="${2:-}"

if [ -z "$TITLE" ] || [ -z "$PUBLISHER" ]; then
  echo "Usage: $0 <title> <publisher>" >&2
  exit 2
fi

RESULT_HTML=$(bash "$SCRIPT_DIR/lib/search-books-or-jp.sh" "$TITLE" "$PUBLISHER") || exit 1

# 検索結果は <h3 class ="result_list_item"> 〜 </h3> が1件分のブロック。
# ブロック内の /book-details/<code> と 出版社名、電子版バッジのactive有無を見て、
# 指定出版社の電子版(20桁のJP-eコード。13桁のISBNは対象外)のみ採用する。
CODE=$(printf '%s\n' "$RESULT_HTML" | awk -v pub="$PUBLISHER" '
  /<h3 class ="result_list_item">/ { code = ""; is_ebook = 0; is_target_publisher = 0 }
  /href = "\/book-details\// {
    line = $0
    sub(/.*book-details\//, "", line)
    sub(/".*/, "", line)
    code = line
  }
  /blue_box active/ { is_ebook = 1 }
  /result_list_discription_publisher/ {
    if (index($0, "出版社：" pub) > 0) { is_target_publisher = 1 }
  }
  /<\/h3>/ {
    if (is_ebook && is_target_publisher && length(code) == 20) {
      print code
    }
  }
' | head -1)

if [ -z "$CODE" ]; then
  echo "NOT_FOUND: '$TITLE'(出版社: $PUBLISHER)の電子版コードが見つかりませんでした(電子版が存在しない可能性があります)" >&2
  exit 1
fi

echo "$CODE"
