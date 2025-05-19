---
title: HugoにTailwind CSS 3をインストール
image: images/Microsoft-Fluentui-Emoji-3d-Leaf-Fluttering-In-Wind-3d.1024.png
publishedAt: 2023-10-09T00:00:00.000Z
updatedAt: 2024-05-03T00:00:00.000Z
---

<p>HugoにTailwind CSS 3をインストールしてみます。</p><h1 id="h3e2f255d1b">Hugoプロジェクトの作成</h1><pre><code class="language-bash">hugo new site hugo-tailwind-css-3-sample
cd hugo-tailwind-css-3-sample
git init</code></pre><h1 id="h6a003ead1f">Tailwindのインストール</h1><p>pnpmでTailwindをインストールします。TailwindはPostCSSのプラグインなので、PostCSSのインストールも必要です。</p><pre><code class="language-bash">pnpm init 
pnpm add -D autoprefixer postcss postcss-cli tailwindcss</code></pre><p><code>.gitignore</code>を作成し、そのファイルに<code>node_modules</code>を追記します。</p><p>これにより<code>node_modules</code>がGitの管理対象から除外されます。</p><pre><code>/node_modules</code></pre><p><code>tailwind.config.js</code>と<code>postcss.config.js</code>を作成します。</p><pre><code class="language-bash">pnpm tailwind init -p</code></pre><h1 id="hdbc0c6b369">Cache bustersの設定</h1><p>Hugo 0.112でTailwind CSS 3がサポートされました。Hugoは、HTMLをビルドするタイミングで<code>hugo_stats.json</code>に自動的にTailwindのクラスを追記します。そして、この<code>hugo_stats.json</code>の内容が変更されるたびに、CSSファイルをビルドし直します。</p><p><code>hugo_stats.json</code>の内容に応じてTailwindのコードを生成させるために、<code>tailwind.config.js</code>を変更します。</p><pre><code class="language-javascript">/** @type {import(&apos;tailwindcss&apos;).Config} */
module.exports = {
  content: [
    &quot;./hugo_stats.json&quot;,
  ],
  theme: {
    extend: {},
  },
  plugins: []
}</code></pre><p>hugo.tomlにCache bustersの設定を追加します。Cache bustersはその名の通り、キャッシュを削除するための設定オプションです。</p><pre><code class="language-toml">baseURL = &apos;https://example.org/&apos;
defaultContentLanguage = &apos;ja&apos;
languageCode = &apos;ja&apos;
title = &apos;タイトル&apos;

[module]
[[module.mounts]]
source = &quot;assets&quot;
target = &quot;assets&quot;
[[module.mounts]]
source = &quot;hugo_stats.json&quot;
target = &quot;assets/watching/hugo_stats.json&quot;

[build]
[build.buildStats]
enable = true
[[build.cachebusters]]
source = &quot;assets/watching/hugo_stats\\.json&quot;
target = &quot;style\\.scss&quot;</code></pre><p>Hugo 0.112から0.114までは<code>writeStats = true</code>というオプションを指定していましたが、Hugo 0.115から<code>build.buildStats</code>に<code>enable = true</code>を記述する形式に置き換えられました。<code>writeStats</code>は現在非推奨のオプションであり、将来的に廃止される予定です。</p><p><code>hugo_stats.json</code>についても<code>node_modules</code>と同様にGitの管理対象から除外しておきましょう。</p><pre><code>/node_modules
/hugo_stats.json</code></pre><p>続いて<code>assets</code>ディレクトリに<code>css</code>ディレクトリを作成し、その中に<code>style.scss</code>を作成します。</p><pre><code class="language-scss">@tailwind base;
@tailwind components;
@tailwind utilities;</code></pre><h1 id="h4a65c38694">レイアウトの作成</h1><p><code>layouts</code>ディレクトリに<code>_default</code>ディレクトリを作成し、その中に<code>baseof.html</code>を作成します。</p><pre><code class="language-html">&lt;!DOCTYPE html&gt;
&lt;html lang=&quot;{{ with .Site.LanguageCode }}{{ . }}{{ end }}&quot;&gt;
  &lt;head&gt;
    &lt;meta charset=&quot;utf-8&quot;&gt;
    &lt;meta name=&quot;viewport&quot; content=&quot;width=device-width, initial-scale=1.0&quot;&gt;
    &lt;title&gt;{{ block &quot;title&quot; . }}{{ end }}&lt;/title&gt;
    {{ $style := resources.Get &quot;css/style.scss&quot; | toCSS | resources.PostCSS }}
    {{ if hugo.IsProduction }}
      {{ $style = $style | minify | fingerprint | resources.PostProcess }}
    {{ end }}
    &lt;link href=&quot;{{ $style.RelPermalink }}&quot; rel=&quot;stylesheet&quot; /&gt;
  &lt;/head&gt;
  &lt;body&gt;
    &lt;main&gt;
      {{ block &quot;main&quot; . }}{{ end }}
    &lt;/main&gt;
  &lt;/body&gt;
&lt;/html&gt;</code></pre><p>同じく<code>layouts/_default</code>ディレクトリに<code>index.html</code>を作成します。</p><pre><code class="language-html">{{ define &quot;title&quot; }}
{{- .Site.Title -}}
{{ end }}
{{ define &quot;main&quot; }}
&lt;div class=&quot;text-5xl text-red-500&quot;&gt;
  {{ .Content }}
&lt;/div&gt;
{{ end }}</code></pre><h1 id="h09d1a5c048">コンテンツの作成</h1><p>ホームに表示するコンテンツを作成するために、マークダウンファイルを作成します。</p><pre><code class="language-bash">hugo new content/_index.md</code></pre><p><code>content/_index.md</code>を開いて<code>Hello, Hugo!</code>という文字列を追記しましょう。追記したら、フロントマターの<code>draft</code>を<code>false</code>に変更してコンテンツが公開されるようにします。</p><pre><code class="language-markdown">+++
title = &apos;&apos;
date = 2023-10-09T01:18:58+09:00
draft = false
+++

Hello, Hugo!</code></pre><h1 id="he4943fbc48">Tailwindの動作確認</h1><p>Hugoのサーバーを起動して、Tailwind CSS 3が動作することを確認します。</p><pre><code class="language-bash">hugo serve</code></pre><p><code>http://localhost:1313</code>にアクセスすると、赤く大きい文字で<code>Hello, Hugo!</code>と表示されます。</p><figure><img src="https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/cea74af6f54842f486e8d424128d18fc/%E3%82%B9%E3%82%AF%E3%83%AA%E3%83%BC%E3%83%B3%E3%82%B7%E3%83%A7%E3%83%83%E3%83%88%202023-10-09%201.31.07.png" alt="" width="3024" height="1890"></figure><h1 id="h3de35099b3">参考</h1><ul><li><a href="https://gohugo.io/getting-started/configuration/#configure-cache-busters">https://gohugo.io/getting-started/configuration/#configure-cache-busters</a></li></ul>
