# backend

Goバックエンド（Connect RPC + sqlc）です。

## 構成

- `cmd/server`: エントリーポイント。`BACKEND_ADDR`(既定`:8080`)でHTTPサーバを起動する。
- `internal/server`: トップレベルのHTTPハンドラ。`/healthz`と各Connect RPCサービスをマウントする。
- `internal/health`: `blog.health.v1.HealthService`の実装。web↔backend間のConnect RPC配線を検証するための最小サービス。
- `gen/`: `proto/`から`buf generate`で生成したコード(コミット対象)。
- `db/migrations`, `db/queries`, `sqlc.yaml`: [sqlc](https://sqlc.dev/)によるDBアクセス層。スキーマ/クエリはPR3で追加する。

## 開発

```sh
go build ./...
go vet ./...
go test ./...
```

Protoからのコード生成はリポジトリルートの`make generate`を使う（`web`のnode_modulesにある`protoc-gen-es`をPATHに載せて`buf generate`を実行する）。

ローカルDBは`docker compose up -d`（ルートの`docker-compose.yaml`）で起動する。ポート5432が使用中の場合は`POSTGRES_PORT`で上書きできる。
