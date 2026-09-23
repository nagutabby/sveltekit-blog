import { describe, expect, it, vi } from 'vitest';

const reviews = Array.from({ length: 21 }, (_, i) => ({ id: `review-${i}` }));
vi.mock('$lib/server/content', () => ({ getAllHTMLData: vi.fn(async () => reviews) }));

import { entries, load } from './+page.server';

describe('書評のページ分割', () => {
  it('2・3ページ目を事前生成し、最終ページの1件を返す', async () => {
    expect(await entries()).toEqual([{ page: '2' }, { page: '3' }]);
    const result = await load({ params: { page: '3' } } as Parameters<typeof load>[0]);
    expect(result).toMatchObject({
      reviews: [reviews[20]],
      pagination: { currentPage: 3, totalPages: 3, hasNextPage: false }
    });
  });

  it.each(['1', '4', 'bad'])('不正なページ番号 %s を404にする', async (page) => {
    await expect(load({ params: { page } } as Parameters<typeof load>[0])).rejects.toMatchObject({ status: 404 });
  });
});
