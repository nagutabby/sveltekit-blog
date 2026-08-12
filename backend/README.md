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
- `db/migrations`: [goose](https://github.com/pressly/goose)のSQLマイグレーション(SQLite方言。本番はCloudflare D1)。`db/migrations.go`で`embed`し、goose CLIとGoテストの両方から使う。
- `db/queries`, `sqlc.yaml`, `internal/db`: [sqlc](https://sqlc.dev/)によるDBアクセス層(`engine: sqlite`, `internal/db`は生成コード)。`db/queries`内の`--`コメントはASCII文字のみにすること(sqlcのSQLiteパーサーがマルチバイト文字を含むコメントで列位置を誤認識してパースエラーになる既知の問題があるため)。
- `internal/db/d1`: 本番でCloudflare D1のHTTP query APIを叩く`db.Querier`実装。ローカル開発・テストではこれを使わず、`internal/db`が生成する`database/sql`実装に`modernc.org/sqlite`(pure Go、cgo/Docker不要)で直接繋ぐ。
- `internal/db/integration_test.go`: 一時ファイルの実SQLiteに対してgoose migrateしてからsqlcクエリを検証する統合テスト(Docker不要)。
- `internal/content/realcontent_test.go`: `content/`配下の実データを実際に読み込み、frontmatterが壊れていないかを検証する回帰テスト。
- `internal/server/federation_integration_test.go`: 一時ファイルの実SQLite + 実HTTPサーバー + 偽のリモートMastodon actorを使い、Followの受信→DB反映→署名付きAcceptの送達までを検証する統合テスト(Docker不要)。

## 環境変数

`SITE_BASE_URL`(既定`https://blog.nagutabby.uk`), `WEB_BASE_URL`(`/actor/outbox`が`/atom.xml`を取得するweb側のURL。未設定時は`SITE_BASE_URL`と同じ), `ACTOR_PUBLIC_KEY_PEM`/`ACTOR_PRIVATE_KEY_PEM`(改行は`\n`エスケープ可)。詳細は`.env.example`を参照。

DBは`CLOUDFLARE_ACCOUNT_ID`/`CLOUDFLARE_D1_DATABASE_ID`/`CLOUDFLARE_D1_API_TOKEN`が3つとも設定されていればCloudflare D1(HTTP query API経由)に接続する。1つでも欠けていれば`SQLITE_PATH`(既定`backend.db`)のローカルSQLiteファイルにフォールバックする。

## 開発

```sh
go build ./...
go vet ./...
go test ./...
```

Protoからのコード生成はリポジトリルートの`make generate`を使う（`web`のnode_modulesにある`protoc-gen-es`をPATHに載せて`buf generate`を実行する）。

DBスキーマを変更したら`sqlc generate`でクエリコードを再生成する。

ローカルDBはDocker不要。`make db-migrate`(ルートの`Makefile`)で`backend.db`にgooseマイグレーションを適用してから`go run ./cmd/server`すればよい。ファイルパスを変えたい場合は`SQLITE_PATH`で上書きする。

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

同じリポジトリを指す2つのRailwayサービスをダッシュボードで作成し、それぞれ以下を設定する。

- `backend`サービス: Settings > Root Directory を`backend`に設定する(Config-as-code file pathはリポジトリルート基準で`backend/railway.json`)。環境変数は`.env.example`を参照(`DATABASE_URL`, `SITE_BASE_URL`, `WEB_BASE_URL`, `ACTOR_PUBLIC_KEY_PEM`/`ACTOR_PRIVATE_KEY_PEM`, `EMAIL_API_TOKEN`等)。`PORT`はRailwayが自動で注入し、`cmd/server/main.go`がそれを`BACKEND_ADDR`未設定時のリスンポートとして使う(明示的に固定したい場合のみ`BACKEND_ADDR`を設定する)。
- `web`サービス: Settings > Root Directory を`web`に設定する(Config-as-code file pathは`web/railway.json`)。`BACKEND_URL`をbackendサービスの公開URL(またはRailwayプライベートネットワーキングのURL)に設定する。
- 公開ドメイン(`blog.nagutabby.uk`)の手前にCloudflare Workerを置き、パスに応じて2つのRailwayサービスへ振り分ける(`web/src/lib/workers/router.ts`、詳細は`web/wrangler.toml`の`WEB_ORIGIN`/`BACKEND_ORIGIN`)。`BACKEND_URL`(web用)と`BACKEND_ORIGIN`(Worker用)は同じbackendサービスのURLを指す。

### 本番DB(Cloudflare D1)のセットアップ

`wrangler`コマンドは`backend/wrangler.jsonc`(`d1_databases`バインディングと`migrations_dir: "db/migrations"`を定義)を使うため、**このディレクトリ(`backend/`)で実行する**。このファイルはD1の作成・マイグレーション適用専用で、backend自体はVercel Functions(Go runtime)で動かすためWorkerとしてdeployすることはない。

```sh
cd backend
wrangler d1 create sveltekit-blog-db  # 出力されたdatabase_idをwrangler.jsoncに設定する
wrangler d1 migrations apply sveltekit-blog-db --remote  # db/migrationsを適用
```

既存のNeon(PostgreSQL)からの実データ移行手順は別途追加する移行スクリプトを参照。
