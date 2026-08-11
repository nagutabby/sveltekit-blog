# proto

`web`(SvelteKit)と`backend`(Go)間のConnect RPC通信に使うProtocol Buffers定義です。

`buf generate`（ルートの`make generate`）で以下を生成する。

- `backend/gen/`: Go stubs (`protoc-gen-go` + `protoc-gen-connect-go`)
- `web/src/lib/gen/`: TypeScript stubs (`protoc-gen-es` v2, Connect-ES v2向け)

生成物はコミット対象。`blog/health/v1/health.proto`はweb↔backend間のConnect RPC配線を検証するための最小サービス。
ContactService/ContentService/FederationAdminServiceは後続PRで追加する。
