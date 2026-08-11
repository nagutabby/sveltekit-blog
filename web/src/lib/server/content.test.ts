import { describe, it, expect, vi, beforeEach } from 'vitest';

const listArticles = vi.fn();
const getArticle = vi.fn();
const listReviews = vi.fn();
const getReview = vi.fn();

vi.mock('$lib/server/backendClient', () => ({
  createBackendClient: () => ({ listArticles, getArticle, listReviews, getReview })
}));

const { getAllRawData, getAllHTMLData, getHTMLData } = await import('./content');

const articlePb = (overrides: Partial<Record<string, unknown>> = {}) => ({
  id: 'my-article',
  title: 'タイトル',
  image: '/content/articles/images/foo.png',
  body: '# 見出し',
  publishedAt: '2025-06-15T00:00:00Z',
  updatedAt: '2025-06-16T00:00:00Z',
  ...overrides
});

const reviewPb = (overrides: Partial<Record<string, unknown>> = {}) => ({
  id: 'my-review',
  title: '本のタイトル',
  description: 'あらすじ',
  jpECode: '1234567890123',
  image: '/content/reviews/images/foo.jpg',
  rating: 5,
  body: '## 概要',
  publishedAt: '2025-03-01T00:00:00Z',
  updatedAt: '2025-03-02T00:00:00Z',
  ...overrides
});

describe('getAllRawData', () => {
  beforeEach(() => {
    listArticles.mockReset();
    listReviews.mockReset();
  });

  it('記事一覧をArticle型にマッピングする', async () => {
    listArticles.mockResolvedValue({ articles: [articlePb()] });

    const result = await getAllRawData('articles');

    expect(result).toEqual([
      {
        id: 'my-article',
        body: '# 見出し',
        title: 'タイトル',
        image: '/content/articles/images/foo.png',
        publishedAt: new Date('2025-06-15T00:00:00Z'),
        updatedAt: new Date('2025-06-16T00:00:00Z')
      }
    ]);
  });

  it('レビュー一覧をReview型にマッピングする(jpECode -> jp_e_code)', async () => {
    listReviews.mockResolvedValue({ reviews: [reviewPb()] });

    const result = await getAllRawData('reviews');

    expect(result).toEqual([
      {
        id: 'my-review',
        body: '## 概要',
        title: '本のタイトル',
        description: 'あらすじ',
        jp_e_code: '1234567890123',
        image: '/content/reviews/images/foo.jpg',
        rating: 5,
        publishedAt: new Date('2025-03-01T00:00:00Z'),
        updatedAt: new Date('2025-03-02T00:00:00Z')
      }
    ]);
  });
});

describe('getAllHTMLData', () => {
  it('bodyをMarkdownからHTMLへ変換する', async () => {
    listArticles.mockResolvedValue({ articles: [articlePb({ body: '# hello' })] });

    const [result] = await getAllHTMLData('articles');

    expect(result.body).toContain('<h1');
  });
});

describe('getHTMLData', () => {
  beforeEach(() => {
    getArticle.mockReset();
    getReview.mockReset();
  });

  it('記事をHTML変換済みのArticleとして返す', async () => {
    getArticle.mockResolvedValue({ article: articlePb({ body: '# hello' }) });

    const result = await getHTMLData('my-article', 'articles');

    expect(result.title).toBe('タイトル');
    expect(result.body).toContain('<h1');
  });

  it('存在しない記事は404を投げる', async () => {
    const { ConnectError, Code } = await import('@connectrpc/connect');
    getArticle.mockRejectedValue(new ConnectError('not found', Code.NotFound));

    await expect(getHTMLData('missing', 'articles')).rejects.toMatchObject({ status: 404 });
  });

  it('バックエンドのエラーは500を投げる', async () => {
    getArticle.mockRejectedValue(new Error('boom'));

    await expect(getHTMLData('my-article', 'articles')).rejects.toMatchObject({ status: 500 });
  });
});
