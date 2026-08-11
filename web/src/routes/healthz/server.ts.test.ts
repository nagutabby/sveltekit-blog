import { GET } from './+server';
import { describe, it, expect } from 'vitest';

describe('GET /healthz', () => {
  it('200とokを返す', async () => {
    const response = GET();

    expect(response.status).toBe(200);
    expect(await response.text()).toBe('ok');
  });
});
