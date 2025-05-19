---
title: NTPによる時刻同期の仕組みを調べてみた
image: images/Microsoft-Fluentui-Emoji-3d-Alarm-Clock-3d.1024.png
publishedAt: 2023-07-05T00:00:00.000Z
updatedAt: 2024-05-01T00:00:00.000Z
---
# 時刻同期

時刻同期(Time Synchronization)は、あるネットワーク上のホストの時刻を、ある基準に従って同期させることです。

基準の例としては、JSTやUTCなどのタイムゾーンがあります。

# NTP

NTP(Network Time Protocol)は、[RFC5095](https://datatracker.ietf.org/doc/html/rfc5905)などで標準化されている、機器の時刻を同期するためのプロトコルです。日本ではNICT(情報通信研究機構)などがNTPサーバを提供しています。NTPは、サーバクライアントモデルで用いられ、NTPクライアントが複数のNTPサーバに時刻を問い合わせることでより正確な時刻同期を行います。

1.  NTPクライアントが複数のNTPサーバに時刻を問い合わせる
2.  NTPサーバから受け取った時刻データを参照し、明らかに時刻がずれているNTPサーバ(falseticker)があればそれを参照先リストから除外する
3.  残ったNTPサーバのリストから時刻同期の精度が高い方から順に最大で3つまでNTPサーバを選択する

## truechimer

truechimerは、NTPのアルゴリズムによって、正しい時刻を配信していると判定されたNTPサーバです。

## falseticker

falsetickerは、NTPのアルゴリズムによって、誤った時刻を配信していると判定されたNTPサーバです。

## NTPサーバの精度

NTPサーバは、DNS(Domain Name System)サーバと同様に、負荷を軽減させるために階層構造になっています。信頼できる階層は全部で**16層**あり、それぞれの階層の総称を**Stratum**と呼び、最上位の階層から順にStratum0、Stratum1、…Stratum15と呼びます。

### Stratum0

Stratum0はGPS(Global Positioning System)あるいは原子時計と直接同期しているNTPサーバの階層です。NTPサーバにおいて、時刻同期の精度が最も高いことが多いです。

### Stratum1〜15

Stratum1はStratum0のNTPサーバを参照して時刻を同期しているNTPサーバの階層です。

StratumNのNTPサーバは、StratumN-1のNTPサーバを参照しています。

## 時刻データの補正

一般的に、NTPサーバとNTPクライアントは地理的に離れています。例えば、NICTのNTPサーバ([ntp.nict.jp](http://ntp.nict.jp))は、兵庫県と東京都で動作しています。これでは、ネットワークの遅延などが起きていた場合に、NTPクライアントが受信した時刻データの正確性が下がってしまいます。できるだけこれらの影響を小さくするために、NTPクライアントは以下の値を計算してNTPサーバから受信した時刻データを補正します。

-   NTPサーバにおける処理による遅延時間
-   NTPサーバとNTPクライアントの間のネットワークによる遅延時間
-   用いられる指標: RTT、(オフセット)
-   NTPサーバとNTPクライアントの間のネットワークの遅延の安定性
-   用いられる指標: ジッター

### オフセット

#### サーバとクライアントの間の時刻の差が０であるとき

実際のシステムではありえませんが、サーバとクライアントの時刻が完全に同期されている場合は、オフセットはRTTを2で割った値になります。

#### サーバとクライアントの間の時刻の差が０ではないとき

実際のシステムでは、サーバとクライアントの間の時刻はずれています。この場合で、かつネットワークの遅延が行きと帰りで等しいとみなせる場合に、オフセットはサーバとクライアントの間の時刻の差と**みなせます**。

この際のオフセットがどれくらい信頼できるかはジッターの大きさによります。

### RTT

RTT(Round Trip Time)は、サーバとクライアントの間でデータを伝送する際の**往復**にかかる時間です。

### ジッター

ジッター(Jitter)はネットワークの遅延の安定性を表す指標です。ジッターが小さいほどネットワークの遅延が安定しており、行きと帰りのネットワークの遅延時間が近いため、**サーバとクライアントの間の時刻の差が０ではないとき**のオフセットを算出する際の精度が高くなります。

# SNTP

SNTP(Simple Network Time Protocol)は、[RFC4330](https://datatracker.ietf.org/doc/rfc4330/)で標準化されている、NTPよりもシンプルに時刻同期を行うためのプロトコルです。NTPでは階層構造を持つ複数のNTPサーバに時刻を問い合わせて、より精度が高い時刻同期を行いますが、SNTPでは階層構造を持たない1つのNTPサーバにのみ時刻を問い合わせます。ミリ秒単位の時刻同期が不要なコンピュータではSNTPを使用することもできます。

# 参照先のNTPサーバはいくつ用意すべきか

NTPを正しく機能させるためには、最低3つのNTPサーバが必要です。NTPの開発チームは4つのNTPサーバから時刻を取得することを推奨しています。

なぜなら、NTPのアルゴリズムでは、多数決によってfalsetickerを選んでおり、3つ以上のNTPサーバが無いとfalsetickerを選ぶことができないからです。

また、4つのNTPサーバを設定することで、NTPサーバの冗長構成が保証されます。

# NTPサーバ、NTPクライアント

## systemd-timesyncd

systemdが動作しているコンピュータにおいて、SNTPクライアントとして動作するデーモンです。

## ntpd(Network Time Protocol daemon)

NTPサーバ、NTPクライアントとして動作するデーモンです。

## chronyd

NTPサーバ、NTPクライアントとして動作するChronyのデーモンです。RedHat系ディストリビューションでは、デフォルトでインストールされています。

# 参考

-   [【図解】NTPプロトコルの概要と仕組み～誤差補正の計算,仕様,シーケンス～](https://milestone-of-se.nesuke.com/l7protocol/ntp/ntp-summary/)
-   [SelectingOffsiteNTPServers < Support < Network Time Foundation's NTP Support Wiki](https://support.ntp.org/Support/SelectingOffsiteNTPServers#Section_5.3.2.)
-   [systemd-timesyncd - ArchWiki](https://wiki.archlinux.jp/index.php/Systemd-timesyncd)
-   [Network Time Protocol daemon - ArchWiki](https://wiki.archlinux.jp/index.php/Network_Time_Protocol_daemon)
-   [chrony - ArchWiki](https://wiki.archlinux.jp/index.php/Chrony)