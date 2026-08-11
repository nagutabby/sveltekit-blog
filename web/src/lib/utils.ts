export const generateDescriptionFromText = (body: string) => {
  const hasHTMLTags = /<[a-z][\s\S]*>/i.test(body);

  const plainText: string = hasHTMLTags ? body.replace(/<[^>]+>/g, "") : body;

  const normalizedText = plainText.replace(/\s+/g, " ").trim();
  let description = normalizedText.slice(0, 100);

  if (normalizedText.length > 100) {
    const lastChar = description.slice(-1);
    if (!/[。、．，!！?？]/.test(lastChar)) {
      description += "…";
    }
  }

  return description;
};

export function getWebpPath(src: string) {
  if (!src) return '';

  // URLやデータURIの場合はそのまま返す
  if (src.startsWith('http') || src.startsWith('data:')) {
    return src;
  }

  const lastDotIndex = src.lastIndexOf('.');
  if (lastDotIndex === -1) return src;

  return src.substring(0, lastDotIndex) + '.webp';
}
