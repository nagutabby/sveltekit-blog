# sveltekit-blog

SvelteKit(`web/`)とGoバックエンド(`backend/`, Connect RPC + sqlc)で構築されたブログのリポジトリです。
このプロジェクトの成果物には、用途に応じて以下の異なるライセンスが適用されます。

## アーキテクチャ

```
                 ┌─────────────────────────┐
  外部リクエスト → │ Cloudflare Worker(router)│
                 └───────────┬─────────────┘
                     パスで振り分け
              ┌───────────────┴───────────────┐
              ▼                               ▼
   /actor*, /.well-known/*,          それ以外(記事ページ等)
   /api/articles/*
              │                               │
              ▼                               ▼
   backend (Go, Railway)  ◀── Connect RPC ──  web (SvelteKit, Railway)
   - ContactService                            - SSR/レンダリング
   - ContentService                            - marked/KaTeXでMarkdown→HTML
   - FederationAdminService(内部専用)
   - ActivityPub公開HTTPエンドポイント
              │
              ▼
           Postgres
```

- `web`と`backend`は別々のRailwayサービスとしてデプロイし、webはbackendをConnect RPC(`BACKEND_URL`)経由で呼ぶ。
- ActivityPub連携(webfinger/actor/inbox/outbox等)や記事のMarkdown/frontmatter読み込みはbackendが担い、Markdown→HTMLのレンダリングはweb側に残している(marked/KaTeX/GFM heading IDへの依存が強く移植コストに見合わないため)。
- 詳細は[`backend/README.md`](backend/README.md)、[`proto/README.md`](proto/README.md)を参照。

## License / ライセンスについて

本リポジトリ内のリソースは、ソースコードと記事コンテンツでデュアルライセンス形式を採用しています。

### 1. ソースコード (Source Code)
ブログを動かすためのプログラム、コンポーネント、スタイルシート等のソースコードはMITライセンスのもとで公開されています。
詳細な許諾条件は[LICENSE-MIT](LICENSE-MIT)をご確認ください。

### 2. 記事・コンテンツ (Articles & Contents)
`backend/content/`(記事・レビュー本文)および`web/static/content/`(画像等の静的アセット)配下のコンテンツはクリエイティブ・コモンズ 表示 4.0 国際ライセンス（CC BY 4.0）のもとで公開されています。
詳細な許諾条件は[LICENSE-CC-BY-4.0](LICENSE-CC-BY-4.0)をご確認ください。
