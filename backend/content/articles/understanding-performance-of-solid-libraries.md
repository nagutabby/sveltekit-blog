---
title: 'Solid Cache, Solid Queue, Solid Cableのパフォーマンスを理解する'
image: images/Microsoft-Fluentui-Emoji-3d-Open-Mailbox-With-Raised-Flag-3d.1024.png
publishedAt: 2025-03-15
updatedAt: 2025-03-15
---
# はじめに

Rails環境では長らく、バックグラウンドジョブ処理、キャッシュ処理、リアルタイム通信などの重要なコンポーネントがRedisに依存していました。しかし、Redisサーバーの運用・監視が必要になるという問題点がありました。この依存関係の複雑さを削減するために、Solidライブラリが生まれました。この記事では、Railsエコシステムの中で注目を集めている3つのSolidライブラリについて解説します。

# Solid Cache

Solid Cacheはキャッシュストアの一種です。データベースを使用してキャッシュを保存するため、インメモリキャッシュに比べて性能が低下しますが、より安価な価格でキャッシュ処理を導入できます。

[37signalsの調査](https://dev.37signals.com/solid-cache/)では、Redisキャッシュに比べて読み取り速度は40%低下しますが、Redisをデータベースに置き換えることで運用コストが80%安くなります。

# Solid Queue

Solid Queueはバックグラウンドジョブシステムの一種です。Solid Queueも同様にデータベースを使用するため、I/Oが低下します。しかし、データベースが提供する[SKIP LOCKED構文](https://dev.mysql.com/blog-archive/mysql-8-0-1-using-skip-locked-and-nowait-to-handle-hot-rows/)を用いることで、ロックされた行をスキップしてジョブを取得できるようになっています。そのため、それまでは実現不可能であったジョブのノンブロッキング処理が実現され、データベースを利用したバックグラウンドジョブ処理のパフォーマンスが向上しています。

# Solid Cable

Solid Cableはリアルタイム通信フレームワークの一種です。Solid Cableも同様にデータベースを使用するため、I/Oが低下します。しかし、[k6を用いたベンチマーキング](https://techracho.bpsinc.jp/hachi8833/2024_11_11/146390)で250ユーザーが同時にメッセージを送信する際のパフォーマンスを比較すると、Solid CableとSQLiteを使用した場合とRedisを使用した場合でパフォーマンスの差はほとんどありません。

具体的には、Solid CableとSQLiteを使用した場合の、RTTの95パーセンタイルが150ms、WebSocketへの接続時間の95パーセンタイルがおよそ291msである一方で、Redisを使用した場合はそれぞれ、135ms、およそ273msでした。

# 参考記事

-   [https://dev.37signals.com/introducing-solid-queue/](https://dev.37signals.com/introducing-solid-queue/)
