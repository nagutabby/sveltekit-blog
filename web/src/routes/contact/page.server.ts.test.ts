import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { actions } from './+page.server';
import { fail } from '@sveltejs/kit';

// 環境変数のモック
vi.mock('$env/static/private', () => ({
  EMAIL_API_TOKEN: 'mock_token',
  FROM_ADDRESS: 'admin@example.com',
  BCC_ADDRESS: 'bcc@example.com'
}));

// fetchのモック
global.fetch = vi.fn();

describe('Contact form actions', () => {
  const mockFormData = new FormData();

  beforeEach(() => {
    // テスト前にモックをリセット
    vi.resetAllMocks();

    // レスポンスのモック
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: vi.fn().mockResolvedValue({})
    } as unknown as Response);
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it('正常なフォームの送信が成功する', async () => {
    // 正常な入力データの準備
    const validFormData = new FormData();
    validFormData.append('im-robot', 'false');
    validFormData.append('name', '氏名');
    validFormData.append('email', 'test@example.com');
    validFormData.append('text', 'これはテストメッセージです。');

    const mockRequest = {
      formData: () => Promise.resolve(validFormData)
    };

    // アクションの実行
    const result = await actions.default({ request: mockRequest as Request } as any);

    // 期待される結果の検証
    expect(result).toEqual({
      name: '氏名',
      email: 'test@example.com',
      text: 'これはテストメッセージです。'
    });

    // fetchが正しく呼ばれたことを確認
    expect(fetch).toHaveBeenCalledTimes(1);
    expect(fetch).toHaveBeenCalledWith('https://send.api.mailtrap.io/api/send', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Api-Token': 'mock_token'
      },
      body: expect.any(String)
    });

    // fetchに渡されたボディが正しいことを確認
    const fetchCall = vi.mocked(fetch).mock.calls[0];
    expect(fetchCall[1]).toBeDefined();
    expect(fetchCall[1]!.body).toBeDefined();

    const requestBody = JSON.parse(fetchCall[1]!.body as string);
    expect(requestBody).toEqual({
      from: {
        email: 'admin@example.com',
        name: 'Hiroto Sasagawa'
      },
      to: [{
        email: 'test@example.com',
        name: '氏名'
      }],
      bcc: [{
        email: 'bcc@example.com',
        name: 'Hiroto Sasagawa'
      }],
      subject: 'お問い合わせを受け付けました',
      html: expect.stringContaining('氏名')
    });
  });

  it('ロボットチェックに失敗する場合のバリデーションエラー', async () => {
    // ロボットチェックが失敗するデータの準備
    const robotFormData = new FormData();
    robotFormData.append('im-robot', 'true');
    robotFormData.append('name', '氏名');
    robotFormData.append('email', 'test@example.com');
    robotFormData.append('text', 'これはテストメッセージです。');

    const mockRequest = {
      formData: () => Promise.resolve(robotFormData)
    };

    // アクションの実行
    const result = await actions.default({ request: mockRequest as Request } as any);

    // 期待される結果の検証
    expect(result).toEqual(
      fail(400, {
        errors: {
          imRobot: 'Botによるメッセージ送信はできません'
        },
        values: {
          name: '氏名',
          email: 'test@example.com',
          text: 'これはテストメッセージです。'
        }
      })
    );

    // fetchが呼ばれていないことを確認
    expect(fetch).not.toHaveBeenCalled();
  });

  it('必須フィールドが欠けている場合のバリデーションエラー', async () => {
    // 名前が欠けているデータの準備
    const invalidFormData = new FormData();
    invalidFormData.append('im-robot', 'false');
    invalidFormData.append('name', ''); // 空の名前
    invalidFormData.append('email', 'test@example.com');
    invalidFormData.append('text', 'これはテストメッセージです。');

    const mockRequest = {
      formData: () => Promise.resolve(invalidFormData)
    };

    // アクションの実行
    const result = await actions.default({ request: mockRequest as Request } as any);

    // 期待される結果の検証
    expect(result).toEqual(
      fail(400, {
        errors: {
          name: '氏名は必須です'
        },
        values: {
          name: '',
          email: 'test@example.com',
          text: 'これはテストメッセージです。'
        }
      })
    );

    // fetchが呼ばれていないことを確認
    expect(fetch).not.toHaveBeenCalled();
  });

  it('メールアドレスの形式が不正な場合のバリデーションエラー', async () => {
    // 無効なメールアドレスを含むデータの準備
    const invalidEmailFormData = new FormData();
    invalidEmailFormData.append('im-robot', 'false');
    invalidEmailFormData.append('name', '氏名');
    invalidEmailFormData.append('email', 'invalid-email');
    invalidEmailFormData.append('text', 'これはテストメッセージです。');

    const mockRequest = {
      formData: () => Promise.resolve(invalidEmailFormData)
    };

    // アクションの実行
    const result = await actions.default({ request: mockRequest as Request } as any);

    // 期待される結果の検証
    expect(result).toEqual(
      fail(400, {
        errors: {
          email: 'メールアドレスの形式が不適切です'
        },
        values: {
          name: '氏名',
          email: 'invalid-email',
          text: 'これはテストメッセージです。'
        }
      })
    );

    // fetchが呼ばれていないことを確認
    expect(fetch).not.toHaveBeenCalled();
  });

  it('複数のバリデーションエラーがある場合', async () => {
    // 複数のエラーを含むデータの準備
    const multipleErrorsFormData = new FormData();
    multipleErrorsFormData.append('im-robot', 'false');
    multipleErrorsFormData.append('name', '');
    multipleErrorsFormData.append('email', '');
    multipleErrorsFormData.append('text', '');

    const mockRequest = {
      formData: () => Promise.resolve(multipleErrorsFormData)
    };

    // アクションの実行
    const result = await actions.default({ request: mockRequest as Request } as any);

    // 期待される結果の検証
    expect(result).toEqual(
      fail(400, {
        errors: {
          name: '氏名は必須です',
          email: 'メールアドレスは必須です',
          text: '本文は必須です'
        },
        values: {
          name: '',
          email: '',
          text: ''
        }
      })
    );

    // fetchが呼ばれていないことを確認
    expect(fetch).not.toHaveBeenCalled();
  });

  it('メール送信APIが失敗した場合', async () => {
    // 正常な入力データの準備
    const validFormData = new FormData();
    validFormData.append('im-robot', 'false');
    validFormData.append('name', '氏名');
    validFormData.append('email', 'test@example.com');
    validFormData.append('text', 'これはテストメッセージです。');

    const mockRequest = {
      formData: () => Promise.resolve(validFormData)
    };

    // API呼び出しが失敗するようにモック
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      json: vi.fn().mockResolvedValue({ error: 'API error' })
    } as unknown as Response);

    // アクションの実行
    const result = await actions.default({ request: mockRequest as Request } as any);

    // 期待される結果の検証
    expect(result).toEqual(fail(500));

    // fetchが呼ばれたことを確認
    expect(fetch).toHaveBeenCalledTimes(1);
  });

  it('メール送信時に例外が発生した場合', async () => {
    // 正常な入力データの準備
    const validFormData = new FormData();
    validFormData.append('im-robot', 'false');
    validFormData.append('name', '氏名');
    validFormData.append('email', 'test@example.com');
    validFormData.append('text', 'これはテストメッセージです。');

    const mockRequest = {
      formData: () => Promise.resolve(validFormData)
    };

    // fetchが例外をスローするようにモック
    vi.mocked(fetch).mockRejectedValue(new Error('Network error'));

    // アクションの実行
    const result = await actions.default({ request: mockRequest as Request } as any);

    // 期待される結果の検証
    expect(result).toEqual(fail(500));

    // fetchが呼ばれたことを確認
    expect(fetch).toHaveBeenCalledTimes(1);
  });
});
