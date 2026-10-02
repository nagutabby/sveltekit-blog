import { beforeEach, describe, expect, it } from 'vitest';
import { render, screen, waitFor } from '@testing-library/svelte';
import '@testing-library/jest-dom/vitest';
import SearchResults from '../../components/SearchResults.svelte';
import type { SearchableArticle } from '$lib/types/search';

const articles: SearchableArticle[] = Array.from({ length: 11 }, (_, index) => ({
  id: `article-${index}`,
  title: `記事${index}`,
  image: '',
  rawBody: index === 10 ? '別の話題' : 'Svelte の話',
  body: '<p>本文</p>'
}));

describe('検索結果', () => {
  beforeEach(() => window.history.replaceState({}, '', '/search'));

  it('本文を大文字小文字を区別せず検索する', async () => {
    window.history.replaceState({}, '', '/search?q=svelte');
    render(SearchResults, { articles });

    await waitFor(() => expect(screen.getAllByRole('link', { name: /記事/ })).toHaveLength(10));
    expect(screen.queryByText('記事10')).not.toBeInTheDocument();
  });

  it('範囲外のページ番号を最終ページに収める', async () => {
    window.history.replaceState({}, '', '/search?page=99');
    render(SearchResults, { articles });

    await waitFor(() => expect(screen.getByText('記事10')).toBeInTheDocument());
    expect(screen.queryByText('記事0')).not.toBeInTheDocument();
  });
});
