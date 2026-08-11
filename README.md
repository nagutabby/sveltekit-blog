# sveltekit-blog

SvelteKit(`web/`)とGoバックエンド(`backend/`, Connect RPC)で構築されたブログのリポジトリです。AWS(Lambda + S3 + CloudFront + DynamoDB)上で動作し、インフラは`infra/`のAWS CDKで管理します。
このプロジェクトの成果物には、用途に応じて以下の異なるライセンスが適用されます。

## アーキテクチャ

```
        blog.nagutabby.uk                              api.nagutabby.uk
               │                                               │
               ▼                                               ▼
         CloudFront                                       CloudFront
               │                                               │
        パスで振り分け                                          │
   ┌───────────┴────────────┐                                  │
   ▼                        ▼                                  │
/actor, /actor/*,      それ以外                                 │
/.well-known/*,        (記事ページ等)                            │
/api/articles/*             │                                  │
   │                        ▼                                  │
   │                  S3(SSGビルド成果物)                        │
   │                                                            │
   └─────────────────────────┬──────────────────────────────────┘
                              ▼
                    Lambda(Go, backend)
                    - ActivityPub公開HTTPエンドポイント
                    - ContentService / ContactService / FederationAdminService
                              │
                              ▼
                          DynamoDB
```

- `web`(SvelteKit)は全ページ`prerender = true`(SSG)でビルドし、サーバーを使わずS3+CloudFrontで配信する。
- `backend`(Go)はLambda上で動作し、DynamoDBを使う。ActivityPub連携(webfinger/actor/inbox/outbox等)や記事のMarkdown/frontmatter読み込みを担う。Markdown→HTMLのレンダリングはweb側に残している(marked/KaTeX/GFM heading IDへの依存が強く移植コストに見合わないため)。
- Actor識別子(`SITE_BASE_URL`)が`blog.nagutabby.uk`に紐づいているため、ActivityPub関連エンドポイントは同ドメインのCloudFrontからLambdaへ振り分ける。Contact等それ以外のAPIは`api.nagutabby.uk`から呼ぶ。
- 詳細は[`backend/README.md`](backend/README.md)、[`infra/README.md`](infra/README.md)、[`proto/README.md`](proto/README.md)を参照。

## License / ライセンスについて

本リポジトリ内のリソースは、ソースコードと記事コンテンツでデュアルライセンス形式を採用しています。

### 1. ソースコード (Source Code)
ブログを動かすためのプログラム、コンポーネント、スタイルシート等のソースコードはMITライセンスのもとで公開されています。
詳細な許諾条件は[LICENSE-MIT](LICENSE-MIT)をご確認ください。

### 2. 記事・コンテンツ (Articles & Contents)
`backend/content/`(記事・レビュー本文)および`web/static/content/`(画像等の静的アセット)配下のコンテンツはクリエイティブ・コモンズ 表示 4.0 国際ライセンス（CC BY 4.0）のもとで公開されています。
詳細な許諾条件は[LICENSE-CC-BY-4.0](LICENSE-CC-BY-4.0)をご確認ください。
