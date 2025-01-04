import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import ContactForm from './+page.svelte';
import type { ActionData } from './$types';
import { tick } from 'svelte';
import '@testing-library/jest-dom/vitest';

describe('お問い合わせフォーム', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it('適切にレンダリングされる', () => {
    const mockForm: ActionData = {
      errors: {},
      values: {
        name: '',
        email: '',
        text: ''
      }
    };
    render(ContactForm, { props: { form: mockForm } });

    expect(screen.getByRole('textbox', { name: '氏名' })).toBeInTheDocument();
    expect(screen.getByRole('textbox', { name: 'メールアドレス' })).toBeInTheDocument();
    expect(screen.getByRole('textbox', { name: '本文' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '送信' })).toBeInTheDocument();

    const honeyPot = screen.getByRole('checkbox', { hidden: true, name: '私はロボットです' });
    expect(honeyPot.closest('.hidden')).toBeTruthy();
  });

  it('適切な値を送信したときに成功メッセージが表示される', () => {
    const mockForm: ActionData = {
      errors: {},
      values: {
        name: '氏名',
        email: 'user@example.com',
        text: '本文'
      }
    };
    render(ContactForm, { props: { form: mockForm } });

    expect(screen.getByRole('alert')).toHaveTextContent('お問い合わせを受け付けました。ご連絡までしばらくお待ちください。');
  });

  it('空文字列を送信したときにエラーメッセージが表示される', () => {
    const mockForm: ActionData = {
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
    };

    render(ContactForm, {
      props: { form: mockForm }
    });

    const alert = screen.getByRole('alert');
    expect(alert).toHaveTextContent('氏名は必須です');
    expect(alert).toHaveTextContent('メールアドレスは必須です');
    expect(alert).toHaveTextContent('本文は必須です');
  });

  it('値を送信した時に送信ボタンのテキストが変化する', async () => {
    const mockForm: ActionData = {
      errors: {},
      values: {
        name: '',
        email: '',
        text: ''
      }
    };
    render(ContactForm, { props: { form: mockForm } });

    const nameInput = screen.getByRole('textbox', { name: '氏名' });
    fireEvent.input(nameInput, { target: { value: '氏名' } });

    const emailInput = screen.getByRole('textbox', { name: 'メールアドレス' });
    fireEvent.input(emailInput, { target: { value: 'user@example.com' } });

    const textArea = screen.getByRole('textbox', { name: '本文' });
    fireEvent.input(textArea, { target: { value: '本文' } });

    const submitButton = screen.getByRole('button', { name: '送信' });
    await fireEvent.click(submitButton);
    expect(screen.getByRole('button', { name: '送信中…' })).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByRole('button', { name: '送信' })).toBeInTheDocument();
    });
  });

  it('空文字列を送信しようとしたときにrequired属性によるバリデーションが機能する', async () => {
    const mockForm: ActionData = {
      errors: {},
      values: {
        name: '',
        email: '',
        text: ''
      }
    };
    render(ContactForm, { props: { form: mockForm } });

    const submitButton = screen.getByRole('button', { name: '送信' });
    fireEvent.click(submitButton);

    const nameInput = screen.getByRole('textbox', { name: '氏名' });
    const emailInput = screen.getByRole('textbox', { name: 'メールアドレス' });
    const textArea = screen.getByRole('textbox', { name: '本文' });

    expect(nameInput).toBeRequired();
    expect(emailInput).toBeRequired();
    expect(textArea).toBeRequired();
  });

  it('値が適切に更新される', async () => {
    const mockForm: ActionData = {
      errors: {},
      values: {
        name: '',
        email: '',
        text: ''
      }
    };
    render(ContactForm, { props: { form: mockForm } });

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

