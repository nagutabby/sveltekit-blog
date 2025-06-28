import { GET } from './+server';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { getAllRawData } from '$lib/utils.js';
import { XMLParser } from 'fast-xml-parser';
import type { Article } from '$lib/types/blog';

beforeEach(() => {
  vi.resetAllMocks();

  vi.mock('$lib/utils.js', () => ({
    getAllRawData: vi.fn(),
  }));
});

// 各テスト後にモックをクリア
afterEach(() => {
  vi.clearAllMocks();
});

describe('GET /sitemap.xml', () => {
  it('生成されたXMLが正しい構造を持つ', async () => {
    const mockSetHeaders = vi.fn();

    // モックデータを準備（Article型に準拠）
    const mockArticles: Article[] = [
      {
        id: 'article1',
        body: 'This is article 1 content.',
        title: 'Article One Title',
        image: '/images/article1.jpg',
        publishedAt: new Date('2023-04-01'),
        updatedAt: new Date('2023-04-01')
      },
      {
        id: 'article2',
        body: 'This is article 2 content.',
        title: 'Article Two Title',
        image: '/images/article2.jpg',
        publishedAt: new Date('2023-04-02'),
        updatedAt: new Date('2023-04-02')
      }
    ];

    // モック関数を設定
    vi.mocked(getAllRawData).mockImplementation((type) => {
      if (type === "articles") {
        return Promise.resolve([...mockArticles]);
      } else if (type === "reviews") {
        return Promise.resolve([]);
      }

      return Promise.resolve([]);
    });

    // GET関数を呼び出し
    const response = await GET({ setHeaders: mockSetHeaders });

    expect(mockSetHeaders).toHaveBeenCalledWith({
      'Content-Type': 'application/xml'
    });

    // XML文字列を取得
    const xmlString = await response.text();

    // XMLパーサーを初期化
    const parser = new XMLParser();
    const xmlObject = parser.parse(xmlString);

    // ルート要素を確認
    expect(xmlObject.urlset).toBeDefined();

    // URLエントリの数を確認
    expect(xmlObject.urlset.url.length).toBe(2);

    // 各エントリを検証
    expect(xmlObject.urlset.url[0].loc).toBe('https://blog.nagutabby.uk/articles/article1');
    expect(xmlObject.urlset.url[0].lastmod).toBe('2023-04-01T00:00:00.000Z');

    expect(xmlObject.urlset.url[1].loc).toBe('https://blog.nagutabby.uk/articles/article2');
    expect(xmlObject.urlset.url[1].lastmod).toBe('2023-04-02T00:00:00.000Z');
  });

  it('更新日時がない記事の処理を確認', async () => {
    const mockSetHeaders = vi.fn();

    const mockArticles: Article[] = [
      {
        id: 'article3',
        body: 'No update date.',
        title: 'Article with no updatedAt',
        image: '/images/article3.jpg',
        publishedAt: new Date('2023-04-03'),
        updatedAt: null as any // 更新日時がnullの場合
      },
    ];

    vi.mocked(getAllRawData).mockResolvedValue(mockArticles);

    const response = await GET({ setHeaders: mockSetHeaders });
    const xmlString = await response.text();

    // '<lastmod>'タグが存在しないことを確認
    expect(xmlString).not.toContain('<lastmod>');
  });
});
