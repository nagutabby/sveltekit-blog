---
title: "Pythonで二分探索を実装する"
image: "https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/bdf32d37ceeb488f8e5bda51690dea93/Microsoft-Fluentui-Emoji-3d-Magnifying-Glass-Tilted-Left-3d.1024.png"
publishedAt: 2024-05-14
updatedAt: 2024-05-15
---

<h1 id="h9707d3a59a">概要</h1><ol><li>昇順でソートされた配列を用意する</li><li>値の探索範囲の真ん中の値を選び、探索する値と比較する</li><li>選んだ配列の要素が探索する値と異なれば、値の探索範囲を変える<ol><li>選んだ配列の要素が探索する値よりも小さければ、探索する値は選んだ配列の要素の右にあるので、探索範囲の左端を選んだ配列の要素の1つ右の要素とする。</li><li>選んだ配列の要素が探索する値よりも大きければ、探索する値は選んだ配列の要素の左にあるので、探索範囲の右端を選んだ配列の要素の1つ左の要素とする。</li></ol></li><li>2から3の操作を繰り返す</li></ol><h1 id="hdadc0eaacf">特徴</h1><ul><li>計算量はO (log n) である</li></ul><h1 id="h922edff87b">実装</h1><pre><code class="language-python">import random

def binary_search(nums: list[int], num: int) -&gt; int:
    left: int = 0
    right: int = len(nums) - 1
    mid: int = len(nums) // 2

    while nums[mid] != num:
        if nums[mid] &lt; num:
            left = mid + 1
        else:
            right = mid - 1
        mid = (left + right) // 2

    return mid

nums: list[int] = []

while len(nums) &lt; 10:
    num: int = random.randint(1, 100)
    if num not in nums:
        nums.append(num)

nums.sort()

print(f&apos;リスト: {nums}, 探索する値: {num}&apos;)
print(f&apos;値がマッチしたインデックス: {binary_search(nums, num)}&apos;)</code></pre>