# Backend API

Backend APIはHonoのCloudflare Worker（`web/src/worker/`）で動作します。このディレクトリには記事・書評のMarkdownと、Cloudflare D1のマイグレーションを置いています。

## API

- ActivityPub: `/.well-known/webfinger`、`/.well-known/nodeinfo`、`/nodeinfo/2.0`、`/nodeinfo/2.1`、`/actor`以下、および`/api/articles/{id}`
- お問い合わせ: `POST /rpc/contact/submit`。JSONで`name`、`email`、`text`、`imRobot`を受け取り、入力エラーは`{ "errors": { ... } }`で返します。
- 記事公開通知: `POST /rpc/federation-admin/publish-article-activity`。`Authorization: Bearer <FEDERATION_ADMIN_TOKEN>`が必要です。

記事通知のリクエスト例:

```sh
curl -X POST https://blog.nagutabby.uk/rpc/federation-admin/publish-article-activity \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $FEDERATION_ADMIN_TOKEN" \
  -d '{"articleId":"goodbye-microcms","changeType":"create"}'
```

`changeType`には`create`、`update`、`delete`を指定します。記事公開フローの外部呼び出し元はこのURLとJSON形式へ切り替え、`FEDERATION_ADMIN_TOKEN`をBearer認証で送る必要があります。旧Connect RPC URLは提供しません。

## ローカル開発

`web/.env.example`を参考に`web/.env`を用意し、フロントエンドの開発サーバーからローカルWorkerへプロキシします。Worker用Secretsは`backend/.dev.vars.example`を`backend/.dev.vars`へコピーして設定します。`.dev.vars`はGit管理対象外です。

```sh
pnpm --dir web install
pnpm --dir web run dev:worker
```

Workerは`http://localhost:8787`で起動します。ローカルD1にスキーマを適用するには、リポジトリルートで次を実行します。

```sh
make db-migrate
```

Astroの開発サーバーは別ターミナルで`pnpm --dir web run dev`を実行します。

## D1マイグレーション

SQLマイグレーションは`db/migrations/`にあります。ローカル適用は`make db-migrate`、既存の本番D1への適用は次のコマンドを使います。

```sh
pnpm --dir web exec wrangler d1 migrations apply sveltekit-blog-db --remote --config ../backend/wrangler.jsonc
```

本番マイグレーション適用はD1のスキーマを変更します。デプロイ作業前に適用対象を確認してください。

ゾーンとDNSレコードの確認にはCloudflare CLIを使います。既存レコードを確認してから、ブログ用ホスト名だけを変更してください。

```sh
cloudflare zones list --name nagutabby.uk
cloudflare dns records list --zone nagutabby.uk
```

## Secretsとデプロイ

Wrangler設定では次のSecretsを必須にしています。値を`wrangler.jsonc`やソースへ書かず、CloudflareダッシュボードまたはWrangler Secretsから登録してください。

- `ACTOR_PUBLIC_KEY_PEM`、`ACTOR_PRIVATE_KEY_PEM`
- `FEDERATION_ADMIN_TOKEN`
- `EMAIL_API_TOKEN`、`FROM_ADDRESS`、`BCC_ADDRESS`

本番への自動デプロイはGitHub Actionsから行います。リポジトリに`CLOUDFLARE_API_TOKEN` Secretと`CLOUDFLARE_ACCOUNT_ID` Actions Variableを登録してください。`main`へのpushで、型確認・テスト・ビルドがすべて成功した後にWorkerと静的アセットをデプロイします。Pull Requestではデプロイしません。

ローカルから手動でデプロイする場合は`pnpm --dir web run deploy`を使います。記事通知APIのリクエスト例とBearer認証は、上記のAPI節を参照してください。
