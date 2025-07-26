---
title: RubyはSmalltalkから何を学んだのか？
image: images/Microsoft-Fluentui-Emoji-3d-Speaker-Low-Volume-3d.1024.png
publishedAt: 2025-07-19
updatedAt: 2025-07-19
---

## はじめに
Rubyの作者が言語設計において最も影響を受けた言語の1つがSmalltalkです。Smalltalkは1970年代に純粋なオブジェクト指向プログラミング言語として登場しました。Smalltalkの「全てがオブジェクト」という純粋なオブジェクト指向の思想は、Rubyの設計思想として息づいています。

この記事では、RubyがSmalltalkから継承した具体的な要素を調査し、それらがどのように現代的なプログラミング言語として実装されたのかを探っていきます。

## オブジェクト
SmalltalkもRubyも、全てがオブジェクトとして表現されます。
Smalltalk:
```smalltalk
Transcript 
    show: 3 class; cr;
    show: true class; cr;
    show: nil class; cr;
    show: Object class; cr.

"
SmallInteger
True
UndefinedObject
Object class
"
```
Ruby:
```ruby
puts 3.class
puts true.class
puts nil.class
puts Object.class

# Integer
# TrueClass
# NilClass
# Class
```
これはシンプルで分かりやすいですね。classメソッドを呼び出すことで、それぞれの値がオブジェクトであることを確認しています。
## メッセージパッシング
メッセージパッシングは、オブジェクト同士がメッセージを送り合うことで処理を実行する仕組みです。メッセージパッシングでは、オブジェクトはレシーバー（受信者）としてメッセージを受信し、そのメッセージに対応するメソッドを実行します。
Smalltalk:
```smalltalk
Transcript
  show: 5 squared; cr;
  show: 5 + 3; cr;
  show: 'hello' size; cr;
  show: 'hello' , ' world'; cr.
```
Ruby:
```ruby
puts 5 ** 2
puts 5 + 3
puts 'hello'.length
puts 'hello' + ' world'
```
Smalltalkは算術演算においてもメッセージをそのまま使用する一方で、Rubyは演算子記法を使用してメッセージパッシングを行います。

## クラス拡張
SmalltalkやRubyでは、実行時にクラスにメソッドを追加できます。これにより、例えば既存のライブラリやフレームワークのクラスを、元のコードを変更することなく拡張できるため、コードの柔軟性が向上します。
Smalltalk:
```smalltalk
Integer compile: 'isPrime
    | i |
    Transcript show: ''Checking if '', self asString, '' is prime''; cr.
    self < 2 ifTrue: [ 
        Transcript show: self asString, '' is less than 2, not prime''; cr.
        ^ false 
    ].
    i := 2.
    [ i * i <= self ] whileTrue: [
        Transcript show: ''Testing divisor: '', i asString; cr.
        self \\ i = 0 ifTrue: [ 
            Transcript show: self asString, '' is divisible by '', i asString, '', not prime''; cr.
            ^ false 
        ].
        i := i + 1
    ].
    Transcript show: self asString, '' is prime!''; cr.
    ^ true'.

12 isPrime.
2 isPrime.
```
Ruby:
```ruby
class Integer
  def prime?(verbose: true)
    log = ->(msg) { puts msg if verbose }
    
    log.call("Checking if #{self} is prime")
    
    if self < 2
      log.call("#{self} is less than 2, not prime")
      return false
    end
    
    (2..Math.sqrt(self)).each do |i|
      log.call("Testing divisor: #{i}")
      if self % i == 0
        log.call("#{self} is divisible by #{i}, not prime")
        return false
      end
    end
    
    log.call("#{self} is prime!")
    true
  end
end

12.prime?
2.prime?
```
Smalltalkでは平方根計算（`i * i <= self`)の処理が分かりにくいですが、RubyではMath.sqrt()を用いられており、より意図が伝わりやすくなっています。また、イテレーターに関しても、Rubyでは範囲オブジェクトが使用されているため、Smalltalkで使用されているwhileTrueと手動インクリメントによる実装に比べて、処理が明示的になっています。
## ブロック構文
処理のまとまりを他のメソッドに渡すための仕組みであるブロック構文は、SmalltalkとRubyの両方で使用でき、両者の間で記法に大きな違いはありません。
Smalltalk:
```smalltalk
#(1 2 3) do: [:x | Transcript show: x; cr]
```
Ruby:
```ruby
[1, 2, 3].each { |x| puts x }
```
## メタプログラミング
メタプログラミングは、プログラムがプログラム自身を操作する技法です。Smalltalkでは[リフレクション](https://e-words.jp/w/%E3%83%AA%E3%83%95%E3%83%AC%E3%82%AF%E3%82%B7%E3%83%A7%E3%83%B3.html)（実行時にプログラム自身の構造を調べたり変更したりする仕組み）、動的ディスパッチ（実行時にどのメソッドを呼び出すかを決定する仕組み）などが実装されています。Rubyでもメタプログラミングが採用されています。
### 動的ディスパッチ
Smalltalk:
```smalltalk
obj := 'hello world'.
methodName := #asUppercase.
result := obj perform: methodName.
Transcript show: result; cr.
result
```
Ruby:
```ruby
obj = 'hello world'
method_name = :upcase
result = obj.send(method_name)
puts result
result
```
## 型システムの考え方
Smalltalkは動的型付け言語であり、柔軟性の高さが言語の特徴でもありましたが、当時はハードウェアの性能が低かったため、Smalltalkはオブジェクト指向プログラミングの概念を示したことによって注目を集めました。一方で、高性能なCPUや多くのメモリを使用できる現在は、オブジェクト指向プログラミングの優れた抽象化に加えて、動的型付け言語の柔軟性も評価されています。もし1970年代から現在までの間で、CPU性能やメモリ量が改善されなかったとしたら、Rubyは現在ほど普及しなかったでしょう。