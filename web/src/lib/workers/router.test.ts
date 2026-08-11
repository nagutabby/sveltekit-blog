import { afterEach, describe, expect, it, vi } from 'vitest';
import { isBackendPath, isRelayInboxRequest, resolveOrigin } from './router';

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('isBackendPath', () => {
  it.each([
    ['/actor', true],
    ['/actor/inbox', true],
    ['/actor/followers', true],
    ['/.well-known/webfinger', true],
    ['/api/articles/my-article', true],
    ['/', false],
    ['/articles/my-article', false],
    ['/atom.xml', false],
    ['/sitemap.xml', false],
    ['/api/articlesnotreally', false]
  ])('%s -> %s', (pathname, expected) => {
    expect(isBackendPath(pathname)).toBe(expected);
  });
});

describe('resolveOrigin', () => {
  const env = { WEB_ORIGIN: 'https://web.example', BACKEND_ORIGIN: 'https://backend.example' };

  it('routes backend paths to BACKEND_ORIGIN', () => {
    expect(resolveOrigin('/actor', env)).toBe('https://backend.example');
  });

  it('routes everything else to WEB_ORIGIN', () => {
    expect(resolveOrigin('/', env)).toBe('https://web.example');
    expect(resolveOrigin('/articles/my-article', env)).toBe('https://web.example');
  });
});

describe('isRelayInboxRequest', () => {
  it('matches a relay user-agent posting to /actor/inbox', () => {
    expect(isRelayInboxRequest('/actor/inbox', 'SomeRelay/1.0')).toBe(true);
    expect(isRelayInboxRequest('/actor/inbox', 'relay-server')).toBe(true);
  });

  it('does not match a non-relay user-agent', () => {
    expect(isRelayInboxRequest('/actor/inbox', 'Mastodon/4.2.0')).toBe(false);
  });

  it('does not match a relay user-agent on other paths', () => {
    expect(isRelayInboxRequest('/actor', 'SomeRelay/1.0')).toBe(false);
  });

  it('is case-insensitive on the user-agent', () => {
    expect(isRelayInboxRequest('/actor/inbox', 'RELAY-BOT')).toBe(true);
  });
});

describe('default.fetch (worker entry point)', () => {
  const env = { WEB_ORIGIN: 'https://web.example', BACKEND_ORIGIN: 'https://backend.example' };

  const setup = async () => {
    const router = (await import('./router')).default;
    return router;
  };

  it('proxies a normal page request to WEB_ORIGIN', async () => {
    const fetchMock = vi.fn().mockResolvedValue(new Response('ok'));
    vi.stubGlobal('fetch', fetchMock);

    const router = await setup();
    const request = new Request('https://blog.nagutabby.uk/articles/my-article');
    await router.fetch(request, env);

    expect(fetchMock).toHaveBeenCalledTimes(1);
    const proxied = fetchMock.mock.calls[0][0] as Request;
    expect(proxied.url).toBe('https://web.example/articles/my-article');
  });

  it('proxies /actor to BACKEND_ORIGIN', async () => {
    const fetchMock = vi.fn().mockResolvedValue(new Response('ok'));
    vi.stubGlobal('fetch', fetchMock);

    const router = await setup();
    const request = new Request('https://blog.nagutabby.uk/actor');
    await router.fetch(request, env);

    const proxied = fetchMock.mock.calls[0][0] as Request;
    expect(proxied.url).toBe('https://backend.example/actor');
  });

  it('proxies a relay Accept to /actor/inbox instead of blocking it', async () => {
    const fetchMock = vi.fn().mockResolvedValue(new Response('ok'));
    vi.stubGlobal('fetch', fetchMock);

    const router = await setup();
    const request = new Request('https://blog.nagutabby.uk/actor/inbox', {
      method: 'POST',
      headers: { 'user-agent': 'SomeRelay/1.0' },
      body: JSON.stringify({ type: 'Accept' })
    });
    const response = await router.fetch(request, env);

    expect(fetchMock).toHaveBeenCalledTimes(1);
    expect(response.status).not.toBe(403);
  });

  it('blocks a relay posting a non-Accept activity to /actor/inbox', async () => {
    const fetchMock = vi.fn().mockResolvedValue(new Response('ok'));
    vi.stubGlobal('fetch', fetchMock);

    const router = await setup();
    const request = new Request('https://blog.nagutabby.uk/actor/inbox', {
      method: 'POST',
      headers: { 'user-agent': 'SomeRelay/1.0' },
      body: JSON.stringify({ type: 'Follow' })
    });
    const response = await router.fetch(request, env);

    expect(fetchMock).not.toHaveBeenCalled();
    expect(response.status).toBe(403);
    const body = await response.json();
    expect(body.error).toBe('Forbidden');
  });

  it('does not block a non-relay user-agent posting to /actor/inbox', async () => {
    const fetchMock = vi.fn().mockResolvedValue(new Response('ok'));
    vi.stubGlobal('fetch', fetchMock);

    const router = await setup();
    const request = new Request('https://blog.nagutabby.uk/actor/inbox', {
      method: 'POST',
      headers: { 'user-agent': 'Mastodon/4.2.0' },
      body: JSON.stringify({ type: 'Follow' })
    });
    const response = await router.fetch(request, env);

    expect(fetchMock).toHaveBeenCalledTimes(1);
    expect(response.status).not.toBe(403);
  });
});
