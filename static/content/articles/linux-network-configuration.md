---
title: "Linuxでネットワークを管理したい"
image: "https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/58069a149226412f83b74ab1c30a88b7/Microsoft-Fluentui-Emoji-3d-Penguin-3d.1024.png"
publishedAt: 2023-07-05
updatedAt: 2024-05-01
---

<h1 id="h097de7bf16">ネットワークを管理するファイル</h1><h2 id="hb1c2bb56a3">/etc/hosts</h2><p>名前解決をするために使われるファイルです。一行に1つのIPアドレスと1つ以上のホスト名を持ちます。</p><p>現在は、DNSによる名前解決と併用されています。</p><p>例:</p><pre><code class="language-shell">127.0.1.1 vultr.guest vultr
127.0.0.1 localhost

# The following lines are desirable for IPv6 capable hosts
::1 localhost ip6-localhost ip6-loopback
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters</code></pre><h2 id="h2aa7c08c5e">/etc/hostname</h2><p>ホスト名を設定します。Debian系で使われるファイルです。</p><p></p><p>例:</p><pre><code class="language-shell">vultr</code></pre><h2 id="hbd4664f19d">/etc/nsswitch.conf</h2><p>名前解決やサービス名の解決を行う際の問い合わせ順序を指定します。</p><p>例えば、<code>hosts</code>に<code>files dns</code>を指定している場合は、<code>/etc/hosts</code>を参照し、そのファイルにIPアドレスとホスト名の対応が記載されていなかった場合に、<code>/etc/resolv.conf</code>に記載されているDNSによる名前解決を行います。</p><p>例:</p><pre><code>passwd:         files systemd
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
</code></pre><h1 id="h493c11690d">ネットワークを管理するコマンド</h1><h2 id="h6094b0b3ec">route</h2><p>ルーティングテーブルを操作します。</p><pre><code class="language-shell">ubuntu@vultr:~$ whatis route
route (8)            - show / manipulate the IP routing table</code></pre><h3 id="he4710622c5">IPv4のルーティングテーブルを表示</h3><pre><code class="language-shell">route</code></pre><pre><code class="language-shell">route -4</code></pre><h3 id="h78c97c9575">IPv6のルーティングテーブルを表示</h3><pre><code class="language-shell">route -6</code></pre><h2 id="h774be72e9c">ip</h2><p>ルーティング、ネットワークデバイス、インターフェースを管理します。</p><pre><code class="language-shell">ubuntu@vultr:~$ whatis ip
ip (8)               - show / manipulate routing, network devices, interfaces and tunnels</code></pre><p>ipコマンドは、iproute2というパッケージに含まれています。</p><p>iproute2は、ifconfig, route, arp, netstatなどのコマンドを置き換えるために開発されました。</p><p>オブジェクト(機能)に対してコマンドを指定します。</p><h3 id="hcec9dd1f5f">主なオブジェクト</h3><ul><li>link, li, l: ネットワークデバイス</li><li>addr, ad, a: ネットワークデバイスのIPアドレス</li><li>route, rou, r: ルーティングテーブルのエントリー</li><li>neighbour, nei, n: IPv4のARPキャッシュ、IPv6のNDキャッシュ</li></ul><h3 id="h8a9f507066">主なコマンド</h3><ul><li>add, ad, a: IPアドレスを追加</li><li>del, de, d: IPアドレスを削除</li><li>show, sh, s: IPアドレスを表示</li></ul><h3 id="h9951201930">ネットワークインターフェースのIPアドレスを表示</h3><pre><code class="language-shell">ip addr show
ip a s</code></pre><p></p><h3 id="h402dc75cc1">ネットワークインターフェースの状況のみを表示</h3><pre><code class="language-shell">ip link show
ip li sh</code></pre><h2 id="h22f236ce0e">nslookup</h2><p>ネームサーバに対して対話的にクエリを実行します。</p><pre><code class="language-shell">ubuntu@vultr:~$ whatis nslookup
nslookup (1)         - query Internet name servers interactively</code></pre><h3 id="hc953a58cc4">正引きを行う</h3><p>ホスト名に対応するIPアドレスを問い合わせることを正引きと言います。</p><pre><code class="language-shell">nslookup example.com</code></pre><h3 id="hdb13e60e56">逆引きを行う</h3><p>IPアドレスに対応するホスト名を問い合わせることを逆引きと言います。PTRレコードを参照します。</p><pre><code class="language-shell">nslookup 172.217.31.174</code></pre><h2 id="h5236cf482a">host</h2><p>DNSルックアップを行います。nslookupと異なり、対話的な実行ができません。名前解決を行うコマンドの中で、最も表示される内容がシンプルです。</p><pre><code class="language-shell">ubuntu@vultr:~$ whatis host
host (1)             - DNS lookup utility</code></pre><h3 id="hc953a58cc4">正引きを行う</h3><p>nslookupと同じく、引数にホスト名を渡します。</p><h3 id="hdb13e60e56">逆引きを行う</h3><p>nslookupと同じく、引数にIPアドレスを渡します。</p><h2 id="ha7ee631e98">dig</h2><p>DNSルックアップを行います。nslookup, hostよりも詳細な情報を得られます。</p><pre><code class="language-shell">ubuntu@vultr:~$ whatis dig
dig (1)              - DNS lookup utility</code></pre><p>digコマンドでは、検索タイプを指定できます。</p><ul><li>a: Aレコード</li><li>aaaa: AAAAレコード</li><li>any: 全てのレコード</li><li>cname: CNAMEレコード</li><li>mx: MXレコード</li><li>ns: NSレコード</li><li>ptr: PTRレコード</li><li>txt: TXTレコード</li></ul><h3 id="hc6df653011">aレコードを問い合わせる</h3><pre><code class="language-shell">dig a example.com</code></pre><h3 id="h747ef66c14">NSレコードを問い合わせる</h3><pre><code class="language-shell">dig ns example.com</code></pre><p>@以降にIPアドレスを指定することで、名前解決に使うネームサーバを指定できます。</p><p>例えば、CloudflareのパブリックDNSを使う場合は@1.1.1.1と指定します。</p><pre><code class="language-shell">dig @1.1.1.1 example.com</code></pre><h2 id="h02938d1622">ifconfig</h2><p>ネットワークインターフェースを設定します。</p><pre><code class="language-shell">ubuntu@vultr:~$ whatis ifconfig
ifconfig (8)         - configure a network interface</code></pre><h3 id="h1fa6c1c697">全てのインターフェースの情報を表示</h3><pre><code class="language-shell">ifconfig</code></pre><h3 id="h6e93416670">一部のインターフェースの情報を表示</h3><pre><code class="language-shell">ifconfig enp1s0</code></pre><h3 id="hc8aaa1ec6a">インターフェースを有効にする</h3><pre><code class="language-shell">ifconfig enp1s0 up </code></pre><h3 id="hcc61dcbc1f">インターフェースを無効にする</h3><pre><code class="language-shell">ifconfig enp1s0 down</code></pre><h2 id="h0abad1822b">traceroute</h2><p>リモートホストまでの経路情報を表示します。</p><pre><code class="language-shell">ubuntu@vultr:~$ whatis traceroute
tcptraceroute (8)    - print the route packets trace to network host</code></pre><h3 id="h33aa146358">example.comまでの経路情報を表示</h3><pre><code class="language-shell">traceroute example.com</code></pre><h2 id="h851c962a77">netstat</h2><p>ネットワークの統計情報を表示します。</p><pre><code class="language-shell">ubuntu@vultr:~$ whatis netstat
netstat (8)          - Print network connections, routing tables, interface statistics, masquerade connections, and multicast memberships</code></pre><h3 id="hcfb66435e9">全ての有効なネットワーク接続を表示</h3><pre><code class="language-shell">netstat
netstat -a</code></pre><h3 id="ha71a49fdbc">UDPの有効なネットワーク接続を表示</h3><pre><code class="language-shell">netstat -u</code></pre><h3 id="he2a24d77da">TCPの有効なネットワーク接続を表示</h3><pre><code class="language-shell">netstat -t</code></pre><h3 id="h51d1f01022">ルーティングテーブルを表示</h3><pre><code class="language-shell">netstat -r</code></pre><h2 id="h9e05a265a8">netplan</h2><p>yamlファイルを使ってネットワーク設定をします。</p><pre><code class="language-shell">ubuntu@vultr:~$ whatis netplan
netplan (5)          - YAML network configuration abstraction for various backends</code></pre><p><code>/etc/netplan/</code>にyamlファイルを作成することで、そのファイルの内容を基にネットワークを設定してくれます。</p><p>例えば、enp1s0にDHCPを設定する場合は以下の内容を記載します。</p><pre><code class="language-yaml">network
    version: 2
    renderer: networkd
    ethernets:
        enp1s0:
            dhcp4: true
            dhcp6: true</code></pre><h1 id="h86e2e5bf36">ネットワークを管理するデーモン</h1><h2 id="h43a00dfc07">systemd-networkd</h2><p>ネットワーク設定を管理するデーモンです。<code>/etc/systemd/network/</code>にあるファイルを参照します。</p><p>例えば、DHCPを設定する場合は<code>/etc/systemd/network/dhcp.network</code>というファイルに以下の設定を記載します。</p><pre><code>[Match]
Name=en*

[Network]
DHCP=yes</code></pre><h2 id="h0eb3f8916e">systemd-resolved</h2><p>スタブリゾルバを使った名前解決を提供するデーモンです。DNSSEC, DNS over TLSにも対応しています。</p><p><code>/etc/systemd/resolved.conf</code>に設定を記載します。</p><p>例:</p><pre><code>[Resolve]
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
#ResolveUnicastSingleLabel=no</code></pre>