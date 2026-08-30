#!/usr/bin/env bash
# backend/content配下の記事/書評Markdownファイルについて、ファイル名とfrontmatterの構造が
# リポジトリの規約(backend/internal/content/loader.go, web/src/lib/types/blog.ts)と
# 一致しているかを検証する。
#
# Usage: validate-content.sh <path/to/backend/content/{articles,reviews}/YYYY-MM-DD[-N].md>
#
# 前提: macOS(BSD date環境)での実行を想定している。

set -uo pipefail

usage() {
  echo "Usage: $0 <path/to/backend/content/{articles,reviews}/YYYY-MM-DD[-N].md>" >&2
  exit 2
}

[ $# -eq 1 ] || usage
FILE="$1"

if [ ! -f "$FILE" ]; then
  echo "ERROR: file not found: $FILE" >&2
  exit 2
fi

REPO_ROOT=$(git -C "$(dirname "$FILE")" rev-parse --show-toplevel 2>/dev/null || pwd)
BASENAME=$(basename "$FILE")
PARENT_DIR=$(basename "$(cd "$(dirname "$FILE")" && pwd)")

errors=0
fail() {
  echo "ERROR: $1" >&2
  errors=$((errors + 1))
}

# 1. 配置ディレクトリ
TYPE=""
case "$PARENT_DIR" in
  articles | reviews)
    TYPE="$PARENT_DIR"
    ;;
  *)
    fail "backend/content/articles/ または backend/content/reviews/ 直下に配置してください(実際の配置先: $PARENT_DIR)"
    ;;
esac

# 2. ファイル名(ISO日付 + 任意の連番)
if [[ "$BASENAME" =~ ^([0-9]{4}-[0-9]{2}-[0-9]{2})(-[0-9]+)?\.md$ ]]; then
  DATE_PART="${BASH_REMATCH[1]}"
  if ! date -j -f "%Y-%m-%d" "$DATE_PART" >/dev/null 2>&1; then
    fail "ファイル名の日付が実在しません: $DATE_PART"
  fi
else
  fail "ファイル名はYYYY-MM-DD.mdまたはYYYY-MM-DD-N.mdの形式にしてください(実際: $BASENAME)"
fi

# 3. frontmatterの抽出(先頭の --- ... --- の間)
if [ "$(sed -n '1p' "$FILE")" != "---" ]; then
  fail "1行目が --- ではありません(frontmatterで始まっていません)"
fi

FRONTMATTER=$(awk '
  NR == 1 && $0 == "---" { infm = 1; next }
  infm && $0 == "---" { exit }
  infm { print }
' "$FILE")

if [ -z "$FRONTMATTER" ]; then
  fail "frontmatterが空か見つかりません"
fi

get_field() {
  echo "$FRONTMATTER" \
    | grep -E "^$1:" \
    | head -1 \
    | sed -E "s/^$1:[[:space:]]*//; s/^\"(.*)\"\$/\1/; s/^'(.*)'\$/\1/"
}

# 4. 種別ごとの必須フィールドと許容フィールド
# required_fields: 値が空であってはいけないフィールド
# present_only_fields: キー自体は必須だが、値が空でもよいフィールド(jp_e_codeは書誌コード未取得時に空を許容)
case "$TYPE" in
  articles)
    required_fields="id title image"
    present_only_fields=""
    allowed_fields="id title image"
    ;;
  reviews)
    required_fields="id title description image rating"
    present_only_fields="jp_e_code"
    allowed_fields="id title description jp_e_code image rating"
    ;;
  *)
    required_fields=""
    present_only_fields=""
    allowed_fields=""
    ;;
esac

if [ -n "$TYPE" ]; then
  for field in $required_fields; do
    value=$(get_field "$field")
    if [ -z "$value" ]; then
      fail "必須フィールド '$field' が空か欠落しています"
    fi
  done

  for field in $present_only_fields; do
    if ! echo "$FRONTMATTER" | grep -qE "^$field:"; then
      fail "フィールド '$field' がありません(値は空でも構いませんが、キー自体は必要です)"
    fi
  done

  present_fields=$(echo "$FRONTMATTER" | grep -oE '^[A-Za-z_]+:' | sed 's/:$//')
  for field in $present_fields; do
    match=0
    for allowed in $allowed_fields; do
      [ "$field" = "$allowed" ] && match=1 && break
    done
    if [ "$match" -eq 0 ]; then
      fail "$TYPE では想定していないフィールドです: '$field'"
    fi
  done
fi

# 5. idの形式(kebab-case)と一意性
ID=$(get_field id)
if [ -n "$ID" ]; then
  if [[ ! "$ID" =~ ^[a-z0-9]+(-[a-z0-9]+)*$ ]]; then
    fail "id は英小文字・数字をハイフンで繋いだkebab-caseにしてください(実際: '$ID')"
  fi

  count=$(grep -h -E "^id:[[:space:]]*\"?${ID}\"?[[:space:]]*\$" \
    "$REPO_ROOT"/backend/content/articles/*.md \
    "$REPO_ROOT"/backend/content/reviews/*.md 2>/dev/null | wc -l | tr -d ' ')
  if [ "$count" -gt 1 ]; then
    fail "id '$ID' が他のファイルと重複しています(${count}件ヒット)"
  fi
fi

# 6. rating(reviewのみ、1〜5の整数)
if [ "$TYPE" = "reviews" ]; then
  rating=$(get_field rating)
  if [ -n "$rating" ] && [[ ! "$rating" =~ ^[1-5]$ ]]; then
    fail "rating は1〜5の整数にしてください(実際: '$rating')"
  fi
fi

# 7. imageが指す画像ファイルの実在確認
# Card.svelte/Header.svelteは image の拡張子を.webpに置き換えたパスのみを<img src>に使う
# (web/src/lib/utils.ts の getWebpPath)。元画像(jpg/png)だけでなく.webpも無いと表示が壊れる。
IMAGE=$(get_field image)
if [ -n "$IMAGE" ] && [ -n "$TYPE" ]; then
  IMG_PATH="$REPO_ROOT/web/static/content/$TYPE/$IMAGE"
  WEBP_PATH="${IMG_PATH%.*}.webp"
  if [ ! -f "$IMG_PATH" ]; then
    fail "image が指す画像ファイルが存在しません: $IMG_PATH"
  fi
  if [ ! -f "$WEBP_PATH" ]; then
    fail "image の.webp版が存在しません(実際の配信に必要): $WEBP_PATH"
  fi
fi

if [ "$errors" -gt 0 ]; then
  echo "NG: $errors 件の問題が見つかりました: $FILE" >&2
  exit 1
fi

echo "OK: $FILE"
exit 0
