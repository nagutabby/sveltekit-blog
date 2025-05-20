---
title: Webアプリのアクセシビリティを向上させる方法
image: images/Microsoft-Fluentui-Emoji-3d-Mobile-Phone-3d.1024.png
publishedAt: 2023-07-05
updatedAt: 2024-05-01
---
# 文章を意味ごとに区切る

スクリーンリーダーが読み上げを適切な位置で区切ってくれるように、文章を意味ごとに区切ります。

不適切な例:

```html
This is a first<br>
sentence.
```

適切な例:

```html
<p>This is a first sentence.</p>
```

# 複雑な表現を使わない

認知障害を抱える人のために、専門用語やスラングなどの複雑な表現をできるだけ使わないようにします。

IT業界には、ときどき自分のブログでカッコつけてTL;DRという用語を使う人がいますが、ああいう表現はしないようにしましょう、ということです。

# 入れ子構造を減らす

HTMLの入れ子構造を複雑にすると、スクリーンリーダーが間違った読み上げをする可能性が高くなります。tableタグなどの入れ子が深くなりやすい要素は、できるだけ使うのを避けましょう。

# コードの役割を明確にする

例えば、ボタンを表したい場合はaタグやdivタグではなく、buttonタグやinputタグで実装しましょう。どうしてもbuttonタグを使いたくない場合は、[WAI-ARIAロール](https://developer.mozilla.org/ja/docs/Web/Accessibility/ARIA/Roles)のbuttonロールを使いましょう。

# ラベルに意味を持たせる

aタグなどのラベルのみがスクリーンリーダーに読み上げられる場合があります。そのため、ラベルと関連付けられたものをラベルのみで説明できるようにします。

不適切な例:

```html
<p>
  Please see<a href="information.html">here</a> to check today's menu.
</p>
```

適切な例:

```html
<p>
  <a href="information.html">Please check today's menu</a>.
</p>
```

また、textarea要素などのフォーム要素の説明には、labelタグを使うことで、Webブラウザがそれらを1つのまとまりとして認識してくれます。

# 代替テキストを使う

画像や動画を説明するためのテキストである代替テキストを使い、視覚障害者や聴覚障害者がWebアプリの内容をより理解できるようにします。

## alt属性

alt属性は画像を説明するためのものです。スクリーンリーダーで読み上げられるため、簡潔で分かりやすい説明をすることを心がけましょう。

不適切な例:

```html
<img href="images/grilled-fish.jpg">
<p>It's tasty!</p>
```

適切な例:

```html
<img href="images/grilled-fish.jpg" alt="Grilled salmon on the while plate.">
<p>It's tasty!</p>
```

# 参考

-   [HTML: アクセシビリティの基礎 - ウェブ開発を学ぶ | MDN](https://developer.mozilla.org/ja/docs/Learn/Accessibility/HTML)
-   [ARIA - アクセシビリティ | MDN](https://developer.mozilla.org/ja/docs/Web/Accessibility/ARIA)
