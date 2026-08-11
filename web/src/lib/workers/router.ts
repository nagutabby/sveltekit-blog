// Cloudflare Worker in front of blog.nagutabby.uk. web (SvelteKit) and
// backend (Go) are deployed as separate Railway services with no shared
// domain, so this worker does path-based routing at the edge in addition
// to its original job of blocking relay servers from calling /actor/inbox
// directly (see isRelayInboxRequest below).
//
// WEB_ORIGIN / BACKEND_ORIGIN must be set to the actual public URLs of
// the two Railway services (wrangler.toml [vars], or `wrangler secret put`
// if you'd rather not have them in source control).

export interface Env {
  WEB_ORIGIN: string;
  BACKEND_ORIGIN: string;
}

// Paths served by backend's public ActivityPub HTTP surface
// (internal/federation), plus Connect RPC services the browser calls
// directly (e.g. the contact form). Everything else is web's.
const BACKEND_PATH_PREFIXES = [
  '/actor',
  '/.well-known/webfinger',
  '/api/articles',
  '/blog.contact.v1.ContactService'
];

export function isBackendPath(pathname: string): boolean {
  return BACKEND_PATH_PREFIXES.some(
    (prefix) => pathname === prefix || pathname.startsWith(`${prefix}/`)
  );
}

export function resolveOrigin(pathname: string, env: Env): string {
  return isBackendPath(pathname) ? env.BACKEND_ORIGIN : env.WEB_ORIGIN;
}

export function isRelayInboxRequest(pathname: string, userAgent: string): boolean {
  return pathname === '/actor/inbox' && userAgent.toLowerCase().includes('relay');
}

function forbiddenRelayResponse(pathname: string, userAgent: string): Response {
  return new Response(
    JSON.stringify({
      error: 'Forbidden',
      status: 403,
      message: 'Relay server access is not permitted',
      details: {
        reason: 'This server does not accept requests from relay servers',
        path: pathname,
        userAgent
      }
    } as const),
    {
      status: 403,
      headers: {
        'Content-Type': 'application/activity+json'
      }
    }
  );
}

// web is fully statically prerendered (SSG); we never want the browser or
// any CDN in front of this worker to cache a stale copy of a page/asset.
// backend's own responses (federation, contact RPC) are left untouched.
function withNoStore(response: Response): Response {
  const noStoreResponse = new Response(response.body, response);
  noStoreResponse.headers.set('Cache-Control', 'no-store');
  noStoreResponse.headers.set('Pragma', 'no-cache');
  return noStoreResponse;
}

function internalErrorResponse(): Response {
  return new Response(
    JSON.stringify({
      error: 'Internal Server Error',
      message: 'An error occurred while processing your request'
    } as const),
    {
      status: 500,
      headers: {
        'Content-Type': 'application/activity+json'
      }
    }
  );
}

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    try {
      const userAgent = request.headers.get('user-agent') || '';
      const url = new URL(request.url);
      const pathname = url.pathname;

      if (isRelayInboxRequest(pathname, userAgent)) {
        const activity = await request
          .clone()
          .json()
          .catch(() => null);
        if (!activity || activity.type !== 'Accept') {
          return forbiddenRelayResponse(pathname, userAgent);
        }
      }

      const targetOrigin = resolveOrigin(pathname, env);
      const targetURL = new URL(pathname + url.search, targetOrigin);
      const response = await fetch(new Request(targetURL, request));

      return isBackendPath(pathname) ? response : withNoStore(response);
    } catch {
      return internalErrorResponse();
    }
  }
};
