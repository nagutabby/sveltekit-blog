import { describe, it, expect, beforeAll, afterAll, afterEach, beforeEach } from 'vitest';
import { setupServer } from 'msw/node';
import { http, HttpResponse } from 'msw';
import { load } from './+page.server';

const server = setupServer();

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }));
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe('slidesバケットのオブジェクトのメタデータ', () => {
  const mockFetch = async (url: string) => {
    const response = await fetch(url);
    return response;
  };
  it('取得した後に適切なデータフォーマットに変換される', async () => {
    const mockSlides = [
      {
        name: "slide1.pdf",
        size: 100000,
        lastModified: "2024-01-01T10:00:00.000Z"
      },
      {
        name: "slide2.pdf",
        size: 200000,
        lastModified: "2024-01-01T11:00:00.000Z"
      }
    ];
    
    server.use(
      http.get('/api/slides', () => {
        return HttpResponse.json(mockSlides);
      })
    );

    const result = await load({ fetch: mockFetch } as any);

    expect(result).toEqual({
      pdfUrls: ['slide1', 'slide2']
    });
  });

  it('取得できなかったときに適切なエラーが発生する', async () => {
    server.use(
      http.get('/api/slides', () => {
        return new HttpResponse(null, { status: 500 });
      })
    );

    await expect(load({ fetch: mockFetch } as any)).rejects.toThrow('スライドのメタデータの取得に失敗しました');
  });
});
