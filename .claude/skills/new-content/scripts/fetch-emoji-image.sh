#!/usr/bin/env bash
# microsoft/fluentui-emoji から3D版の絵文字画像を取得し、article用の画像として配置する。
#
# Usage: fetch-emoji-image.sh <EmojiName>   (例: "White Flag")
#
# 成功時: web/static/content/articles/images/Microsoft-Fluentui-Emoji-3d-<Name>-3d.1024.png
#         と同名.webp を作成し、標準出力に "images/Microsoft-Fluentui-Emoji-3d-<Name>-3d.1024.png"
#         (frontmatterのimage用の相対パス)を出力して終了コード0。
# 失敗時(該当する絵文字がリポジトリに存在しない等): 何も出力せず、標準エラーに理由を出して
#         終了コード1。中間生成物は残さない。
#
# 前提: macOS(sips)での実行を想定している。1024x1024へのリサイズにsips、
# webp変換にcwebp(brew install webp)を使う。

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"

EMOJI_NAME="${1:-}"

if [ -z "$EMOJI_NAME" ]; then
  echo "Usage: $0 <EmojiName>" >&2
  exit 2
fi

HYPHEN=$(echo "$EMOJI_NAME" | tr ' ' '-')

PNG_FILE=$(gh api "repos/microsoft/fluentui-emoji/contents/assets/$EMOJI_NAME/3D" --jq '.[].name' 2>/dev/null | grep '\.png$' | head -1)

if [ -z "$PNG_FILE" ]; then
  echo "NOT_FOUND: fluentui-emoji に '$EMOJI_NAME' の3D画像が見つかりませんでした" >&2
  exit 1
fi

TMP_SOURCE=$(mktemp -t emoji_source).png
trap 'rm -f "$TMP_SOURCE"' EXIT

curl -sL "https://raw.githubusercontent.com/microsoft/fluentui-emoji/main/assets/$EMOJI_NAME/3D/$PNG_FILE" -o "$TMP_SOURCE"

if ! file "$TMP_SOURCE" | grep -q "PNG image data"; then
  echo "NOT_FOUND: '$EMOJI_NAME' の画像ダウンロードに失敗しました" >&2
  exit 1
fi

OUT_DIR="$REPO_ROOT/web/static/content/articles/images"
mkdir -p "$OUT_DIR"
BASENAME="Microsoft-Fluentui-Emoji-3d-${HYPHEN}-3d.1024"
PNG_PATH="$OUT_DIR/$BASENAME.png"
WEBP_PATH="$OUT_DIR/$BASENAME.webp"

sips -s format png -z 1024 1024 "$TMP_SOURCE" --out "$PNG_PATH" >/dev/null

if ! cwebp -quiet -q 80 "$PNG_PATH" -o "$WEBP_PATH" 2>/dev/null; then
  echo "ERROR: webp変換に失敗しました($PNG_PATH)" >&2
  rm -f "$PNG_PATH" "$WEBP_PATH"
  exit 1
fi

echo "images/$BASENAME.png"
