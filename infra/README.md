# infra

AWS CDK(TypeScript)によるインフラ定義です。RailwayからAWS(Lambda + S3 + CloudFront + DynamoDB)への移行(`migrate/11`〜`migrate/18`)で追加した。

## スタック構成

すべて`us-east-1`固定(CloudFrontにアタッチするACM証明書がus-east-1必須のため、リージョン分割を避けて統一している)。

- `DynamoDbStack`: `Follower`/`RelayConnection`テーブル(`migrate/12`)。
- `ApiStack`: Goバックエンドを動かすLambda + Function URL、Secrets Manager(`migrate/13`)。
- `SiteStack`: `blog.nagutabby.uk`向けS3(SSGビルド成果物)+ CloudFront。Federation関連パス(`/actor*`, `/.well-known/webfinger`等)はビヘイビアでLambda Function URLへ振り分ける(`migrate/14`)。Actor識別子(`SITE_BASE_URL`)がこのドメインに紐づいているため、ホスト名は変更できない。
- `ApiDistributionStack`: `api.nagutabby.uk`向けCloudFront。同じLambda Function URLをオリジンにし、Contact/ContentサービスなどFederation以外のAPIを配信する(`migrate/15`)。
- `DnsStack`: `blog.nagutabby.uk`/`api.nagutabby.uk`用のRoute 53ホストゾーンとACM証明書(`migrate/16`)。ゾーンはすべて新規作成であり、既存リソースの`fromLookup`は使わない(AWS認証情報なしでも`cdk synth`が通るようにするため)。

`nagutabby.uk`ゾーン自体はCloudflareに残る。`blog`と`api`サブドメインのみNSレコードでRoute 53に委任する(手順はPRの説明を参照)。

## 開発

```sh
pnpm install
pnpm run build   # tsc --noEmit
pnpm run synth   # cdk synth
```

初回デプロイ前に対象AWSアカウント/リージョンで一度だけ実行する。

```sh
pnpm exec cdk bootstrap
```

デプロイ/差分確認。

```sh
pnpm run diff
pnpm run deploy
```

## 環境変数移行(Railway → AWS)

各サービスの環境変数をAWSに移行する際の対応表。実際にRailway上で設定されている値は`railway variables --service <name>`で確認する(コード参照だけでは棚卸し漏れの可能性があるため)。

| 変数名 | 現在 | 移行先 |
|---|---|---|
| `ACTOR_PRIVATE_KEY_PEM` / `ACTOR_PUBLIC_KEY_PEM` / `EMAIL_API_TOKEN` | backend(機密) | Secrets Manager(`ApiStack`) |
| `FROM_ADDRESS` / `BCC_ADDRESS` / `SITE_BASE_URL` / `WEB_BASE_URL` / `CONTENT_DIR` | backend(非機密) | Lambda環境変数(`ApiStack`)。`SITE_BASE_URL`/`WEB_BASE_URL`の値は`https://blog.nagutabby.uk`のまま変更しない |
| `DATABASE_URL` / `BACKEND_ADDR` / `PORT` / `POSTGRES_PORT` | backend | 廃止(DynamoDB化・Lambda化に伴い不要) |
| `BACKEND_URL` | web(SSGビルド時) | ビルドパイプラインの参照先を`api.nagutabby.uk`に更新 |
| `WEB_ORIGIN` / `BACKEND_ORIGIN` | Cloudflare Worker | Worker撤去(`migrate/18`)に伴い削除 |
