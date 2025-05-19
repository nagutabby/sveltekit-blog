---
title: Linuxでネットワークを管理したい
image: images/Microsoft-Fluentui-Emoji-3d-Penguin-3d.1024.png
publishedAt: 2023-07-05T00:00:00.000Z
updatedAt: 2024-05-01T00:00:00.000Z
---
# ネットワークを管理するファイル

## /etc/hosts

名前解決をするために使われるファイルです。一行に1つのIPアドレスと1つ以上のホスト名を持ちます。

現在は、DNSによる名前解決と併用されています。

例:

```shell
127.0.1.1 vultr.guest vultr
127.0.0.1 localhost

# The following lines are desirable for IPv6 capable hosts
::1 localhost ip6-localhost ip6-loopback
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
```

## /etc/hostname

ホスト名を設定します。Debian系で使われるファイルです。

例:

```shell
vultr
```

## /etc/nsswitch.conf

名前解決やサービス名の解決を行う際の問い合わせ順序を指定します。

例えば、`hosts`に`files dns`を指定している場合は、`/etc/hosts`を参照し、そのファイルにIPアドレスとホスト名の対応が記載されていなかった場合に、`/etc/resolv.conf`に記載されているDNSによる名前解決を行います。

例:

```
passwd:         files systemd
group:          files systemd
shadow:         files
gshadow:        files

hosts:          files dns
networks:       files

protocols:      db files
services:       db files
ethers:         db files
rpc:            db files

netgroup:       nis
```

# ネットワークを管理するコマンド

## route

ルーティングテーブルを操作します。

```shell
ubuntu@vultr:~$ whatis route
route (8)            - show / manipulate the IP routing table
```

### IPv4のルーティングテーブルを表示

```shell
route
```

```shell
route -4
```

### IPv6のルーティングテーブルを表示

```shell
route -6
```

## ip

ルーティング、ネットワークデバイス、インターフェースを管理します。

```shell
ubuntu@vultr:~$ whatis ip
ip (8)               - show / manipulate routing, network devices, interfaces and tunnels
```

ipコマンドは、iproute2というパッケージに含まれています。

iproute2は、ifconfig, route, arp, netstatなどのコマンドを置き換えるために開発されました。

オブジェクト(機能)に対してコマンドを指定します。

### 主なオブジェクト

-   link, li, l: ネットワークデバイス
-   addr, ad, a: ネットワークデバイスのIPアドレス
-   route, rou, r: ルーティングテーブルのエントリー
-   neighbour, nei, n: IPv4のARPキャッシュ、IPv6のNDキャッシュ

### 主なコマンド

-   add, ad, a: IPアドレスを追加
-   del, de, d: IPアドレスを削除
-   show, sh, s: IPアドレスを表示

### ネットワークインターフェースのIPアドレスを表示

```shell
ip addr show
ip a s
```

### ネットワークインターフェースの状況のみを表示

```shell
ip link show
ip li sh
```

## nslookup

ネームサーバに対して対話的にクエリを実行します。

```shell
ubuntu@vultr:~$ whatis nslookup
nslookup (1)         - query Internet name servers interactively
```

### 正引きを行う

ホスト名に対応するIPアドレスを問い合わせることを正引きと言います。

```shell
nslookup example.com
```

### 逆引きを行う

IPアドレスに対応するホスト名を問い合わせることを逆引きと言います。PTRレコードを参照します。

```shell
nslookup 172.217.31.174
```

## host

DNSルックアップを行います。nslookupと異なり、対話的な実行ができません。名前解決を行うコマンドの中で、最も表示される内容がシンプルです。

```shell
ubuntu@vultr:~$ whatis host
host (1)             - DNS lookup utility
```

### 正引きを行う

nslookupと同じく、引数にホスト名を渡します。

### 逆引きを行う

nslookupと同じく、引数にIPアドレスを渡します。

## dig

DNSルックアップを行います。nslookup, hostよりも詳細な情報を得られます。

```shell
ubuntu@vultr:~$ whatis dig
dig (1)              - DNS lookup utility
```

digコマンドでは、検索タイプを指定できます。

