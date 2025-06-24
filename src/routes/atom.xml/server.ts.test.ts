import { describe, it, expect, vi, beforeEach } from 'vitest';
import { GET } from './+server';
import { XMLParser } from 'fast-xml-parser';
import { generateDescriptionFromText, getAllHTMLData } from '$lib/utils';

// モックのインポートをセットアップ
vi.mock('$lib/utils', () => {
  return {
    generateDescriptionFromText: vi.fn((text) => `Description for: ${text.substring(0, 50)}...`),
    getAllHTMLData: vi.fn()
  };
});

describe('/atom.xml endpoint', () => {
  // テスト用のダミー記事データ
  const mockArticles = [
    {
      id: 'article-1',
      title: 'テスト記事1',
      body: 'これはテスト記事1の本文です。',
      image: '/images/article-1.jpg',
      publishedAt: new Date('2023-01-10T09:00:00Z'),
      updatedAt: new Date('2023-01-15T14:30:00Z')
    },
    {
      id: 'article-2',
      title: 'テスト記事2',
      body: 'これはテスト記事2の本文です。これは長い本文のテストです。',
      image: '/images/article-2.png',
      publishedAt: new Date('2023-02-20T10:00:00Z'),
      updatedAt: new Date('2023-02-25T15:45:00Z')
    }
  ];

  // XMLパーサーの設定
  const parser = new XMLParser({
    ignoreAttributes: false,
    attributeNamePrefix: '@_',
    isArray: (name) => name === 'entry',
    cdataPropName: '__cdata'
  });

  // 各テストの前にモックをリセット
  beforeEach(() => {
    vi.resetAllMocks();
    // モックの getAllHTMLData が記事データを返すように設定
    vi.mocked(getAllHTMLData).mockResolvedValue([...mockArticles]);
  });

  it('Content-Typeヘッダーが正しい', async () => {
    const mockSetHeaders = vi.fn();

    await GET({ setHeaders: mockSetHeaders });

    expect(mockSetHeaders).toHaveBeenCalledWith({
      'Content-Type': 'application/xml'
    });
  });

  it('有効なAtomフィードを返す', async () => {
    const mockSetHeaders = vi.fn();

    const response = await GET({ setHeaders: mockSetHeaders });
    const xmlContent = await response.text();

    // XMLが解析可能かチェック
    expect(() => parser.parse(xmlContent)).not.toThrow();
  });

  it('フィードのメタデータが正しい', async () => {
    const mockSetHeaders = vi.fn();

    const response = await GET({ setHeaders: mockSetHeaders });
    const xmlContent = await response.text();
    const parsedXml = parser.parse(xmlContent);

    const feed = parsedXml.feed;

    // フィードの基本情報をチェック
    expect(feed.title).toBe('nagutabbyの考え事');
    expect(feed.link).toEqual(expect.arrayContaining([
      expect.objectContaining({
        '@_href': 'https://blog.nagutabby.uk'
      }),
      expect.objectContaining({
        '@_href': 'https://blog.nagutabby.uk/atom.xml',
        '@_type': 'application/rss+xml',
        '@_rel': 'self'
      })
    ]));
    expect(feed.author.name).toBe('nagutabby');
    expect(feed.id).toBe('tag:blog.nagutabby.uk,2023-01-01:/');
  });

  it('最新の更新日時が正しく設定される', async () => {
    const mockSetHeaders = vi.fn();

    const response = await GET({ setHeaders: mockSetHeaders });
    const xmlContent = await response.text();
    const parsedXml = parser.parse(xmlContent);

    // 最新の更新日時は2番目の記事のupdatedAt
    const latestDate = new Date(mockArticles[1].updatedAt).toISOString();
    expect(parsedXml.feed.updated).toBe(latestDate);
  });

  it('すべての記事エントリが含まれている', async () => {
    const mockSetHeaders = vi.fn();

    const response = await GET({ setHeaders: mockSetHeaders });
    const xmlContent = await response.text();
    const parsedXml = parser.parse(xmlContent);

    // 記事エントリの数をチェック
    expect(parsedXml.feed.entry.length).toBe(2);
  });

  it('記事エントリのフォーマットが正しい', async () => {
    const mockSetHeaders = vi.fn();

    const response = await GET({ setHeaders: mockSetHeaders });
    const xmlContent = await response.text();
    const parsedXml = parser.parse(xmlContent);

    const firstEntry = parsedXml.feed.entry[0];
    const firstArticle = mockArticles[0];

    // 最初の記事エントリの内容をチェック
    expect(firstEntry.title).toBe(firstArticle.title);
    expect(firstEntry.summary.__cdata).toContain('Description for:');
    expect(firstEntry.link['@_href']).toBe(`https://blog.nagutabby.uk/articles/${firstArticle.id}`);
    expect(firstEntry.link['@_rel']).toBe('alternate');
    expect(firstEntry.published).toBe(new Date(firstArticle.publishedAt).toISOString());
    expect(firstEntry.updated).toBe(new Date(firstArticle.updatedAt).toISOString());
    expect(firstEntry.id).toBe(`tag:blog.nagutabby.uk,${new Date(firstArticle.publishedAt).toISOString().substring(0, 10)}:/articles/${firstArticle.id}`);
  });

  it('記事がない場合でも有効なフィードを返す', async () => {
    vi.mocked(getAllHTMLData).mockResolvedValue([]);

    const mockSetHeaders = vi.fn();

    const response = await GET({ setHeaders: mockSetHeaders });
    const xmlContent = await response.text();
    const parsedXml = parser.parse(xmlContent);

    // 記事がなくても基本的なフィード情報があること
    expect(parsedXml.feed.title).toBe('nagutabbyの考え事');
    expect(parsedXml.feed.entry).toBeUndefined(); // 記事がないのでentryはundefined

    console.log(parsedXml);
    // 現在の日時が使われているはず
    const now = new Date();
    const todayStr = now.toISOString().substring(0, 10); // YYYY-MM-DD部分
    expect(parsedXml.feed.updated).toMatch(new RegExp(`^${todayStr}`));
  });

  it('generateDescriptionFromTextが各記事に対して呼び出される', async () => {
    const mockSetHeaders = vi.fn();

    await GET({ setHeaders: mockSetHeaders });

    // generateDescriptionFromTextが各記事の本文で呼び出されたことを確認
    expect(generateDescriptionFromText).toHaveBeenCalledTimes(2);
    expect(generateDescriptionFromText).toHaveBeenCalledWith(mockArticles[0].body);
    expect(generateDescriptionFromText).toHaveBeenCalledWith(mockArticles[1].body);
  });
});
