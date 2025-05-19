---
title: Pythonで挿入ソートを実装する
image: images/Microsoft-Fluentui-Emoji-3d-Chart-Increasing-3d.1024.png
publishedAt: 2024-05-03T00:00:00.000Z
updatedAt: 2024-05-03T00:00:00.000Z
---

<h1 id="h9707d3a59a">概要</h1><ol><li>ソート済み範囲の末尾から先頭に向かって、ソート範囲の左端の要素とソート済みの要素を比較する</li><li>ソート範囲の左端の要素がソート済みの要素よりも小さければ要素を交換する</li><li>ソート範囲の左端の要素が、ソート済みの要素と同じかソート済みの要素よりも大きければ5に進む</li><li>2を繰り返してソート範囲の左端の要素がソート済み範囲の先頭に到達するか、あるいは3の条件を満たしたら要素を適切な位置に挿入できたことになる</li><li>ソート範囲の左端を1つ狭める</li><li>ソート範囲の要素数が0になるまで2から5を繰り返す</li></ol><h1 id="hdadc0eaacf">特徴</h1><ul><li>in-placeアルゴリズムである</li><li>計算量は O (n^2) である</li><li>入力データが昇順に近ければ近いほど、ソート範囲の左端の要素とソート済み範囲の要素の比較回数や交換回数が減少するため、バブルソートや選択ソートよりも高速である</li></ul><h1 id="h922edff87b">実装</h1><pre><code class="language-python">import random

def insertion_sort(nums: list[int]) -&gt; list[int]:
    for i in range(len(nums)):
        for j in range(i, 0, -1):
            if nums[j] &lt; nums[j - 1]:
                nums[j], nums[j - 1] = nums[j - 1], nums[j]
            else:
                break

    return nums

nums: list[int] = []

while len(nums) &lt; 10:
    num: int = random.randint(1, 100)
    if num not in nums:
        nums.append(num)

print(f&apos;入力: {nums}&apos;)
print(f&apos;出力: {insertion_sort(nums)}&apos;)</code></pre>
