---
title: Pythonでバブルソートを実装する
image: images/Microsoft-Fluentui-Emoji-3d-Chart-Increasing-3d.1024.png
publishedAt: 2024-05-01T00:00:00.000Z
updatedAt: 2024-05-03T00:00:00.000Z
---
# 概要

1.  ソート範囲の末尾から順に隣り合う要素の大小を比較する
2.  右側の要素が左側の要素よりも小さければ2つの要素を交換する
3.  ソート範囲の左端に到達するまで1から2を繰り返す
4.  ソート範囲の左端に到達したら、ソート範囲を左から1つ狭める
5.  ソート範囲に含まれる要素数が1になるまで1から4を繰り返す

# 特徴

-   in-placeアルゴリズムである
-   計算量はO (n^2) である

# 実装

```python
import random

def bubble_sort(nums: list[int]) -> list[int]:
    for i in range(len(nums) - 1):
        for j in range(len(nums) - 1, i, -1):
            if nums[j] < nums[j - 1]:
                nums[j], nums[j - 1] = nums[j - 1], nums[j]

    return nums

nums: list[int] = []

while len(nums) < 10:
    num: int = random.randint(1, 100)
    if num not in nums:
        nums.append(num)

print(f'入力: {nums}')
print(f'出力: {bubble_sort(nums)}')
```