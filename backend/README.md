# backend

Goバックエンド（Connect RPC + sqlc）を配置するディレクトリです。

移行計画の詳細は各stacked PRの説明を参照してください。

- PR2でGoモジュールの骨組み（`go.mod`, `cmd/server`, `sqlc`, `goose`）を追加します。
- PR3以降でActivityPub連携、コンテンツ配信、お問い合わせ送信のロジックを段階的に移行します。
