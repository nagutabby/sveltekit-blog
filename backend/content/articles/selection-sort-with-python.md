---
title: Pythonで選択ソートを実装する
image: images/Microsoft-Fluentui-Emoji-3d-Chart-Increasing-3d.1024.png
publishedAt: 2024-05-03
updatedAt: 2024-05-03
---
# 概要

1.  ソート範囲における最小値を探す
2.  最小値をソート範囲の左端の要素と交換する
3.  ソート範囲を左から1つ狭める
4.  ソート範囲に含まれる要素数が1になるまで1から3を繰り返す

# 特徴

-   in-placeアルゴリズムである
-   計算量はO (n^2) である
-   バブルソートと比べると、要素の比較回数は同じであるが、要素の交換回数が少ないためバブルソートよりも高速である

# 実装

```python
import random

def selection_sort(nums: list[int]) -> list[int]:
    for i in range(len(nums) - 1):
        current_min_index: int = i
        for j in range(i + 1, len(nums)):
            if nums[current_min_index] > nums[j]:
                current_min_index = j

        nums[i], nums[current_min_index] = nums[current_min_index], nums[i]

    return nums

nums: list[int] = []

while len(nums) < 10:
    num: int = random.randint(1, 100)
    if num not in nums:
        nums.append(num)

print(f'入力: {nums}')
print(f'出力: {selection_sort(nums)}')
```
