import { Hono } from 'hono';
import { handleActivityPubRequest } from './federation';
import { rpcRoutes } from './rpc';
import type { WorkerEnv } from './bindings';

const app = new Hono<{ Bindings: WorkerEnv }>();

app.route('/rpc', rpcRoutes);
app.get('/healthz', (context) => context.text('ok'));

app.get('/.well-known/webfinger', (context) => handleActivityPubRequest(context.req.raw, context.env, context.req.path)
  .then((response) => response ?? context.notFound()));
app.get('/.well-known/nodeinfo', (context) => handleActivityPubRequest(context.req.raw, context.env, context.req.path)
  .then((response) => response ?? context.notFound()));
app.get('/nodeinfo/:version', (context) => handleActivityPubRequest(context.req.raw, context.env, context.req.path)
  .then((response) => response ?? context.notFound()));
app.get('/actor', (context) => handleActivityPubRequest(context.req.raw, context.env, context.req.path)
  .then((response) => response ?? context.notFound()));
app.get('/actor/followers', (context) => handleActivityPubRequest(context.req.raw, context.env, context.req.path)
  .then((response) => response ?? context.notFound()));
app.get('/actor/following', (context) => handleActivityPubRequest(context.req.raw, context.env, context.req.path)
  .then((response) => response ?? context.notFound()));
app.get('/actor/outbox', (context) => handleActivityPubRequest(context.req.raw, context.env, context.req.path)
  .then((response) => response ?? context.notFound()));
app.post('/actor/inbox', (context) => handleActivityPubRequest(context.req.raw, context.env, context.req.path)
  .then((response) => response ?? context.notFound()));
app.get('/api/articles/:name', (context) => handleActivityPubRequest(context.req.raw, context.env, context.req.path)
  .then((response) => response ?? context.notFound()));

app.onError((error, context) => {
  console.error(JSON.stringify({
    message: 'unhandled Worker request error',
    path: context.req.path,
    error: error.message
  }));
  return context.json({ error: 'Internal Server Error' }, 500);
});

export type AppType = typeof app;
export default app;
