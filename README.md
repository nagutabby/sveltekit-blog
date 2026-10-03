# astro-svelte-blog

Astro + Svelte (`web/`)で構築したブログです。フロントエンドは静的出力し、HonoのCloudflare WorkerがActivityPub API、お問い合わせAPI、記事公開通知APIを提供します。

## アーキテクチャ

```
外部リクエスト
     │
     ▼
Cloudflare Worker (Hono + 静的アセット)
     ├─ /actor*, /.well-known/*, /nodeinfo/*, /api/articles/*  ActivityPub
     ├─ /rpc/contact/submit                                    お問い合わせ
     ├─ /rpc/federation-admin/publish-article-activity         記事公開通知
     └─ その他のページ・画像                                   Astro静的出力
     │
     ▼
Cloudflare D1
```

- Workerのエントリーポイントは`web/src/worker/index.ts`、設定は[`backend/wrangler.jsonc`](backend/wrangler.jsonc)です。静的アセットは`web/dist/`から同じWorker経由で配信します。
- 記事・書評のMarkdownは`backend/content/`にあり、AstroのビルドとWorker用の記事メタデータ生成で読み込みます。MarkdownのHTMLレンダリングは引き続きAstro側で行います。
- D1のスキーマとマイグレーションは`backend/db/migrations/`にあります。ローカル開発と本番で同じD1バインディングを使います。
- ローカル開発、API、D1マイグレーションの手順は[`backend/README.md`](backend/README.md)を参照してください。

## ライセンス

プログラム、コンポーネント、スタイルシート等のソースコードはMITライセンスです。詳細は[LICENSE-MIT](LICENSE-MIT)を確認してください。

`backend/content/`の記事・書評、および`web/static/content/`の画像等の静的アセットはCC BY 4.0です。詳細は[LICENSE-CC-BY-4.0](LICENSE-CC-BY-4.0)を確認してください。
