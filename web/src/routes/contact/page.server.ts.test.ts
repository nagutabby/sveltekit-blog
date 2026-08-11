import { describe, it, expect, vi, beforeEach } from 'vitest';
import { fail } from '@sveltejs/kit';

const submitContact = vi.fn();

vi.mock('$lib/server/backendClient', () => ({
  createBackendClient: () => ({ submitContact })
}));

const { actions } = await import('./+page.server');

describe('Contact form actions', () => {
  beforeEach(() => {
    submitContact.mockReset();
  });

  const buildRequest = (fields: Record<string, string>) => {
    const formData = new FormData();
    for (const [key, value] of Object.entries(fields)) {
      formData.append(key, value);
    }
    return { formData: () => Promise.resolve(formData) } as unknown as Request;
  };

  it('正常なフォームの送信が成功する', async () => {
    submitContact.mockResolvedValue({ errors: {} });

    const request = buildRequest({
      'im-robot': 'false',
      name: '氏名',
      email: 'test@example.com',
      text: 'これはテストメッセージです。'
    });

    const result = await actions.default({ request } as any);

    expect(result).toEqual({
      name: '氏名',
      email: 'test@example.com',
      text: 'これはテストメッセージです。'
    });
    expect(submitContact).toHaveBeenCalledWith({
      name: '氏名',
      email: 'test@example.com',
      text: 'これはテストメッセージです。',
      imRobot: false
    });
  });

  it('バックエンドがバリデーションエラーを返した場合は400を返す', async () => {
    submitContact.mockResolvedValue({
      errors: { imRobot: 'Botによるメッセージ送信はできません' }
    });

    const request = buildRequest({
      'im-robot': 'true',
      name: '氏名',
      email: 'test@example.com',
      text: 'これはテストメッセージです。'
    });

    const result = await actions.default({ request } as any);

    expect(result).toEqual(
      fail(400, {
        errors: { imRobot: 'Botによるメッセージ送信はできません' },
        values: {
          name: '氏名',
          email: 'test@example.com',
          text: 'これはテストメッセージです。'
        }
      })
    );
  });

  it('複数のバリデーションエラーがある場合', async () => {
    submitContact.mockResolvedValue({
      errors: {
        name: '氏名は必須です',
        email: 'メールアドレスは必須です',
        text: '本文は必須です'
      }
    });

    const request = buildRequest({ 'im-robot': 'false', name: '', email: '', text: '' });

    const result = await actions.default({ request } as any);

    expect(result).toEqual(
      fail(400, {
        errors: {
          name: '氏名は必須です',
          email: 'メールアドレスは必須です',
          text: '本文は必須です'
        },
        values: { name: '', email: '', text: '' }
      })
    );
  });

  it('バックエンド呼び出しがエラーを返した場合は500を返す', async () => {
    submitContact.mockRejectedValue(new Error('connect error'));

    const request = buildRequest({
      'im-robot': 'false',
      name: '氏名',
      email: 'test@example.com',
      text: 'これはテストメッセージです。'
    });

    const result = await actions.default({ request } as any);

    expect(result).toEqual(fail(500));
  });
});
