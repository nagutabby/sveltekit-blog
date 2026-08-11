import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import ContactForm from './+page.svelte';
import { tick } from 'svelte';
import '@testing-library/jest-dom/vitest';

const submitContact = vi.fn();

vi.mock('$lib/client/backendClient', () => ({
  createBackendClient: () => ({ submitContact })
}));

describe('お問い合わせフォーム', () => {
  beforeEach(() => {
    vi.resetAllMocks();

    vi.mock('$lib/components/Header.svelte', () => ({
      default: vi.fn().mockImplementation((props) => {
        return { props };
      })
    }));

    vi.mock('$lib/components/OpenGraph.svelte', () => ({
      default: vi.fn().mockImplementation((props) => {
        return { props };
      })
    }));
  });

  const fillForm = () => {
    fireEvent.input(screen.getByRole('textbox', { name: '氏名' }), { target: { value: '氏名' } });
    fireEvent.input(screen.getByRole('textbox', { name: 'メールアドレス' }), { target: { value: 'user@example.com' } });
    fireEvent.input(screen.getByRole('textbox', { name: '本文' }), { target: { value: '本文' } });
  };

  it('適切にレンダリングされる', () => {
    render(ContactForm);

    expect(screen.getByRole('textbox', { name: '氏名' })).toBeInTheDocument();
    expect(screen.getByRole('textbox', { name: 'メールアドレス' })).toBeInTheDocument();
    expect(screen.getByRole('textbox', { name: '本文' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '送信' })).toBeInTheDocument();

    const honeyPot = screen.getByRole('checkbox', { hidden: true, name: '私はロボットです' });
    expect(honeyPot.closest('.hidden')).toBeTruthy();
  });

  it('適切な値を送信したときに成功メッセージが表示される', async () => {
    submitContact.mockResolvedValue({ errors: {} });
    render(ContactForm);
    fillForm();

    await fireEvent.click(screen.getByRole('button', { name: '送信' }));

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('お問い合わせを受け付けました。ご連絡までしばらくお待ちください。');
    });
    expect(submitContact).toHaveBeenCalledWith({
      name: '氏名',
      email: 'user@example.com',
      text: '本文',
      imRobot: false
    });
  });

  it('バックエンドがエラーを返したときにエラーメッセージが表示される', async () => {
    submitContact.mockResolvedValue({
      errors: {
        name: '氏名は必須です',
        email: 'メールアドレスは必須です',
        text: '本文は必須です'
      }
    });
    render(ContactForm);
    fillForm();

    await fireEvent.click(screen.getByRole('button', { name: '送信' }));

    await waitFor(() => {
      const alert = screen.getByRole('alert');
      expect(alert).toHaveTextContent('氏名は必須です');
      expect(alert).toHaveTextContent('メールアドレスは必須です');
      expect(alert).toHaveTextContent('本文は必須です');
    });
  });

  it('送信中は送信ボタンのテキストが変化する', async () => {
    let resolveSubmit: (value: { errors: Record<string, string> }) => void;
    submitContact.mockReturnValue(
      new Promise((resolve) => {
        resolveSubmit = resolve;
      })
    );
    render(ContactForm);
    fillForm();

    const submitButton = screen.getByRole('button', { name: '送信' });
    await fireEvent.click(submitButton);
    expect(screen.getByRole('button', { name: '送信中…' })).toBeInTheDocument();

    resolveSubmit!({ errors: {} });

    await waitFor(() => {
      expect(screen.getByRole('button', { name: '送信' })).toBeInTheDocument();
    });
  });

  it('空文字列を送信しようとしたときにrequired属性によるバリデーションが機能する', async () => {
    render(ContactForm);

    const submitButton = screen.getByRole('button', { name: '送信' });
    fireEvent.click(submitButton);

    const nameInput = screen.getByRole('textbox', { name: '氏名' });
    const emailInput = screen.getByRole('textbox', { name: 'メールアドレス' });
    const textArea = screen.getByRole('textbox', { name: '本文' });

    expect(nameInput).toBeRequired();
    expect(emailInput).toBeRequired();
    expect(textArea).toBeRequired();
    expect(submitContact).not.toHaveBeenCalled();
  });

  it('値が適切に更新される', async () => {
    render(ContactForm);

    const nameInput = screen.getByRole('textbox', { name: '氏名' });
    const emailInput = screen.getByRole('textbox', { name: 'メールアドレス' });
    const textArea = screen.getByRole('textbox', { name: '本文' });

    fireEvent.input(nameInput, { target: { value: '氏名' } });
    fireEvent.input(emailInput, { target: { value: 'user@example.com' } });
    fireEvent.input(textArea, { target: { value: '本文' } });
    await tick();

    expect(nameInput).toHaveValue('氏名');
    expect(emailInput).toHaveValue('user@example.com');
    expect(textArea).toHaveValue('本文');
  });
});