-   a: Aレコード
-   aaaa: AAAAレコード
-   any: 全てのレコード
-   cname: CNAMEレコード
-   mx: MXレコード
-   ns: NSレコード
-   ptr: PTRレコード
-   txt: TXTレコード

### aレコードを問い合わせる

```shell
dig a example.com
```

### NSレコードを問い合わせる

```shell
dig ns example.com
```

@以降にIPアドレスを指定することで、名前解決に使うネームサーバを指定できます。

例えば、CloudflareのパブリックDNSを使う場合は@1.1.1.1と指定します。

```shell
dig @1.1.1.1 example.com
```

## ifconfig

ネットワークインターフェースを設定します。

```shell
ubuntu@vultr:~$ whatis ifconfig
ifconfig (8)         - configure a network interface
```

### 全てのインターフェースの情報を表示

```shell
ifconfig
```

### 一部のインターフェースの情報を表示

```shell
ifconfig enp1s0
```

### インターフェースを有効にする

```shell
ifconfig enp1s0 up 
```

### インターフェースを無効にする

```shell
ifconfig enp1s0 down
```

## traceroute

リモートホストまでの経路情報を表示します。

```shell
ubuntu@vultr:~$ whatis traceroute
tcptraceroute (8)    - print the route packets trace to network host
```

### example.comまでの経路情報を表示

```shell
traceroute example.com
```

## netstat

ネットワークの統計情報を表示します。

```shell
ubuntu@vultr:~$ whatis netstat
netstat (8)          - Print network connections, routing tables, interface statistics, masquerade connections, and multicast memberships
```

### 全ての有効なネットワーク接続を表示

```shell
netstat
netstat -a
```

### UDPの有効なネットワーク接続を表示

```shell
netstat -u
```

### TCPの有効なネットワーク接続を表示

```shell
netstat -t
```

### ルーティングテーブルを表示

```shell
netstat -r
```

## netplan

yamlファイルを使ってネットワーク設定をします。

```shell
ubuntu@vultr:~$ whatis netplan
netplan (5)          - YAML network configuration abstraction for various backends
```

`/etc/netplan/`にyamlファイルを作成することで、そのファイルの内容を基にネットワークを設定してくれます。

例えば、enp1s0にDHCPを設定する場合は以下の内容を記載します。

```yaml
network
    version: 2
    renderer: networkd
    ethernets:
        enp1s0:
            dhcp4: true
            dhcp6: true
```

# ネットワークを管理するデーモン

## systemd-networkd

ネットワーク設定を管理するデーモンです。`/etc/systemd/network/`にあるファイルを参照します。

例えば、DHCPを設定する場合は`/etc/systemd/network/dhcp.network`というファイルに以下の設定を記載します。

```
[Match]
Name=en*

[Network]
DHCP=yes
```

## systemd-resolved

スタブリゾルバを使った名前解決を提供するデーモンです。DNSSEC, DNS over TLSにも対応しています。

`/etc/systemd/resolved.conf`に設定を記載します。

例:

```
[Resolve]
# Some examples of DNS servers which may be used for DNS= and FallbackDNS=:
# Cloudflare: 1.1.1.1#cloudflare-dns.com 1.0.0.1#cloudflare-dns.com 2606:4700:4700::1111#cloudflare-dns.com 2606:4700:4700::1001#cloudflare-dns.com
# Google:     8.8.8.8#dns.google 8.8.4.4#dns.google 2001:4860:4860::8888#dns.google 2001:4860:4860::8844#dns.google
# Quad9:      9.9.9.9#dns.quad9.net 149.112.112.112#dns.quad9.net 2620:fe::fe#dns.quad9.net 2620:fe::9#dns.quad9.net
#DNS=
#FallbackDNS=
#Domains=
#DNSSEC=no
#DNSOverTLS=no
#MulticastDNS=no
#LLMNR=no
#Cache=no-negative
#CacheFromLocalhost=no
#DNSStubListener=yes
#DNSStubListenerExtra=
#ReadEtcHosts=yes
#ResolveUnicastSingleLabel=no
```