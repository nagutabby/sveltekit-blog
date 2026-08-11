# backend

Goバックエンド（Connect RPC）です。

## 構成

- `cmd/server`: エントリーポイント。`BACKEND_ADDR`(既定`:8080`)でHTTPサーバを起動する。
- `cmd/create-tables`: ローカルの`dynamodb-local`にFollower/RelayConnectionテーブルを作成するワンショットコマンド(本番のテーブル定義は`infra/lib/dynamodb-stack.ts`のCDKで管理する)。
- `internal/server`: トップレベルのHTTPハンドラ。`/healthz`と各Connect RPCサービスをマウントする。
- `internal/health`: `blog.health.v1.HealthService`の実装。web↔backend間のConnect RPC配線を検証するための最小サービス。
- `internal/contact`: `blog.contact.v1.ContactService`の実装。お問い合わせフォームの入力検証とMailtrap呼び出し。
- `internal/content`: `blog.content.v1.ContentService`の実装。記事/レビューのMarkdown+frontmatterを読み込み、raw bodyのまま返す(HTML変換はwebのmarked/KaTeXパイプラインに残す)。`CONTENT_DIR`(既定`content`)を読む。スライドはfrontmatterを持たない静的PDFのため対象外(webに残す)。
- `internal/federation`: ActivityPub連携の公開HTTPエンドポイント(`webfinger`, `actor`, `actor/followers`, `actor/following`, `actor/inbox`, `actor/outbox`, `api/articles/{name}`)。外部のMastodon/リレーサーバーはJSON-LD/HTTPしか話せないため、Connect RPCではなくプレーンな`net/http`ハンドラとして実装する。HTTP Signature(`signRequest.ts`相当)の署名・検証、Follower/RelayConnectionのDB更新(DynamoDB経由)、リモートactorのfetch、署名済みAcceptの送達を行う。
- `internal/federationadmin`: `blog.federationadmin.v1.FederationAdminService`の実装(内部専用のConnect RPC)。記事のCreate/Update/Delete Activityを組み立て、LD-Signature(`signActivity`相当、RFC 8785のJSON Canonicalizationで署名)を付けて、DB上の全リレーへHTTP Signature付きで配送する。公開のActivityPub HTTPエンドポイントではないため`internal/federation`とは別パッケージ。**公開トリガーの仕組みはこのリポジトリの外にある想定**(元のSvelteKit実装の`/api/activitypub/sender`も認証なしの手動/外部トリガー呼び出しだったため、そのまま踏襲している)。呼び出し例は本ファイル末尾を参照。
- `content/`: 記事/レビューのMarkdownソース(コミット対象)。画像等の静的アセットは`web/static/content/**/images`に残る。
- `gen/`: `proto/`から`buf generate`で生成したコード(コミット対象)。
- `internal/db`: `Follower`/`RelayConnection`のDynamoDBアクセス層。[aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2)を直接使い、ORM/コード生成は使わない。主キーは`ActorId`(Postgres時代のSERIAL `id`は廃止し、実質的なユニークキーだった`ActorId`をそのまま主キーにした)。
- `internal/db/integration_test.go`: [testcontainers-go](https://golang.testcontainers.org/)で使い捨ての`dynamodb-local`を起動し、テーブル作成からDynamoDBクエリまでを検証する統合テスト。
- `internal/content/realcontent_test.go`: `content/`配下の実データを実際に読み込み、frontmatterが壊れていないかを検証する回帰テスト。
- `internal/server/federation_integration_test.go`: testcontainers-goの実`dynamodb-local` + 実HTTPサーバー + 偽のリモートMastodon actorを使い、Followの受信→DB反映→署名付きAcceptの送達までを検証する統合テスト。

## 環境変数

`SITE_BASE_URL`(既定`https://blog.nagutabby.uk`。ActivityPub Actor識別子の起点になるため値を変更しない), `WEB_BASE_URL`(`/actor/outbox`が`/atom.xml`を取得するweb側のURL。未設定時は`SITE_BASE_URL`と同じ), `ACTOR_PUBLIC_KEY_PEM`/`ACTOR_PRIVATE_KEY_PEM`(改行は`\n`エスケープ可), `FOLLOWER_TABLE_NAME`(既定`Follower`), `RELAY_CONNECTION_TABLE_NAME`(既定`RelayConnection`), `DYNAMODB_ENDPOINT`(`dynamodb-local`向けのエンドポイント上書き。本番では未設定のままにする)。AWS認証情報/リージョンはSDKのデフォルト解決(環境変数・IAMロール等、[aws-sdk-go-v2/config](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/config)に従う)。詳細は`.env.example`を参照。

## 開発

```sh
go build ./...
go vet ./...
go test ./...
```

Protoからのコード生成はリポジトリルートの`make generate`を使う（`web`のnode_modulesにある`protoc-gen-es`をPATHに載せて`buf generate`を実行する）。

ローカルのDynamoDBは`docker compose up -d`（ルートの`docker-compose.yaml`、`amazon/dynamodb-local`）で起動する。ポート8000が使用中の場合は`DYNAMODB_PORT`で上書きできる。起動後、`make db-init`（`backend/cmd/create-tables`）でFollower/RelayConnectionテーブルを作成する。

### Rancher Desktop / Colima など非Docker-Desktop環境での`go test`

`internal/db`の統合テストはtestcontainers-goでDockerコンテナを起動する。Rancher Desktopなど、Docker
socketがVM内にあり`docker context`経由でしか見えない環境では、`DOCKER_HOST`を明示し、testcontainersの
reaper(ryuk)サイドカーが socket を bind mount できずに失敗するため無効化する必要がある。

```sh
DOCKER_HOST="unix:///Users/$(whoami)/.rd/docker.sock" TESTCONTAINERS_RYUK_DISABLED=true go test ./...
```

GitHub Actions(ubuntu-latest)の標準Docker環境では追加設定は不要。

## 記事公開時にFederation通知を送る

`FederationAdminService.PublishArticleActivity`はConnectのJSON+HTTPフォールバックでも呼べるため、`curl`で直接叩ける。

```sh
curl -X POST http://localhost:8080/blog.federationadmin.v1.FederationAdminService/PublishArticleActivity \
  -H "Content-Type: application/json" \
  -d '{"articleId": "goodbye-microcms", "changeType": "CHANGE_TYPE_CREATE"}'
```

`changeType`は`CHANGE_TYPE_CREATE` / `CHANGE_TYPE_UPDATE` / `CHANGE_TYPE_DELETE`。

## デプロイ

`web`と`backend`は別々のRailwayサービスとしてデプロイする(サービス間の通信は各サービスに設定する環境変数のURLで行う。Railwayのプライベートネットワーキングが使えるならそちらでもよい)。ビルド・デプロイ設定はRailwayの[Config as Code](https://docs.railway.com/guides/config-as-code)機能で`backend/railway.json` / `web/railway.json`にコミットしている(`Dockerfile`をビルダーに指定し、`/healthz`へのヘルスチェックと再起動ポリシーを定義)。

**この構成はAWS(Lambda + S3 + CloudFront + DynamoDB, `infra/`のCDK管理)への移行中の中間状態であり、最終形ではない。** `migrate/12`時点ではストレージだけをDynamoDBに切り替え、コンピュートはRailwayに残したままにしている(変更を一度に1軸だけ進めることで問題発生時の切り戻しを容易にするため)。詳細は`infra/README.md`を参照。

同じリポジトリを指す2つのRailwayサービスをダッシュボードで作成し、それぞれ以下を設定する。

- `backend`サービス: Settings > Root Directory を`backend`に設定する(Config-as-code file pathはリポジトリルート基準で`backend/railway.json`)。環境変数は`.env.example`を参照(`SITE_BASE_URL`, `WEB_BASE_URL`, `ACTOR_PUBLIC_KEY_PEM`/`ACTOR_PRIVATE_KEY_PEM`, `EMAIL_API_TOKEN`, `FOLLOWER_TABLE_NAME`, `RELAY_CONNECTION_TABLE_NAME`, AWS認証情報等)。`PORT`はRailwayが自動で注入し、`cmd/server/main.go`がそれを`BACKEND_ADDR`未設定時のリスンポートとして使う(明示的に固定したい場合のみ`BACKEND_ADDR`を設定する)。**このサービスにデプロイする前に、対象AWSアカウントで`FollowerTable`/`RelayConnectionTable`が作成済み(`infra`で`cdk deploy`済み)かつ既存データが移行済みであることを確認すること。**
- `web`サービス: Settings > Root Directory を`web`に設定する(Config-as-code file pathは`web/railway.json`)。`BACKEND_URL`をbackendサービスの公開URL(またはRailwayプライベートネットワーキングのURL)に設定する。
- 公開ドメイン(`blog.nagutabby.uk`)の手前にCloudflare Workerを置き、パスに応じて2つのRailwayサービスへ振り分ける(`web/src/lib/workers/router.ts`、詳細は`web/wrangler.toml`の`WEB_ORIGIN`/`BACKEND_ORIGIN`)。`BACKEND_URL`(web用)と`BACKEND_ORIGIN`(Worker用)は同じbackendサービスのURLを指す。
