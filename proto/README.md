# proto

`web`(SvelteKit)と`backend`(Go)間のConnect RPC通信に使うProtocol Buffers定義です。

`buf generate`（ルートの`make generate`）で以下を生成する。

- `backend/gen/`: Go stubs (`protoc-gen-go` + `protoc-gen-connect-go`)
- `web/src/lib/gen/`: TypeScript stubs (`protoc-gen-es` v2, Connect-ES v2向け)

生成物はコミット対象。

- `blog/health/v1`: web↔backend間のConnect RPC配線を検証するための最小サービス
- `blog/contact/v1`: お問い合わせフォームの送信
- `blog/content/v1`: 記事/レビューのMarkdown+frontmatter取得
- `blog/federationadmin/v1`: 記事公開時のActivityPub通知(内部専用、公開HTTPエンドポイントではない)

ActivityPubの公開HTTPエンドポイント(webfinger/actor/inbox/outbox等)は外部のMastodon/リレーサーバーがJSON-LD/HTTPしか話せないため、Connect RPCではなく`backend/internal/federation`のプレーンな`net/http`ハンドラとして実装している(ここには対応するprotoはない)。
