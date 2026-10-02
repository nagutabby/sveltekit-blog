import { describe, expect, it } from 'vitest';
import { GET } from '../../pages/healthz';

describe('GET /healthz', () => {
  it('200とokを返す', async () => {
    const response = await GET({} as never);
    expect(response.status).toBe(200);
    expect(await response.text()).toBe('ok');
  });
});
