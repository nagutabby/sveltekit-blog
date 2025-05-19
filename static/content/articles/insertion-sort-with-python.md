---
title: Pythonで挿入ソートを実装する
image: images/Microsoft-Fluentui-Emoji-3d-Chart-Increasing-3d.1024.png
publishedAt: 2024-05-03
updatedAt: 2024-05-03
---
# 概要

1.  ソート済み範囲の末尾から先頭に向かって、ソート範囲の左端の要素とソート済みの要素を比較する
2.  ソート範囲の左端の要素がソート済みの要素よりも小さければ要素を交換する
3.  ソート範囲の左端の要素が、ソート済みの要素と同じかソート済みの要素よりも大きければ5に進む
4.  2を繰り返してソート範囲の左端の要素がソート済み範囲の先頭に到達するか、あるいは3の条件を満たしたら要素を適切な位置に挿入できたことになる
5.  ソート範囲の左端を1つ狭める
6.  ソート範囲の要素数が0になるまで2から5を繰り返す

# 特徴

-   in-placeアルゴリズムである
-   計算量は O (n^2) である
-   入力データが昇順に近ければ近いほど、ソート範囲の左端の要素とソート済み範囲の要素の比較回数や交換回数が減少するため、バブルソートや選択ソートよりも高速である

# 実装

```python
import random

def insertion_sort(nums: list[int]) -> list[int]:
    for i in range(len(nums)):
        for j in range(i, 0, -1):
            if nums[j] < nums[j - 1]:
                nums[j], nums[j - 1] = nums[j - 1], nums[j]
            else:
                break

    return nums

nums: list[int] = []

while len(nums) < 10:
    num: int = random.randint(1, 100)
    if num not in nums:
        nums.append(num)

print(f'入力: {nums}')
print(f'出力: {insertion_sort(nums)}')
```
