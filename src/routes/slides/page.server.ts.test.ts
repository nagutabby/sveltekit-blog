import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { load } from './+page.server';
import fs from 'fs';
import path from 'path';
import { error } from '@sveltejs/kit';

// fsとpathモジュールをモック化
vi.mock('fs');
vi.mock('path');
vi.mock('@sveltejs/kit', () => ({
  error: vi.fn()
}));

describe('スライド一覧', () => {
  // 各テストの前にモックをリセット
  beforeEach(() => {
    vi.resetAllMocks();

    // path.joinのモック実装
    vi.mocked(path.join).mockReturnValue('static/content/slides');
  });

  // 各テスト後にモックをクリア
  afterEach(() => {
    vi.clearAllMocks();
  });

  it('スライドディレクトリが存在する場合、PDFファイルのリストを返すこと', async () => {
    // existsSyncがtrueを返すようにモック化
    vi.mocked(fs.existsSync).mockReturnValue(true);

    // readdirSyncのモック実装 - PDFファイル名のリストを返す
    vi.mocked(fs.readdirSync).mockReturnValue(['slide1.pdf', 'slide2.pdf', 'presentation.pdf'] as unknown as fs.Dirent[]);

    const result = await load({} as any);

    // path.joinが正しく呼ばれたか確認
    expect(path.join).toHaveBeenCalledWith('static/content/slides');

    // fs.existsSyncが正しく呼ばれたか確認
    expect(fs.existsSync).toHaveBeenCalledWith('static/content/slides');

    // fs.readdirSyncが正しく呼ばれたか確認
    expect(fs.readdirSync).toHaveBeenCalledWith('static/content/slides');

    // 結果が期待通りか確認
    expect(result).toEqual({
      image: "images/Microsoft-Fluentui-Emoji-3d-Package-3d.1024.png",
      title: "スライド一覧",
      body: "LT会で使ったスライドをまとめています",
      pdfs: [
        { id: 'slide1', url: '/content/slides/slide1.pdf' },
        { id: 'slide2', url: '/content/slides/slide2.pdf' },
        { id: 'presentation', url: '/content/slides/presentation.pdf' }
      ]
    });
  });

  it('スライドディレクトリが存在しない場合、エラーをスローすること', async () => {
    // existsSyncがfalseを返すようにモック化
    vi.mocked(fs.existsSync).mockReturnValue(false);

    // errorがエラーオブジェクトを返すようにモック化
    vi.mocked(error).mockImplementation((code, message) => {
      throw new Error(`${code}: ${message}`);
    });

    // load関数がエラーをスローすることを期待
    await expect(load({} as any)).rejects.toThrow('500: スライドディレクトリが見つかりません');

    // path.joinが正しく呼ばれたか確認
    expect(path.join).toHaveBeenCalledWith('static/content/slides');

    // fs.existsSyncが正しく呼ばれたか確認
    expect(fs.existsSync).toHaveBeenCalledWith('static/content/slides');

    // @sveltejs/kitのerror関数が正しく呼ばれたか確認
    expect(error).toHaveBeenCalledWith(500, 'スライドディレクトリが見つかりません');
  });

  it('PDFファイル以外のファイルがフィルタリングされること', async () => {
    // existsSyncがtrueを返すようにモック化
    vi.mocked(fs.existsSync).mockReturnValue(true);

    // readdirSyncのモック実装 - 様々な拡張子のファイル名を含むリストを返す
    vi.mocked(fs.readdirSync).mockReturnValue(['slide1.pdf', 'image.jpg', 'document.txt', 'slide2.pdf'] as unknown as fs.Dirent[]);

    const result: any = await load({} as any);

    expect(result.pdfs).toHaveLength(2);
    expect(result.pdfs).toEqual([
      { id: 'slide1', url: '/content/slides/slide1.pdf' },
      { id: 'slide2', url: '/content/slides/slide2.pdf' }
    ]);
  });
});
