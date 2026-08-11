# proto

`web`(SvelteKit)と`backend`(Go)間のConnect RPC通信に使うProtocol Buffers定義を配置するディレクトリです。

`buf`を使い、`backend/gen`（Go stubs）と`web/src/lib/gen`（TypeScript stubs）の両方を生成します。設定はPR2で追加します。
