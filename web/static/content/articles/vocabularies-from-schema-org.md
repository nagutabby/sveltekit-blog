---
title: Schema.orgが提供する語彙
image: images/Microsoft-Fluentui-Emoji-3d-Robot-3d.1024.png
publishedAt: 2025-05-03
updatedAt: 2025-05-03
---
# はじめに

人間だけでなくコンピューターも理解できるWebサイトを構築することは、Webの相互運用性を高め、Webのビジョンの1つであるセマンティックWebを実現する上で非常に重要です。最近ではLLM（Large Language Model）の利用率が急増し、セマンティックWebのニーズがますます高まっています。

そこで今回は、[Schema.org](https://schema.org)というコミュニティが開発している語彙（vocabulary）を学びます。語彙を使用することでWebサイトのコンテンツが整理され、コンピューターがコンテンツをより正確に解析できるようになります。

# Schema.orgって何？

Schema.orgはインターネット上の構造化データのスキーマを作成し、普及させることを目的に活動しているコミュニティです。Google、Microsoft、Yahoo、Yandexなどの企業とWebコミュニティの人々で構成されています。Schema.orgの参加者は、[GitHub](https://github.com/schemaorg)やメーリングリストで語彙に関する提案や議論を行っています。

# 基本用語

## 語彙

Schema.orgにおける語彙とは、スキーマの構成要素を指します。あらゆるスキーマを作成するための基本となるプロパティやデータ構造は語彙の一部です。

## スキーマ

スキーマとは、データを記述するための仕組みです。Schema.orgは、開発者が様々な概念を記述できるようにするために、複数のスキーマを提供しています。

## 構造化データ

構造化データとは、スキーマに従って記述されたデータです。Schema.orgの語彙に基づいたスキーマを使用して構造化データを作成する場合は、複数の形式で記述できます。

# 構造化データの具体例

構造化データは語彙とスキーマを包含しているため、具体的な構造化データを用いることで用語の意味が理解しやすくなります。ここでは、[Articleスキーマ](https://schema.org/Article)を参考にしながら用語を整理します。

## Articleスキーマに基づいた構造化データ

```json
{
  "@context": "https://schema.org",
  "@type": "Article",
  "headline": "SEOの基礎知識",
  "datePublished": "2023-10-05T09:00:00+09:00",
  "author": {
    "@type": "Person",
    "name": "佐藤花子"
  },
  "publisher": {
    "@type": "Organization",
    "name": "WebTech Inc.",
    "logo": {
      "@type": "ImageObject",
      "url": "https://webtechinc.example.com/logo.png",
      "width": "600",
      "height": "60"
    }
  },
  "description": "SEOの基本概念と実践方法について解説します",
  "image": {
    "@type": "ImageObject",
    "url": "https://webtechinc.example.com/images/seo-basics.jpg",
    "width": "1200",
    "height": "800"
  },
  "mainEntityOfPage": {
    "@type": "WebPage",
    "@id": "https://webtechinc.example.com/articles/seo-basics"
  }
}
```

この構造化データでは、Articleスキーマの以下の語彙（プロパティ）が使用されています。

-   @context: Schema.orgの語彙を使用することを宣言するURL
-   @type: データの種類（この場合は「Article」）
-   headline: 記事のタイトル
-   datePublished: 記事が公開された日時
-   author: 記事の著者情報（Person型）
-   publisher: 記事の発行者情報（Organization型）
-   description: 記事の簡潔な説明
-   image: 記事に関連する画像情報（ImageObject型）
-   mainEntityOfPage: この構造化データが表すページのURL
-   logo: 発行組織のロゴ情報（ImageObject型）

これらのプロパティからなる構造化データは、JSON-LD形式やMicrodata形式などの複数の形式で表現できますが、現在はマークアップ形とコンテンツを分離するために、JSON-LD形式で表現することが推奨されています。先ほどの例はJSON-LD式で記述されており、scriptタグを使用することでHTMLに埋め込めます。

```html
<!DOCTYPE html>
<html>
  <head>
    <title>SEOの基礎知識 - WebTech Inc.</title>
    <script type="application/ld+json">
      {
        "@context": "https://schema.org",
        "@type": "Article",
        "headline": "SEOの基礎知識",
        "datePublished": "2023-10-05T09:00:00+09:00",
        "author": {
          "@type": "Person",
          "name": "佐藤花子"
        },
        "publisher": {
          "@type": "Organization",
          "name": "WebTech Inc.",
          "logo": {
            "@type": "ImageObject",
            "url": "https://webtechinc.example.com/logo.png",
            "width": "600",
            "height": "60"
          }
        },
        "description": "SEOの基本概念と実践方法について解説します",
        "image": {
          "@type": "ImageObject",
          "url": "https://webtechinc.example.com/images/seo-basics.jpg",
          "width": "1200",
          "height": "800"
        },
        "mainEntityOfPage": {
          "@type": "WebPage",
          "@id": "https://webtechinc.example.com/articles/seo-basics"
        }
      }
    </script>
  </head>
  <body>
    <!-- ページコンテンツ -->
  </body>
</html>
```
