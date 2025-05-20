import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/svelte';
import SlidesPage from './+page.svelte';
import '@testing-library/jest-dom/vitest';

// モックデータ
vi.mock('$app/state', () => {
  return {
    page: {
      data: {
        image: 'https://example.com/image.jpg',
        title: 'テストタイトル',
        body: 'テスト本文'
      }
    }
  };
});

// PDFThumbnailコンポーネントをモック
vi.mock('$lib/components/PDFThumbnail.svelte', () => ({
  default: vi.fn().mockImplementation((props) => {
    return { props };
  })
}));

// Headerコンポーネントをモック
// PDFThumbnailコンポーネントをモック
vi.mock('$lib/components/Header.svelte', () => ({
  default: vi.fn().mockImplementation((props) => {
    return { props };
  })
}));

// OpenGraphコンポーネントをモック
vi.mock('$lib/components/OpenGraph.svelte', () => ({
  default: vi.fn().mockImplementation((props) => {
    return { props };
  })
}));

describe('スライド一覧', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it('PDFが存在する場合、PDFサムネイルが表示される', () => {
    const mockData = {
      image: 'https://example.com/image.jpg',
      title: 'テストタイトル',
      body: 'テスト本文',
      pdfs: [
        { id: '1', url: 'https://example.com/pdf1.pdf' },
        { id: '2', url: 'https://example.com/pdf2.pdf' },
        { id: '3', url: 'https://example.com/pdf3.pdf' }
      ]
    };

    render(SlidesPage, { data: mockData });

    // グリッドが存在することを確認
    const gridElement = document.querySelector('.grid');
    expect(gridElement).not.toBeNull();

    // PDFの数だけアンカータグが存在することを確認
    const links = document.querySelectorAll('a');
    expect(links.length).toBe(3);

    // 各リンクが正しいhref属性を持っていることを確認
    expect(links[0].getAttribute('href')).toBe('/slides/1');
    expect(links[1].getAttribute('href')).toBe('/slides/2');
    expect(links[2].getAttribute('href')).toBe('/slides/3');
  });

  it('PDFが存在しない場合、メッセージが表示される', () => {
    const mockData = {
      image: 'https://example.com/image.jpg',
      title: 'テストタイトル',
      body: 'テスト本文',
      pdfs: []
    };

    render(SlidesPage, { data: mockData });

    // 「スライドがありません」というメッセージが表示されることを確認
    expect(screen.getByText('スライドがありません')).toBeInTheDocument();

    // グリッドが存在しないことを確認
    const gridElement = document.querySelector('.grid');
    expect(gridElement).toBeNull();
  });
});
