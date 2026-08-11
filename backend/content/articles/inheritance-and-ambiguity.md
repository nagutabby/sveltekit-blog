---
title: 継承と曖昧さ
image: images/Microsoft-Fluentui-Emoji-3d-Bug-3d.1024.png
publishedAt: 2025-05-17
updatedAt: 2025-05-18
---
# はじめに

オブジェクト指向プログラミング（OOP）の登場とともに、継承という概念が普及しました。継承によってコードの再利用性が高まり、階層的な設計ができるようになりました。この分割的な設計手法は、抽象化の強力な手段の1つとして定着しました。

しかし、継承がプログラムに提供したのは抽象化だけではありません。継承はプログラムに複雑さも導入し、それによって曖昧な基準が生まれました。

# 継承固有の曖昧さ

継承の曖昧な要素には様々なものがあります。有名なものをいくつかご紹介します。

## 菱形継承問題

菱形継承問題（Diamond problem）とは、オブジェクト指向プログラミングにおいて多重継承を使用した際に発生する問題です。基底クラスがあり、そのクラスを継承するサブクラスAとBがある状況を考えます。このときに、サブクラスAとBを継承する派生クラスを作成すると菱形継承が成立します。菱形継承には2つの問題があります。

### 基底クラスの二重継承

基底クラスの変数やメソッドが、派生クラスに二重に継承されます。これにより、基底クラスのメンバーを呼び出すときの経路が不明確になります。

![](images/mermaid-diagram-2025-05-17-141432.png)

### 同名メソッドの多重継承

サブクラスAとBにある同じ名前のメソッドが継承されます。これにより、派生クラスからそのメソッドを呼び出そうとしたときに呼び出し先メソッドを区別できなくなります。

![](images/mermaid-diagram-2025-05-17-143301.png)

## リスコフの置換原則とis-a関係

リスコフの置換原則は、is-a関係と競合することがあります。まずはそれぞれの用語を整理しましょう。

### リスコフの置換原則

リスコフの置換原則（Liskov Substitution principle, LSP）とは、「サブタイプのオブジェクトはスーパータイプのオブジェクトと置き換え可能でなければならない」という原則です。サブタイプとはサブクラス（派生クラス）のことで、スーパータイプとはスーパークラス（基底クラス）やインターフェースのことです。

リスコフは、オブジェクト指向設計の堅牢性を向上させるために、型の階層関係に意味を持たせることが重要であると考えました。そのため、このような原則を提唱したのです。

リスコフの置換原則に従って派生クラスを記述するために、守るべきルールがあります。主要なものは以下の通りです。

-   事前条件（Precondition）について
    -   サブタイプは、スーパータイプよりも強い事前条件を要求してはならない
-   事後条件（Postcondition）について
    -   サブタイプは、スーパータイプよりも弱い事後条件を保証してはならない
-   不変条件（Invariant）について
    -   サブタイプは、スーパータイプの不変条件を維持しなければならない

事前条件とは、メソッドが正しく実行されるために呼び出し前に満たされるべき条件です。例えば、「数値を計算するメソッドの引数はnullであってはならない」、「配列の要素を取得するメソッドのインデックスは配列の範囲内でなければならない」という事前条件があります。

事後条件とは、メソッドが正しく実行された後に保証されるべき条件です。例えば、「配列をソートするメソッドを実行すると、配列の要素は昇順に並ぶ」、「ユーザーを登録するメソッドを実行すると、データベースにはそのユーザーが存在する」という事後条件があります。

不変条件とは、オブジェクトの生存期間を通じて常に真であるべき条件です。例えば、「銀行口座クラス残高は常に0以上である」、「スタッククラスのサイズは常に0以上かつ最大容量以下である」という不変条件があります。

次に、リスコフの置換原則に従った継承関係を例を示します。

```java
// 形状の抽象インターフェース
interface Shape {
    double area();
    double perimeter();
}

class Rectangle implements Shape {
    private double width;
    private double height;

    public Rectangle(double width, double height) {
        if (width <= 0 || height <= 0) {
            throw new IllegalArgumentException("幅と高さは正の値である必要があります");
        }
        this.width = width;
        this.height = height;
    }

    public double getWidth() {
        return width;
    }

    public void setWidth(double width) {
        if (width <= 0) {
            throw new IllegalArgumentException("幅は正の値である必要があります");
        }
        this.width = width;
    }

    public double getHeight() {
        return height;
    }

    public void setHeight(double height) {
        if (height <= 0) {
            throw new IllegalArgumentException("高さは正の値である必要があります");
        }
        this.height = height;
    }

    @Override
    public double area() {
        return width * height;
    }

    @Override
    public double perimeter() {
        return 2 * (width + height);
    }
}

class Square implements Shape {
    private double side;

    public Square(double side) {
        if (side <= 0) {
            throw new IllegalArgumentException("辺の長さは正の値である必要があります");
        }
        this.side = side;
    }

    public double getSide() {
        return side;
    }

    public void setSide(double side) {
        if (side <= 0) {
            throw new IllegalArgumentException("辺の長さは正の値である必要があります");
        }
        this.side = side;
    }

    @Override
    public double area() {
        return side * side;
    }

    @Override
    public double perimeter() {
        return 4 * side;
    }
}
```

この例では、Shapeインターフェースを実装したRectangleクラスとSquareクラスが定義されています。どちらのクラスもShapeインターフェースで定義されたメソッドを全て実装しているため、Shape（形状）として期待される振る舞いを提供することが保証されています。

### is-a関係

is-a関係とは、継承によって得られる、サブタイプとスーパータイプの関係です。例えば、「犬は動物である（Dog is an animal）」という関係は、DogクラスがAnimalクラスを継承していることを表します。このときに、派生クラスであるDogは、基底クラスであるAnimalの全ての特性と振る舞いを継承します。犬は動物なので、食べたり眠ったりできます。

![](images/mermaid-diagram-2025-05-17-155039.png)

### リスコフの置換原則を満たさないis-a関係

一見すると、is-a関係とリスコフの置換原則は両立できそうです。しかし、is-a関係を満たしていても、リスコフの置換原則を満たしていないことがあります。

正方形と長方形はまさにその好例です。正方形は長方形の一種（Square is a Rectangle）しかし、正方形は幅と高さのどちらか一方があれば面積が求まるのに対して、長方形は幅と高さの両方が必要です。そのため、正方形は長方形の事前条件を強化しており、リスコフの置換原則のルールに違反しています。従って正方形と長方形はis-a関係でありながら、リスコフの置換原則を満たしていません。

リスコフの置換原則に違反する例は以下の通りです。

```java
class Rectangle {
    protected int width;
    protected int height;

    public Rectangle(int width, int height) {
        this.width = width;
        this.height = height;
    }

    public void setWidth(int width) {
        this.width = width;
    }

    public void setHeight(int height) {
        this.height = height;
    }

    public int getWidth() {
        return width;
    }

    public int getHeight() {
        return height;
    }

    public int area() {
        return width * height;
    }
}

// LSP違反の例
class Square extends Rectangle {
    public Square(int side) {
        super(side, side);
    }

    // 親クラスの振る舞いを変更している
    @Override
    public void setWidth(int width) {
        super.setWidth(width);
        super.setHeight(width); // 正方形の制約を維持するため
    }

    @Override
    public void setHeight(int height) {
        super.setHeight(height);
        super.setWidth(height); // 正方形の制約を維持するため
    }
}
```
