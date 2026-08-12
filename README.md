# sveltekit-blog

SvelteKit(`web/`)とGoバックエンド(`backend/`, Connect RPC + sqlc)で構築されたブログのリポジトリです。
このプロジェクトの成果物には、用途に応じて以下の異なるライセンスが適用されます。

## アーキテクチャ

```
                      ┌───────────────────────────┐
       外部リクエスト → │ Vercel (vercel.json rewrites)│
                      └─────────────┬─────────────┘
                             パスで振り分け
              ┌───────────────────┴───────────────────┐
              ▼                                       ▼
   /actor*, /.well-known/webfinger,             それ以外(記事ページ等)
   /api/articles/*, /blog.contact.v1.ContactService/*,
   /blog.federationadmin.v1.FederationAdminService/*
              │                                       │
              ▼                                       ▼
   backend (Go, Vercel Functions)              web (SvelteKit adapter-static, Vercel)
   - ContactService                             - SSGで全ページ静的化
   - FederationAdminService(共有シークレット保護) - marked/KaTeXでMarkdown→HTML
   - ActivityPub公開HTTPエンドポイント            - ビルド時にbackend/content/を直接読む
              │
              ▼
       Cloudflare D1 (SQLite)
```

- `web`(SvelteKit adapter-static)と`backend`(Go)は1つのVercelプロジェクト内の別サービス(`services`)としてデプロイし、ルートの`vercel.json`のrewritesでパスに応じて振り分ける。webはビルド時に`backend/content/`のMarkdownを直接読み込むため、backendへのネットワーク呼び出しは発生しない。
- ActivityPub連携(webfinger/actor/inbox/outbox等)はbackendが担い、Markdown→HTMLのレンダリングはweb側に残している(marked/KaTeX/GFM heading IDへの依存が強く移植コストに見合わないため)。
- DBはCloudflare D1(SQLite互換)。backendからはD1のHTTP query API経由でアクセスする(詳細は[`backend/README.md`](backend/README.md))。
- 詳細は[`backend/README.md`](backend/README.md)、[`proto/README.md`](proto/README.md)を参照。

## License / ライセンスについて

本リポジトリ内のリソースは、ソースコードと記事コンテンツでデュアルライセンス形式を採用しています。

### 1. ソースコード (Source Code)
ブログを動かすためのプログラム、コンポーネント、スタイルシート等のソースコードはMITライセンスのもとで公開されています。
詳細な許諾条件は[LICENSE-MIT](LICENSE-MIT)をご確認ください。

### 2. 記事・コンテンツ (Articles & Contents)
`backend/content/`(記事・レビュー本文)および`web/static/content/`(画像等の静的アセット)配下のコンテンツはクリエイティブ・コモンズ 表示 4.0 国際ライセンス（CC BY 4.0）のもとで公開されています。
詳細な許諾条件は[LICENSE-CC-BY-4.0](LICENSE-CC-BY-4.0)をご確認ください。
