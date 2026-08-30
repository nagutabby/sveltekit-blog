#!/usr/bin/env bash
# microsoft/fluentui-emoji から2D(Color)版の絵文字画像を取得し、article用の画像として配置する。
#
# Usage: fetch-emoji-image.sh <EmojiName>   (例: "White Flag"。実際のリポジトリ側のフォルダ名は
#        "White flag" のように先頭のみ大文字のsentence caseだが、大文字小文字は問わず検索する)
#
# 成功時: web/static/content/articles/images/Microsoft-Fluentui-Emoji-Color-<Name>.512.png
#         と同名.webp を作成し、標準出力に "images/Microsoft-Fluentui-Emoji-Color-<Name>.512.png"
#         (frontmatterのimage用の相対パス)を出力して終了コード0。
# 失敗時(該当する絵文字がリポジトリに存在しない等): 何も出力せず、標準エラーに理由を出して
#         終了コード1。中間生成物は残さない。
#
# 前提: macOS(homebrew)での実行を想定している。fluentui-emojiの2D版はSVGのみ配布されているため、
# rsvg-convert(brew install librsvg)で512x512 PNGにラスタライズし、webp変換にcwebp(brew install webp)を使う。

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"

INPUT_NAME="${1:-}"

if [ -z "$INPUT_NAME" ]; then
  echo "Usage: $0 <EmojiName>" >&2
  exit 2
fi

# fluentui-emojiのフォルダ名は "White flag" のように先頭のみ大文字のsentence caseで、
# 入力の表記(Title Case等)と一致しないことが多いため、リポジトリ全体のパス一覧から
# 大文字小文字を無視して一致するColorフォルダを探す。
ALL_COLOR_SVGS=$(gh api "repos/microsoft/fluentui-emoji/git/trees/main?recursive=1" --jq '.tree[].path' 2>/dev/null | grep -E '^assets/[^/]+/Color/[^/]+\.svg$')

if [ -z "$ALL_COLOR_SVGS" ]; then
  echo "ERROR: fluentui-emojiのファイル一覧を取得できませんでした" >&2
  exit 1
fi

MATCH_PATH=$(printf '%s\n' "$ALL_COLOR_SVGS" | grep -iF "/${INPUT_NAME}/Color/" | head -1)

if [ -z "$MATCH_PATH" ]; then
  echo "NOT_FOUND: fluentui-emoji に '$INPUT_NAME' の2D(Color)画像が見つかりませんでした" >&2
  exit 1
fi

# MATCH_PATH の実例: assets/White flag/Color/white_flag_color.svg
REAL_NAME=$(printf '%s\n' "$MATCH_PATH" | sed -E 's#^assets/(.+)/Color/[^/]+\.svg$#\1#')
SVG_FILE=$(basename "$MATCH_PATH")

# ファイル名にはリポジトリの命名規則(Title-Case-Hyphenated)を使うため、入力側の表記を採用する。
HYPHEN=$(echo "$INPUT_NAME" | tr ' ' '-')

TMP_SOURCE=$(mktemp -t emoji_source).svg
trap 'rm -f "$TMP_SOURCE"' EXIT

# URLパスにはスペースをそのまま入れられないためパーセントエンコードする
ENCODED_NAME=$(printf '%s' "$REAL_NAME" | sed 's/ /%20/g')

curl -sL --max-time 20 "https://raw.githubusercontent.com/microsoft/fluentui-emoji/main/assets/$ENCODED_NAME/Color/$SVG_FILE" -o "$TMP_SOURCE"

if ! file "$TMP_SOURCE" | grep -qi "SVG"; then
  echo "NOT_FOUND: '$INPUT_NAME' の画像ダウンロードに失敗しました" >&2
  exit 1
fi

OUT_DIR="$REPO_ROOT/web/static/content/articles/images"
mkdir -p "$OUT_DIR"
BASENAME="Microsoft-Fluentui-Emoji-Color-${HYPHEN}.512"
PNG_PATH="$OUT_DIR/$BASENAME.png"
WEBP_PATH="$OUT_DIR/$BASENAME.webp"

if ! rsvg-convert -w 512 -h 512 -o "$PNG_PATH" "$TMP_SOURCE"; then
  echo "ERROR: SVGのラスタライズに失敗しました($INPUT_NAME)" >&2
  rm -f "$PNG_PATH"
  exit 1
fi

if ! cwebp -quiet -q 80 "$PNG_PATH" -o "$WEBP_PATH" 2>/dev/null; then
  echo "ERROR: webp変換に失敗しました($PNG_PATH)" >&2
  rm -f "$PNG_PATH" "$WEBP_PATH"
  exit 1
fi

echo "images/$BASENAME.png"
