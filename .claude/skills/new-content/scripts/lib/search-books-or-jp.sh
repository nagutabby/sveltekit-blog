#!/usr/bin/env bash
# books.or.jp(日本書籍出版協会 出版書誌データベース)をタイトル・出版社で検索し、
# 検索結果ページの生HTMLを標準出力に返す内部ヘルパー。
# lookup-jp-e-code.sh / fetch-review-cover.sh から共通利用する。
#
# Usage: search-books-or-jp.sh <title> <publisher>

set -uo pipefail

TITLE="${1:-}"
PUBLISHER="${2:-}"

if [ -z "$TITLE" ] || [ -z "$PUBLISHER" ]; then
  echo "Usage: $0 <title> <publisher>" >&2
  exit 2
fi

COOKIE_JAR=$(mktemp)
trap 'rm -f "$COOKIE_JAR"' EXIT

TOP_HTML=$(curl -sL -c "$COOKIE_JAR" -A "Mozilla/5.0" "https://www.books.or.jp/")
TOKEN=$(printf '%s\n' "$TOP_HTML" | grep -o 'name="_token" value="[^"]*"' | head -1 | sed -E 's/.*value="([^"]*)"/\1/')

if [ -z "$TOKEN" ]; then
  echo "ERROR: books.or.jp からCSRFトークンを取得できませんでした" >&2
  exit 1
fi

curl -sL -b "$COOKIE_JAR" -c "$COOKIE_JAR" -A "Mozilla/5.0" \
  -X POST "https://www.books.or.jp/search-results" \
  --data-urlencode "_token=$TOKEN" \
  --data-urlencode "searchforbooks_title=$TITLE" \
  --data-urlencode "searchforbooks_publisher=$PUBLISHER" \
  --data-urlencode "publishtype1=on" \
  --data-urlencode "publishtype2=on" \
  --data-urlencode "publishtype3=on" \
  --data-urlencode "publishtype4=on" \
  --data-urlencode "accessible_search_flag=0" \
  --data-urlencode "first_books_search_flag=1"
