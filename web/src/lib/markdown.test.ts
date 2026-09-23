import { describe, expect, it } from 'vitest';
import { convertMarkdownToHtml } from './markdown';

describe('Markdown images', () => {
  it('ローカル画像をWebPの公開URLに変換する', () => {
    expect(convertMarkdownToHtml('![説明](images/photo.png)')).toContain(
      '<img src="/content/articles/images/photo.webp" alt="説明">'
    );
  });

  it('画像属性に含まれるHTML特殊文字をエスケープする', () => {
    const html = convertMarkdownToHtml('![a" onerror="alert(1)](https://example.com/a.png?x=1&y=2 "title & <tag>")');

    expect(html).toContain('alt="a&quot; onerror=&quot;alert(1)"');
    expect(html).toContain('src="https://example.com/a.png?x=1&amp;y=2"');
    expect(html).toContain('title="title &amp; &lt;tag&gt;"');
    expect(html).not.toContain(' onerror="');
  });
});
