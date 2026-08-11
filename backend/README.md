# backend

Goバックエンド（Connect RPC + sqlc）です。

## 構成

- `cmd/server`: エントリーポイント。`BACKEND_ADDR`(既定`:8080`)でHTTPサーバを起動する。
- `internal/server`: トップレベルのHTTPハンドラ。`/healthz`と各Connect RPCサービスをマウントする。
- `internal/health`: `blog.health.v1.HealthService`の実装。web↔backend間のConnect RPC配線を検証するための最小サービス。
- `internal/contact`: `blog.contact.v1.ContactService`の実装。お問い合わせフォームの入力検証とMailtrap呼び出し。
- `internal/content`: `blog.content.v1.ContentService`の実装。記事/レビューのMarkdown+frontmatterを読み込み、raw bodyのまま返す(HTML変換はwebのmarked/KaTeXパイプラインに残す)。`CONTENT_DIR`(既定`content`)を読む。スライドはfrontmatterを持たない静的PDFのため対象外(webに残す)。
- `internal/federation`: ActivityPub連携の公開HTTPエンドポイント(`webfinger`, `actor`, `actor/followers`, `actor/following`, `actor/inbox`, `actor/outbox`, `api/articles/{name}`)。外部のMastodon/リレーサーバーはJSON-LD/HTTPしか話せないため、Connect RPCではなくプレーンな`net/http`ハンドラとして実装する。HTTP Signature(`signRequest.ts`相当)の署名・検証、Follower/RelayConnectionのDB更新(sqlc経由)、リモートactorのfetch、署名済みAcceptの送達を行う。
- `internal/federationadmin`: `blog.federationadmin.v1.FederationAdminService`の実装(内部専用のConnect RPC)。記事のCreate/Update/Delete Activityを組み立て、LD-Signature(`signActivity`相当、RFC 8785のJSON Canonicalizationで署名)を付けて、DB上の全リレーへHTTP Signature付きで配送する。公開のActivityPub HTTPエンドポイントではないため`internal/federation`とは別パッケージ。**公開トリガーの仕組みはこのリポジトリの外にある想定**(元のSvelteKit実装の`/api/activitypub/sender`も認証なしの手動/外部トリガー呼び出しだったため、そのまま踏襲している)。呼び出し例は本ファイル末尾を参照。
- `content/`: 記事/レビューのMarkdownソース(コミット対象)。画像等の静的アセットは`web/static/content/**/images`に残る。
- `gen/`: `proto/`から`buf generate`で生成したコード(コミット対象)。
- `db/migrations`: [goose](https://github.com/pressly/goose)のSQLマイグレーション。`db/migrations.go`で`embed`し、goose CLIとGoテストの両方から使う。
- `db/queries`, `sqlc.yaml`, `internal/db`: [sqlc](https://sqlc.dev/)によるDBアクセス層(`internal/db`は生成コード)。
- `internal/db/integration_test.go`: [testcontainers-go](https://golang.testcontainers.org/)で使い捨てのPostgresを起動し、goose migrateしてからsqlcクエリを検証する統合テスト。
- `internal/content/realcontent_test.go`: `content/`配下の実データを実際に読み込み、frontmatterが壊れていないかを検証する回帰テスト。
- `internal/server/federation_integration_test.go`: testcontainers-goの実Postgres + 実HTTPサーバー + 偽のリモートMastodon actorを使い、Followの受信→DB反映→署名付きAcceptの送達までを検証する統合テスト。

## 環境変数

`DATABASE_URL`, `SITE_BASE_URL`(既定`https://blog.nagutabby.uk`), `WEB_BASE_URL`(`/actor/outbox`が`/atom.xml`を取得するweb側のURL。未設定時は`SITE_BASE_URL`と同じ), `ACTOR_PUBLIC_KEY_PEM`/`ACTOR_PRIVATE_KEY_PEM`(改行は`\n`エスケープ可)。詳細は`.env.example`を参照。

## 開発

```sh
go build ./...
go vet ./...
go test ./...
```

Protoからのコード生成はリポジトリルートの`make generate`を使う（`web`のnode_modulesにある`protoc-gen-es`をPATHに載せて`buf generate`を実行する）。

DBスキーマを変更したら`sqlc generate`でクエリコードを再生成する。

ローカルDBは`docker compose up -d`（ルートの`docker-compose.yaml`）で起動する。ポート5432が使用中の場合は`POSTGRES_PORT`で上書きできる。

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
