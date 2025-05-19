---
title: "Pythonで選択ソートを実装する"
image: "https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/d4219785c85749ea8db95d86b998f2b8/Microsoft-Fluentui-Emoji-3d-Chart-Increasing-3d.1024.png"
publishedAt: 2024-05-03
updatedAt: 2024-05-03
---

<h1 id="h9707d3a59a">概要</h1><ol><li>ソート範囲における最小値を探す</li><li>最小値をソート範囲の左端の要素と交換する</li><li>ソート範囲を左から1つ狭める</li><li>ソート範囲に含まれる要素数が1になるまで1から3を繰り返す</li></ol><h1 id="hdadc0eaacf">特徴</h1><ul><li>in-placeアルゴリズムである</li><li>計算量はO (n^2) である</li><li>バブルソートと比べると、要素の比較回数は同じであるが、要素の交換回数が少ないためバブルソートよりも高速である</li></ul><h1 id="h922edff87b">実装</h1><pre><code class="language-python">import random

def selection_sort(nums: list[int]) -&gt; list[int]:
    for i in range(len(nums) - 1):
        current_min_index: int = i
        for j in range(i + 1, len(nums)):
            if nums[current_min_index] &gt; nums[j]:
                current_min_index = j

        nums[i], nums[current_min_index] = nums[current_min_index], nums[i]

    return nums

nums: list[int] = []

while len(nums) &lt; 10:
    num: int = random.randint(1, 100)
    if num not in nums:
        nums.append(num)

print(f&apos;入力: {nums}&apos;)
print(f&apos;出力: {selection_sort(nums)}&apos;)</code></pre>