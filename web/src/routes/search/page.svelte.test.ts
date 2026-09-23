import { beforeEach, describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/svelte';
import '@testing-library/jest-dom/vitest';
import SearchPage from './+page.svelte';

const route = vi.hoisted(() => ({ url: new URL('https://example.com/search') }));
vi.mock('$app/state', () => ({ page: route }));
vi.mock('$app/environment', () => ({ browser: true }));
vi.mock('$lib/components/Timeline.svelte', () => ({ default: vi.fn().mockImplementation(() => ({})) }));
vi.mock('$lib/components/Header.svelte', () => ({ default: vi.fn().mockImplementation(() => ({})) }));
vi.mock('$lib/components/OpenGraph.svelte', () => ({ default: vi.fn().mockImplementation(() => ({})) }));

const articles = Array.from({ length: 11 }, (_, i) => ({
  id: `article-${i}`, title: `記事${i}`, image: '',
  rawBody: i === 10 ? '別の話題' : 'Svelte の話', body: '<p>本文</p>',
  publishedAt: new Date('2026-01-01')
}));
const data = { image: '/image.png', articles };

describe('検索結果', () => {
  beforeEach(() => { route.url = new URL('https://example.com/search'); });

  it('本文を大文字小文字を区別せず検索する', () => {
    route.url = new URL('https://example.com/search?q=svelte');
    render(SearchPage, { data } as any);

    expect(screen.getAllByRole('link', { name: /記事/ })).toHaveLength(10);
    expect(screen.queryByText('記事10')).not.toBeInTheDocument();
  });

  it('範囲外のページ番号を最終ページに収める', () => {
    route.url = new URL('https://example.com/search?page=99');
    render(SearchPage, { data } as any);

    expect(screen.getByText('記事10')).toBeInTheDocument();
    expect(screen.queryByText('記事0')).not.toBeInTheDocument();
  });
});
