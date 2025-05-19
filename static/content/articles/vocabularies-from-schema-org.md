---
title: "Schema.orgが提供する語彙"
image: "https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/c90f23a3ec4c47738de4752c6d5471d2/Microsoft-Fluentui-Emoji-3d-Robot-3d.1024.png"
publishedAt: 2025-05-03
updatedAt: 2025-05-03
---

<h1 id="h8d027c8ed3">はじめに</h1><p>人間だけでなくコンピューターも理解できるWebサイトを構築することは、Webの相互運用性を高め、Webのビジョンの1つであるセマンティックWebを実現する上で非常に重要です。最近ではLLM（Large Language Model）の利用率が急増し、セマンティックWebのニーズがますます高まっています。</p><p>そこで今回は、<a href="https://schema.org">Schema.org</a>というコミュニティが開発している語彙（vocabulary）を学びます。語彙を使用することでWebサイトのコンテンツが整理され、コンピューターがコンテンツをより正確に解析できるようになります。</p><h1 id="hd314f24766">Schema.orgって何？</h1><p>Schema.orgはインターネット上の構造化データのスキーマを作成し、普及させることを目的に活動しているコミュニティです。Google、Microsoft、Yahoo、Yandexなどの企業とWebコミュニティの人々で構成されています。Schema.orgの参加者は、<a href="https://github.com/schemaorg">GitHub</a>やメーリングリストで語彙に関する提案や議論を行っています。</p><h1 id="hfffca13dec">基本用語</h1><h2 id="h28c1bd2d68">語彙</h2><p>Schema.orgにおける語彙とは、スキーマの構成要素を指します。あらゆるスキーマを作成するための基本となるプロパティやデータ構造は語彙の一部です。</p><h2 id="h7ae0bad91d">スキーマ</h2><p>スキーマとは、データを記述するための仕組みです。Schema.orgは、開発者が様々な概念を記述できるようにするために、複数のスキーマを提供しています。</p><h2 id="h4e9efdf585">構造化データ</h2><p>構造化データとは、スキーマに従って記述されたデータです。Schema.orgの語彙に基づいたスキーマを使用して構造化データを作成する場合は、複数の形式で記述できます。</p><h1 id="h4e70bfec67">構造化データの具体例</h1><p>構造化データは語彙とスキーマを包含しているため、具体的な構造化データを用いることで用語の意味が理解しやすくなります。ここでは、<a href="https://schema.org/Article">Articleスキーマ</a>を参考にしながら用語を整理します。</p><h2 id="h3505fe394a">Articleスキーマに基づいた構造化データ</h2><pre><code class="language-json">{
  &quot;@context&quot;: &quot;https://schema.org&quot;,
  &quot;@type&quot;: &quot;Article&quot;,
  &quot;headline&quot;: &quot;SEOの基礎知識&quot;,
  &quot;datePublished&quot;: &quot;2023-10-05T09:00:00+09:00&quot;,
  &quot;author&quot;: {
    &quot;@type&quot;: &quot;Person&quot;,
    &quot;name&quot;: &quot;佐藤花子&quot;
  },
  &quot;publisher&quot;: {
    &quot;@type&quot;: &quot;Organization&quot;,
    &quot;name&quot;: &quot;WebTech Inc.&quot;,
    &quot;logo&quot;: {
      &quot;@type&quot;: &quot;ImageObject&quot;,
      &quot;url&quot;: &quot;https://webtechinc.example.com/logo.png&quot;,
      &quot;width&quot;: &quot;600&quot;,
      &quot;height&quot;: &quot;60&quot;
    }
  },
  &quot;description&quot;: &quot;SEOの基本概念と実践方法について解説します&quot;,
  &quot;image&quot;: {
    &quot;@type&quot;: &quot;ImageObject&quot;,
    &quot;url&quot;: &quot;https://webtechinc.example.com/images/seo-basics.jpg&quot;,
    &quot;width&quot;: &quot;1200&quot;,
    &quot;height&quot;: &quot;800&quot;
  },
  &quot;mainEntityOfPage&quot;: {
    &quot;@type&quot;: &quot;WebPage&quot;,
    &quot;@id&quot;: &quot;https://webtechinc.example.com/articles/seo-basics&quot;
  }
}</code></pre><p>この構造化データでは、Articleスキーマの以下の語彙（プロパティ）が使用されています。</p><ul><li>@context: Schema.orgの語彙を使用することを宣言するURL</li><li>@type: データの種類（この場合は「Article」）</li><li>headline: 記事のタイトル</li><li>datePublished: 記事が公開された日時</li><li>author: 記事の著者情報（Person型）</li><li>publisher: 記事の発行者情報（Organization型）</li><li>description: 記事の簡潔な説明</li><li>image: 記事に関連する画像情報（ImageObject型）</li><li>mainEntityOfPage: この構造化データが表すページのURL</li><li>logo: 発行組織のロゴ情報（ImageObject型）</li></ul><p>これらのプロパティからなる構造化データは、JSON-LD形式やMicrodata形式などの複数の形式で表現できますが、現在はマークアップ形とコンテンツを分離するために、JSON-LD形式で表現することが推奨されています。先ほどの例はJSON-LD式で記述されており、scriptタグを使用することでHTMLに埋め込めます。</p><pre><code class="language-html">&lt;!DOCTYPE html&gt;
&lt;html&gt;
  &lt;head&gt;
    &lt;title&gt;SEOの基礎知識 - WebTech Inc.&lt;/title&gt;
    &lt;script type=&quot;application/ld+json&quot;&gt;
      {
        &quot;@context&quot;: &quot;https://schema.org&quot;,
        &quot;@type&quot;: &quot;Article&quot;,
        &quot;headline&quot;: &quot;SEOの基礎知識&quot;,
        &quot;datePublished&quot;: &quot;2023-10-05T09:00:00+09:00&quot;,
        &quot;author&quot;: {
          &quot;@type&quot;: &quot;Person&quot;,
          &quot;name&quot;: &quot;佐藤花子&quot;
        },
        &quot;publisher&quot;: {
          &quot;@type&quot;: &quot;Organization&quot;,
          &quot;name&quot;: &quot;WebTech Inc.&quot;,
          &quot;logo&quot;: {
            &quot;@type&quot;: &quot;ImageObject&quot;,
            &quot;url&quot;: &quot;https://webtechinc.example.com/logo.png&quot;,
            &quot;width&quot;: &quot;600&quot;,
            &quot;height&quot;: &quot;60&quot;
          }
        },
        &quot;description&quot;: &quot;SEOの基本概念と実践方法について解説します&quot;,
        &quot;image&quot;: {
          &quot;@type&quot;: &quot;ImageObject&quot;,
          &quot;url&quot;: &quot;https://webtechinc.example.com/images/seo-basics.jpg&quot;,
          &quot;width&quot;: &quot;1200&quot;,
          &quot;height&quot;: &quot;800&quot;
        },
        &quot;mainEntityOfPage&quot;: {
          &quot;@type&quot;: &quot;WebPage&quot;,
          &quot;@id&quot;: &quot;https://webtechinc.example.com/articles/seo-basics&quot;
        }
      }
    &lt;/script&gt;
  &lt;/head&gt;
  &lt;body&gt;
    &lt;!-- ページコンテンツ --&gt;
  &lt;/body&gt;
&lt;/html&gt;</code></pre><p></p>