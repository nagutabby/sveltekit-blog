import { describe, it, expect, afterAll } from 'vitest';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';

// content.ts reads CONTENT_DIR once at module load, so the fixture
// directory must exist and the env var must be set *before* the dynamic
// import below runs (a top-level `beforeAll` would run too late).
const contentDir = fs.mkdtempSync(path.join(os.tmpdir(), 'content-test-'));
process.env.CONTENT_DIR = contentDir;

const writeArticle = (id: string, frontmatter: Record<string, unknown>, body: string) => {
  const yaml = Object.entries(frontmatter)
    .map(([key, value]) => `${key}: ${value}`)
    .join('\n');
  fs.writeFileSync(path.join(contentDir, 'articles', `${id}.md`), `---\n${yaml}\n---\n${body}`);
};

const writeReview = (id: string, frontmatter: Record<string, unknown>, body: string) => {
  const yaml = Object.entries(frontmatter)
    .map(([key, value]) => `${key}: ${value}`)
    .join('\n');
  fs.writeFileSync(path.join(contentDir, 'reviews', `${id}.md`), `---\n${yaml}\n---\n${body}`);
};

fs.mkdirSync(path.join(contentDir, 'articles'), { recursive: true });
fs.mkdirSync(path.join(contentDir, 'reviews'), { recursive: true });

writeArticle(
  'my-article',
  { title: 'タイトル', image: 'images/foo.png', publishedAt: '2025-06-15', updatedAt: '2025-06-16' },
  '# 見出し'
);
writeReview(
  'my-review',
  {
    title: '本のタイトル',
    description: 'あらすじ',
    jp_e_code: '"1234567890123"',
    image: 'images/foo.jpg',
    rating: 5,
    publishedAt: '2025-03-01',
    updatedAt: '2025-03-02'
  },
  '## 概要'
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
        publishedAt: new Date('2025-06-15'),
        updatedAt: new Date('2025-06-16')
      }
    ]);
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
        publishedAt: new Date('2025-03-01'),
        updatedAt: new Date('2025-03-02')
      }
    ]);
  });
});

describe('getAllHTMLData', () => {
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

  it('frontmatterが壊れている記事は500を投げる', async () => {
    fs.writeFileSync(
      path.join(contentDir, 'articles', 'broken.md'),
      '---\ntitle: ["unterminated\n---\nbody'
    );

    await expect(getHTMLData('broken', 'articles')).rejects.toMatchObject({ status: 500 });
  });
});
