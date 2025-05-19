---
title: Pythonで二分探索を実装する
image: images/Microsoft-Fluentui-Emoji-3d-Magnifying-Glass-Tilted-Left-3d.1024.png
publishedAt: 2024-05-14T00:00:00.000Z
updatedAt: 2024-05-15T00:00:00.000Z
---
# 概要

1.  昇順でソートされた配列を用意する
2.  値の探索範囲の真ん中の値を選び、探索する値と比較する
3.  選んだ配列の要素が探索する値と異なれば、値の探索範囲を変える
    1.  選んだ配列の要素が探索する値よりも小さければ、探索する値は選んだ配列の要素の右にあるので、探索範囲の左端を選んだ配列の要素の1つ右の要素とする。
    2.  選んだ配列の要素が探索する値よりも大きければ、探索する値は選んだ配列の要素の左にあるので、探索範囲の右端を選んだ配列の要素の1つ左の要素とする。
4.  2から3の操作を繰り返す

# 特徴

-   計算量はO (log n) である

# 実装

```python
import random

def binary_search(nums: list[int], num: int) -> int:
    left: int = 0
    right: int = len(nums) - 1
    mid: int = len(nums) // 2

    while nums[mid] != num:
        if nums[mid] < num:
            left = mid + 1
        else:
            right = mid - 1
        mid = (left + right) // 2

    return mid

nums: list[int] = []

while len(nums) < 10:
    num: int = random.randint(1, 100)
    if num not in nums:
        nums.append(num)

nums.sort()

print(f'リスト: {nums}, 探索する値: {num}')
print(f'値がマッチしたインデックス: {binary_search(nums, num)}')
```