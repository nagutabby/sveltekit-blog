import { describe, expect, it, vi } from 'vitest';

const articles = Array.from({ length: 11 }, (_, i) => ({ id: `article-${i}` }));
vi.mock('$lib/server/content', () => ({ getAllHTMLData: vi.fn(async () => articles) }));

import { entries, load } from './+page.server';

describe('記事のページ分割', () => {
  it('2ページ目だけを事前生成し、残り1件を返す', async () => {
    expect(await entries()).toEqual([{ page: '2' }]);
    const result = await load({ params: { page: '2' } } as Parameters<typeof load>[0]);
    expect(result).toMatchObject({
      articles: [articles[10]],
      pagination: { currentPage: 2, totalPages: 2, hasNextPage: false }
    });
  });

  it.each(['1', '3', 'abc', '2.5'])('不正なページ番号 %s を404にする', async (page) => {
    await expect(load({ params: { page } } as Parameters<typeof load>[0])).rejects.toMatchObject({ status: 404 });
  });
});
