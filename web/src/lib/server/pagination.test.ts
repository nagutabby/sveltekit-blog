import { describe, expect, it } from 'vitest';
import { paginate, totalPagesFor } from './pagination';

describe('pagination', () => {
  it('空・ちょうど1ページ・次ページありの件数を計算する', () => {
    expect(totalPagesFor(0, 10)).toBe(0);
    expect(totalPagesFor(10, 10)).toBe(1);
    expect(totalPagesFor(11, 10)).toBe(2);
  });

  it('最終ページの残り1件と前後のページ状態を返す', () => {
    const result = paginate(Array.from({ length: 11 }, (_, i) => i), 2, 10);
    expect(result).toEqual({
      items: [10], currentPage: 2, totalPages: 2,
      hasNextPage: false, hasPrevPage: true
    });
  });

  it('空の一覧でも余分なページを作らない', () => {
    expect(paginate([], 1, 10)).toEqual({
      items: [], currentPage: 1, totalPages: 0,
      hasNextPage: false, hasPrevPage: false
    });
  });
});
