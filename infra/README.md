# infra

AWS CDK(TypeScript)によるインフラ定義です。RailwayからAWS(Lambda + S3 + CloudFront + DynamoDB)への移行(`migrate/11`〜`migrate/19`)で追加した。

## スタック構成

すべて`us-east-1`固定(CloudFrontにアタッチするACM証明書がus-east-1必須のため、リージョン分割を避けて統一している)。

- `DynamoDbStack`: `Follower`/`RelayConnection`テーブル(`migrate/12`)。
- `ApiStack`: Goバックエンドを動かすLambda + Function URL(`AWS_IAM`認証、CloudFrontのOrigin Access Control経由以外からは呼べない)、Secrets Manager(`migrate/13`)。
- `SiteStack`: `blog.nagutabby.uk`向けS3(SSGビルド成果物)+ CloudFront。Federation関連パス(`/actor`, `/actor/*`, `/.well-known/*`, `/api/articles/*`)はビヘイビアで同じLambda Function URLへ振り分ける(`migrate/15`)。Actor識別子(`SITE_BASE_URL`)がこのドメインに紐づいているため、ホスト名は変更できない。
- `ApiDistributionStack`: `api.nagutabby.uk`向けCloudFront。同じLambda Function URLを唯一のオリジンとする単純なパススルー(`migrate/16`)。現時点ではContactService(お問い合わせ)専用。CORSは(CloudFrontではなく)backend側`internal/server/cors.go`が`CORS_ALLOWED_ORIGIN`(=`https://blog.nagutabby.uk`)を許可することで実現している。
- `DnsStack`: `blog.nagutabby.uk`/`api.nagutabby.uk`用のRoute 53ホストゾーンとACM証明書(`migrate/17`)。ゾーンはすべて新規作成であり、既存リソースの`fromLookup`は使わない(AWS認証情報なしでも`cdk synth`が通るようにするため)。

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

デプロイ/差分確認。お問い合わせフォームの送信元/BCC先アドレスはリポジトリに書かず、`-c`でcontextとして渡す(`fromAddress`/`bccAddress`を省略した場合`cdk synth`は空文字列のまま通るが、`cdk deploy`前には必ず指定する)。

```sh
pnpm run diff -- -c fromAddress="no-reply@example.com" -c bccAddress="admin@example.com"
pnpm run deploy -- -c fromAddress="no-reply@example.com" -c bccAddress="admin@example.com"
```

`ApiStack`が作る`ActorKeysSecret`/`EmailApiTokenSecret`は空のSecrets Managerシークレットとして作成されるだけで、実際の値は入らない。デプロイ後に一度だけ手動で投入する。

```sh
aws secretsmanager put-secret-value --secret-id <ActorKeysSecretのARN> \
  --secret-string '{"ACTOR_PUBLIC_KEY_PEM":"...","ACTOR_PRIVATE_KEY_PEM":"..."}'
aws secretsmanager put-secret-value --secret-id <EmailApiTokenSecretのARN> \
  --secret-string "<Mailtrapのトークン>"
```

`SiteStack`をデプロイする前に`web/`で`pnpm build`を実行し、`web/build`に実際のSSGビルド成果物を用意しておくこと(`BucketDeployment`がこのディレクトリをそのままS3にアップロードする)。api.nagutabby.uk稼働後は、そのビルドを`PUBLIC_API_BASE_URL=https://api.nagutabby.uk`を指定して実行する(`web/src/lib/client/backendClient.ts`がビルド時にこの値を埋め込む)。

## 環境変数移行(Railway → AWS)

各サービスの環境変数をAWSに移行する際の対応表。実際にRailway上で設定されている値は`railway variables --service <name>`で確認する(コード参照だけでは棚卸し漏れの可能性があるため)。

| 変数名 | 現在 | 移行先 |
|---|---|---|
| `ACTOR_PRIVATE_KEY_PEM` / `ACTOR_PUBLIC_KEY_PEM` | backend(機密) | `ActorKeysSecret`(Secrets Manager)にJSONで格納。LambdaにはARNのみ`ACTOR_KEYS_SECRET_ARN`として渡す(`backend/internal/protectedconfig`が解決する) |
| `EMAIL_API_TOKEN` | backend(機密) | `EmailApiTokenSecret`(Secrets Manager)。Lambdaには`EMAIL_API_TOKEN_SECRET_ARN`のみ渡す |
| `FROM_ADDRESS` / `BCC_ADDRESS` | backend(非機密だが実際のアドレスは非公開) | `cdk deploy -c fromAddress=... -c bccAddress=...`経由でLambda環境変数に設定 |
| `SITE_BASE_URL` / `WEB_BASE_URL` / `CONTENT_DIR` | backend(非機密) | Lambda環境変数(`ApiStack`)。`SITE_BASE_URL`/`WEB_BASE_URL`の値は`https://blog.nagutabby.uk`のまま変更しない |
| (新規)`CORS_ALLOWED_ORIGIN` | — | Lambda環境変数、`https://blog.nagutabby.uk`固定。ContactServiceがapi.nagutabby.ukから別オリジンで呼ばれるようになったため追加 |
| `DATABASE_URL` / `BACKEND_ADDR` / `PORT` / `POSTGRES_PORT` | backend | 廃止(DynamoDB化・Lambda化に伴い不要) |
| `BACKEND_URL` | web(SSGビルド時、`$lib/server/content.ts`がbackendのContentServiceを叩く) | ビルド実行環境からLambda(またはRailway backend)へ到達可能なURLを指定。api.nagutabby.uk経由でも直接Function URLでもよい |
| (新規)`PUBLIC_API_BASE_URL` | — | webのビルド時、`$lib/client/backendClient.ts`(ブラウザから呼ぶContactService)の宛先。`https://api.nagutabby.uk`を指定する。未設定時は相対パス(同一オリジン)のまま |
| `WEB_ORIGIN` / `BACKEND_ORIGIN` | Cloudflare Worker | Worker撤去(`migrate/19`)に伴い削除 |
