import { describe, it, expect, afterAll, vi } from 'vitest';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';

// content.ts reads CONTENT_DIR once at module load, so the fixture
// directory must exist and the env var must be set *before* the dynamic
// import below runs (a top-level `beforeAll` would run too late).
const contentDir = fs.mkdtempSync(path.join(os.tmpdir(), 'content-test-'));
process.env.CONTENT_DIR = contentDir;
vi.mock('$app/environment', () => ({ dev: false }));

const writeArticle = (
  filename: string,
  frontmatter: Record<string, unknown>,
  body: string
) => {
  const yaml = Object.entries(frontmatter)
    .map(([key, value]) => `${key}: ${value}`)
    .join('\n');
  fs.writeFileSync(path.join(contentDir, 'articles', filename), `---\n${yaml}\n---\n${body}`);
};

const writeReview = (
  filename: string,
  frontmatter: Record<string, unknown>,
  body: string
) => {
  const yaml = Object.entries(frontmatter)
    .map(([key, value]) => `${key}: ${value}`)
    .join('\n');
  fs.writeFileSync(path.join(contentDir, 'reviews', filename), `---\n${yaml}\n---\n${body}`);
};

fs.mkdirSync(path.join(contentDir, 'articles'), { recursive: true });
fs.mkdirSync(path.join(contentDir, 'reviews'), { recursive: true });

writeArticle(
  '2025-06-15.md',
  { id: 'my-article', title: 'タイトル', image: 'images/foo.png', is_draft: false },
  '# 見出し'
);
writeReview(
  '2025-03-01.md',
  {
    id: 'my-review',
    title: '本のタイトル',
    description: 'あらすじ',
    jp_e_code: '"1234567890123"',
    image: 'images/foo.jpg',
    rating: 5,
    is_draft: false
  },
  '## 概要'
);
writeArticle(
  '2025-06-16.md',
  { id: 'draft-article', title: '下書き', image: 'images/bar.png', is_draft: true },
  '# 下書き本文'
);
writeArticle(
  '2025-06-17.md',
  { id: 'missing-is-draft', title: 'is_draft未指定', image: 'images/baz.png' },
  '# 本文'
);

afterAll(() => {
  fs.rmSync(contentDir, { recursive: true, force: true });
});

const { getAllRawData, getAllHTMLData, getHTMLData } = await import('./content');

describe('getAllRawData', () => {
  it('記事一覧をArticle型にマッピングする', async () => {
    const result = await getAllRawData('articles');

    expect(result).toEqual([
      {
        id: 'my-article',
        body: '# 見出し',
        title: 'タイトル',
        image: '/content/articles/images/foo.png',
        publishedAt: new Date('2025-06-15')
      }
    ]);
  });

  it('is_draftがtrue、または未指定の記事は一覧から除外する', async () => {
    const result = await getAllRawData('articles');

    expect(result.map((a) => a.id)).toEqual(['my-article']);
  });

  it('レビュー一覧をReview型にマッピングする', async () => {
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
        publishedAt: new Date('2025-03-01')
      }
    ]);
  });
});

describe('getAllHTMLData', () => {
  it('キャッシュされた生データを変更せず、検索用のMarkdownを保持する', async () => {
    const raw = await getAllRawData('articles');
    const originalBody = raw[0].body;

    const html = await getAllHTMLData('articles');

    expect(html[0].body).toContain('<h1');
    expect((await getAllRawData('articles'))[0].body).toBe(originalBody);
  });

  it('bodyをMarkdownからHTMLへ変換する', async () => {
    const [result] = await getAllHTMLData('articles');

    expect(result.body).toContain('<h1');
  });
});

describe('getHTMLData', () => {
  it('記事をHTML変換済みのArticleとして返す', async () => {
    const result = await getHTMLData('my-article', 'articles');

    expect(result.title).toBe('タイトル');
    expect(result.body).toContain('<h1');
  });

  it('存在しない記事は404を投げる', async () => {
    await expect(getHTMLData('missing', 'articles')).rejects.toMatchObject({ status: 404 });
  });

  it('is_draftがtrueの記事は404を投げる', async () => {
    await expect(getHTMLData('draft-article', 'articles')).rejects.toMatchObject({ status: 404 });
  });

  it('is_draft未指定の記事は404を投げる(デフォルトtrue)', async () => {
    await expect(getHTMLData('missing-is-draft', 'articles')).rejects.toMatchObject({
      status: 404
    });
  });

  it('frontmatterが壊れている記事は500を投げる', async () => {
    fs.writeFileSync(
      path.join(contentDir, 'articles', '2025-01-01.md'),
      '---\ntitle: ["unterminated\n---\nbody'
    );

    await expect(getHTMLData('broken', 'articles')).rejects.toMatchObject({ status: 500 });
  });
});
