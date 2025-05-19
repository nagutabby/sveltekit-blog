---
title: Pythonでバブルソートを実装する
image: images/Microsoft-Fluentui-Emoji-3d-Chart-Increasing-3d.1024.png
publishedAt: 2024-05-01T00:00:00.000Z
updatedAt: 2024-05-03T00:00:00.000Z
---

<h1 id="h9707d3a59a">概要</h1><ol><li>ソート範囲の末尾から順に隣り合う要素の大小を比較する</li><li>右側の要素が左側の要素よりも小さければ2つの要素を交換する</li><li>ソート範囲の左端に到達するまで1から2を繰り返す</li><li>ソート範囲の左端に到達したら、ソート範囲を左から1つ狭める</li><li>ソート範囲に含まれる要素数が1になるまで1から4を繰り返す</li></ol><h1 id="hdadc0eaacf">特徴</h1><ul><li>in-placeアルゴリズムである</li><li>計算量はO (n^2) である</li></ul><h1 id="h922edff87b">実装</h1><pre><code class="language-python">import random

def bubble_sort(nums: list[int]) -&gt; list[int]:
    for i in range(len(nums) - 1):
        for j in range(len(nums) - 1, i, -1):
            if nums[j] &lt; nums[j - 1]:
                nums[j], nums[j - 1] = nums[j - 1], nums[j]

    return nums

nums: list[int] = []

while len(nums) &lt; 10:
    num: int = random.randint(1, 100)
    if num not in nums:
        nums.append(num)

print(f&apos;入力: {nums}&apos;)
print(f&apos;出力: {bubble_sort(nums)}&apos;)</code></pre><p></p>
