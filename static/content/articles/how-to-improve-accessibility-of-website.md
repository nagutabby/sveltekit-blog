---
title: "Webアプリのアクセシビリティを向上させる方法"
image: "https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/93654bdb0a4b4d6ba654235ffecfc643/Microsoft-Fluentui-Emoji-3d-Mobile-Phone-3d.1024.png"
publishedAt: 2023-07-05
updatedAt: 2024-05-01
---

<h1 id="h4f7f573348">文章を意味ごとに区切る</h1><p>スクリーンリーダーが読み上げを適切な位置で区切ってくれるように、文章を意味ごとに区切ります。</p><p>不適切な例:</p><pre><code class="language-html">This is a first&lt;br&gt;
sentence.</code></pre><p>適切な例:</p><pre><code class="language-html">&lt;p&gt;This is a first sentence.&lt;/p&gt;
</code></pre><h1 id="hdc78befd33">複雑な表現を使わない</h1><p>認知障害を抱える人のために、専門用語やスラングなどの複雑な表現をできるだけ使わないようにします。</p><p>IT業界には、ときどき自分のブログでカッコつけてTL;DRという用語を使う人がいますが、ああいう表現はしないようにしましょう、ということです。</p><h1 id="h753a22adef">入れ子構造を減らす</h1><p>HTMLの入れ子構造を複雑にすると、スクリーンリーダーが間違った読み上げをする可能性が高くなります。tableタグなどの入れ子が深くなりやすい要素は、できるだけ使うのを避けましょう。</p><h1 id="had4b67d0e8">コードの役割を明確にする</h1><p>例えば、ボタンを表したい場合はaタグやdivタグではなく、buttonタグやinputタグで実装しましょう。どうしてもbuttonタグを使いたくない場合は、<a href="https://developer.mozilla.org/ja/docs/Web/Accessibility/ARIA/Roles" target="_blank" rel="noopener noreferrer nofollow"><u>WAI-ARIAロール</u></a>のbuttonロールを使いましょう。</p><h1 id="h8f303ec441">ラベルに意味を持たせる</h1><p>aタグなどのラベルのみがスクリーンリーダーに読み上げられる場合があります。そのため、ラベルと関連付けられたものをラベルのみで説明できるようにします。</p><p>不適切な例:</p><pre><code class="language-html">&lt;p&gt;
  Please see&lt;a href=&quot;information.html&quot;&gt;here&lt;/a&gt; to check today&apos;s menu.
&lt;/p&gt;</code></pre><p>適切な例:</p><pre><code class="language-html">&lt;p&gt;
  &lt;a href=&quot;information.html&quot;&gt;Please check today&apos;s menu&lt;/a&gt;.
&lt;/p&gt;</code></pre><p>また、textarea要素などのフォーム要素の説明には、labelタグを使うことで、Webブラウザがそれらを1つのまとまりとして認識してくれます。</p><h1 id="hd64eeae975">代替テキストを使う</h1><p>画像や動画を説明するためのテキストである代替テキストを使い、視覚障害者や聴覚障害者がWebアプリの内容をより理解できるようにします。</p><h2 id="hb9cc2f8bb1">alt属性</h2><p>alt属性は画像を説明するためのものです。スクリーンリーダーで読み上げられるため、簡潔で分かりやすい説明をすることを心がけましょう。</p><p>不適切な例:</p><pre><code class="language-html">&lt;img href=&quot;images/grilled-fish.jpg&quot;&gt;
&lt;p&gt;It&apos;s tasty!&lt;/p&gt;</code></pre><p>適切な例:</p><pre><code class="language-html">&lt;img href=&quot;images/grilled-fish.jpg&quot; alt=&quot;Grilled salmon on the while plate.&quot;&gt;
&lt;p&gt;It&apos;s tasty!&lt;/p&gt;
</code></pre><h1 id="h3de35099b3">参考</h1><ul><li><a href="https://developer.mozilla.org/ja/docs/Learn/Accessibility/HTML" target="_blank" rel="noopener noreferrer nofollow"><u>HTML: アクセシビリティの基礎 - ウェブ開発を学ぶ | MDN</u></a></li><li><a href="https://developer.mozilla.org/ja/docs/Web/Accessibility/ARIA" target="_blank" rel="noopener noreferrer nofollow"><u>ARIA - アクセシビリティ | MDN</u></a></li></ul>