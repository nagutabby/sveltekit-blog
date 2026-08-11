# backend

Goバックエンド（Connect RPC + sqlc）です。

## 構成

- `cmd/server`: エントリーポイント。`BACKEND_ADDR`(既定`:8080`)でHTTPサーバを起動する。
- `internal/server`: トップレベルのHTTPハンドラ。`/healthz`と各Connect RPCサービスをマウントする。
- `internal/health`: `blog.health.v1.HealthService`の実装。web↔backend間のConnect RPC配線を検証するための最小サービス。
- `gen/`: `proto/`から`buf generate`で生成したコード(コミット対象)。
- `db/migrations`: [goose](https://github.com/pressly/goose)のSQLマイグレーション。`db/migrations.go`で`embed`し、goose CLIとGoテストの両方から使う。
- `db/queries`, `sqlc.yaml`, `internal/db`: [sqlc](https://sqlc.dev/)によるDBアクセス層(`internal/db`は生成コード)。
- `internal/db/integration_test.go`: [testcontainers-go](https://golang.testcontainers.org/)で使い捨てのPostgresを起動し、goose migrateしてからsqlcクエリを検証する統合テスト。

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
