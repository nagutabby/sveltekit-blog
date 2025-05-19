---
title: "マークダウンで職務経歴書を作る"
image: "https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/fea647b538764d80a31802e1419652b7/Microsoft-Fluentui-Emoji-3d-Bookmark-Tabs-3d.1024.png"
publishedAt: 2024-05-01
updatedAt: 2024-05-04
---

<p>就活をしていると職務経歴書が必要になる場面に遭遇します。</p><p>簡単に職務経歴書を作れるWebサービスを利用するのも手ですが、自分の知的好奇心や技術力をアピールできるものを作りたいと思う方も多いのではないでしょうか？今回は、「ITエンジニアにやさしい職務経歴書」をテーマにマークダウンで職務経歴書を作る方法を解説していきます。</p><h1 id="h5b3e73489f">準備するもの</h1><ul><li>VSCode</li><li><a href="https://marketplace.visualstudio.com/items?itemName=tomoki1207.pdf" target="_blank" rel="noopener noreferrer nofollow">vscode-pdf</a></li><li>pnpmなどのNode.jsのパッケージ管理システム</li><li>好きなCSSフレームワーク<ul><li>何を使ったらいいか分からない方にはClassless CSSフレームワークをおすすめします</li></ul></li></ul><h1 id="hd2a56d7ce9">セットアップ</h1><p>プロジェクトディレクトリを作ります。</p><pre><code class="language-bash">mkdir ~/resume-example
cd ~/resume-example</code></pre><p>Node.jsのプロジェクトを初期化します。</p><pre><code class="language-bash">pnpm init</code></pre><p><a href="https://github.com/simonhaenisch/md-to-pdf" target="_blank" rel="noopener noreferrer nofollow">md-to-pdf</a>をインストールします。</p><pre><code class="language-bash">pnpm i -D md-to-pdf</code></pre><p>職務経歴書として使うマークダウンファイルを作ります。</p><p>動作するかテストするために適当な文字列を書き込んでおきましょう。</p><pre><code class="language-bash">echo &apos;# Hello, World!&apos; &gt; resume.md</code></pre><p>マークダウンファイルからpdfファイルを生成してみます。</p><pre><code class="language-bash">pnpm md-to-pdf resume.md</code></pre><p>例えば、resume.mdを引数として指定した場合はresume.pdfが生成されます。</p><p>VSCodeを開いてresume.pdfの内容を見てみましょう。</p><pre><code class="language-bash">code .</code></pre><figure><img src="https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/b936e3c5b8f942a8aead3a4f7abc17bd/%E3%82%B9%E3%82%AF%E3%83%AA%E3%83%BC%E3%83%B3%E3%82%B7%E3%83%A7%E3%83%83%E3%83%88%202024-05-01%2013.37.01.png" alt="" width="3024" height="1890"></figure><p>md-to-pdfに渡したマークダウンファイルの内容に基づいて、pdfファイルが出力されることが分かります。</p><p></p><p>動作確認ができたら、md-to-pdfのコマンドをパッケージマネージャーにスクリプトとして登録しましょう。スクリプトを登録することで、pdfファイルを生成するためにコマンドを入力する時間を削減できます。</p><p>package.jsonを開いて、<code>&quot;scripts&quot;</code>ブロックを以下の内容に変更します。</p><pre><code class="language-json">&quot;scripts&quot;: {
  &quot;make-resume&quot;: &quot;pnpm md-to-pdf resume.md&quot;
},</code></pre><p><code>make-resume</code>スクリプトを追加すると、以下のコマンドでpdfファイルを生成できるようになります。</p><pre><code class="language-bash">pnpm make-resume</code></pre><h1 id="hcc1af24286">スタイリング</h1><h2 id="h5caccfc304">レイアウトを作る</h2><p>md-to-pdfにはフロントマターで指定できるオプションが用意されています。</p><p>一例として、出力されるpdfファイルのレイアウトを設定する方法を示します。</p><p>resume.mdを以下のように変更します。</p><pre><code class="language-markdown">---
pdf_options:
  format: A4
  margin: 20mm
---

# Hello, World!
This is a dummy text.

This is a dummy text.

This is a [link text](https://example.com/).

This is a dummy text.

This is a dummy text.

## Table
| Header 1 | Header 2 |
| ---- | ---- |
| Data 1 | Data 2 |
| Data 3 | Data 4 |

## Code

```
print(&apos;Hello, World!&apos;)
```

## List
- This is a dummy text.
- This is a dummy text.
- This is a dummy text.</code></pre><figure><img src="https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/0b6c48a048924db693b2fd8f4f4d496e/resume_without_css.jpg" alt="" width="1242" height="1756"></figure><p>上記の例では、A4のページに上下左右20mmの余白が設けられます。</p><p>オプションの詳細は<a href="https://github.com/simonhaenisch/md-to-pdf/blob/master/readme.md" target="_blank" rel="noopener noreferrer nofollow">md-to-pdfのREADME</a>を参照してください。</p><h2 id="h96054fbb83">デザインを工夫する</h2><p>今のままではデザインが味気ないですね。CSSフレームワークを使って見栄えを良くしましょう。</p><p>stylesheetオプションにCSSファイルのパスを指定することでCSSを読み込めます。</p><pre><code class="language-markdown">---
pdf_options:
  format: A4
  margin: 20mm
stylesheet:
  - https://cdnjs.cloudflare.com/ajax/libs/picocss/2.0.6/pico.classless.indigo.min.css
---

# Hello, World!
This is a dummy text.

This is a dummy text.

This is a [link text](https://example.com/).

This is a dummy text.

This is a dummy text.

## Table
| Header 1 | Header 2 |
| ---- | ---- |
| Data 1 | Data 2 |
| Data 3 | Data 4 |

## Code

```
print(&apos;Hello, World!&apos;)
```

## List
- This is a dummy text.
- This is a dummy text.
- This is a dummy text.</code></pre><figure><img src="https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/50f606e16c214048b6851d6080cb9b3c/resume_with_css.jpg" alt="" width="1242" height="1756"></figure><p>モダンなデザインになりました。</p><h1 id="h9be0c3393d">おわりに</h1><p><a href="https://github.com/nagutabby/resume" target="_blank" rel="noopener noreferrer nofollow">マークダウンで作る職務経歴書の例</a>をGitHubで公開しています。参考までにどうぞ。</p><p></p><p><br></p><p></p>