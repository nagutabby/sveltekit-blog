import { describe, expect, it } from 'vitest';
import { slidesFromFiles } from './slides';

describe('スライド一覧', () => {
  it('PDFだけを名前順にスライドURLへ変換する', () => {
    expect(slidesFromFiles(['zeta.pdf', 'cover.png', 'alpha.pdf', 'notes.txt'])).toEqual([
      { id: 'alpha', url: '/content/slides/alpha.pdf' },
      { id: 'zeta', url: '/content/slides/zeta.pdf' }
    ]);
  });
});
