---
title: ファクトリーメソッドの使い方
image: images/Microsoft-Fluentui-Emoji-3d-Factory-3d.1024.png
publishedAt: 2024-06-30T00:00:00.000Z
updatedAt: 2024-06-30T00:00:00.000Z
---
# ファクトリーメソッド

ファクトリーメソッド (Factory Method) はその名の通り、製品を生産する工場のような役割を持つメソッドです。例えば、缶詰を生産する工場を考えてみます。缶詰を作るには、まず缶に入れる食品を加工する必要があります。さらに、加工した食品を缶に詰めなければなりません。しかし、工場の機械化が進んだ現代では、食品の加工方法や食品を缶に詰める方法を知らなくても、缶詰の生産数を入力すれば機械が代わりに缶詰を作ってくれます。この仕組みをオブジェクト指向プログラミングで実現する方法がファクトリーメソッドです。

# Pyhonで鯖の缶詰を作る

具体例として、鯖の缶詰を作ってみましょう。鯖の缶詰を作るには、鯖、鯖の缶詰、鯖の缶詰を作る人が必要です。これらの役割を、Mackerelクラス、CannedMackerelクラス、Workerクラスに順に割り当てます。

```python
class Mackerel:
    def __init__(self) -> None:
        # 色
        self.color: str = 'blue'
        # 重さ
        self.weight: int = 200
```

```python
from mackerel import Mackerel

class CannedMackerel:
    def __init__(self) -> None:
        # 缶詰の中身: 鯖2匹
        self.content: list[Mackerel] = [Mackerel(), Mackerel()]

    # 缶詰をラッピングする
    def wrap(self, material: str) -> None:
        self.wrapping: str = material
```

```python
from canned_mackerel import CannedMackerel

class Worker:
    # 鯖の缶詰を作る
    def produce_canned_mackerel(self) -> CannedMackerel:
        # 鯖の缶詰を用意する
        canned_mackerel: CannedMackerel = CannedMackerel()
        # 缶詰を木箱でラッピングする
        canned_mackerel.wrap('crate')
        return canned_mackerel
```

最後に、Workerクラスのproduce\_canned\_mackerelメソッドを呼び出します。

```python
from worker import Worker
from canned_mackerel import CannedMackerel

canned_mackerel: CannedMackerel = Worker().produce_canned_mackerel()

print(vars(canned_mackerel))
```

# 利点

先ほど紹介したMackerelクラス、CannedMackerelクラス、Workerクラスはそれぞれ以下の役割を持ちます。

-   Mackerelクラス: 生成されるオブジェクトに共通する要素を定義する
-   CannedMackerelクラス: 実際に生成されるオブジェクトを定義する
-   Workerクラス: ファクトリーメソッドを定義する

それぞれのクラスはたった1つの役割を持っているため、このコードは単一責任の原則を満たしており、クラス同士が疎結合になっています。また、もしも鯖を缶ではなく瓶に詰めたくなったとしても、新しくBottledMackerelクラスを作ることで対応できるため、コードの拡張性が高いです。さらに、Workerクラスがproduce\_canned\_mackerelメソッドを呼び出すだけで鯖の缶詰が作られるため、CannedMackerelクラスがカプセル化されていることが分かります。これにより、鯖の缶詰の作り方を変更したくなったとしてもコードをわずかに修正するだけで済みます。