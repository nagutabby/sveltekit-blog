import { describe, it, expect } from 'vitest';
import { generateDescriptionFromText, getWebpPath } from './utils';

describe('generateDescriptionFromText', () => {
  it('HTMLタグを除去する', () => {
    expect(generateDescriptionFromText('<p>こんにちは</p>')).toBe('こんにちは');
  });

  it('100文字を超える場合は省略記号を付ける', () => {
    const long = 'あ'.repeat(150);
    const result = generateDescriptionFromText(long);
    expect(result.endsWith('…')).toBe(true);
    expect(result.length).toBe(101);
  });
});

describe('getWebpPath', () => {
  it('拡張子をwebpに変換する', () => {
    expect(getWebpPath('/images/foo.png')).toBe('/images/foo.webp');
  });

  it('http/dataスキームはそのまま返す', () => {
    expect(getWebpPath('https://example.com/foo.png')).toBe('https://example.com/foo.png');
  });

  it('空文字はそのまま返す', () => {
    expect(getWebpPath('')).toBe('');
  });
});
