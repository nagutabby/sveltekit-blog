---
title: 'Windows, MacOS,  Linuxにおけるプロキシサーバーの設定'
image: images/Microsoft-Fluentui-Emoji-3d-Laptop-3d.1024.png
publishedAt: 2023-08-07T00:00:00.000Z
updatedAt: 2024-05-01T00:00:00.000Z
---

<h1 id="hdd801a1d1b">全体設定</h1><h2 id="h48ba0d2347">Windows</h2><h3 id="h5e0a70c55b">コマンドプロンプト</h3><p>プロキシサーバーを設定する</p><pre><code>set HTTP_PROXY=http[s]://&lt;FQDN&gt;:&lt;PORT&gt;
set HTTPS_PROXY=http[s]://&lt;FQDN&gt;:&lt;PORT&gt;
set FTP_PROXY=http[s]://&lt;FQDN&gt;:&lt;PORT&gt;</code></pre><p>プロキシサーバーの設定を表示する</p><pre><code>echo %HTTP_PROXY%
echo %HTTPS_PROXY%
echo %FTP_PROXY%</code></pre><p>プロキシサーバーの設定を削除する</p><pre><code>set HTTP_PROXY=
set HTTPS_PROXY=
set FTP_PROXY=</code></pre><h3 id="hc0162bd196">Powershell</h3><pre><code>$env:HTTP_PROXY=http[s]://&lt;FQDN&gt;:&lt;PORT&gt;
$env:HTTPS_PROXY=http[s]://&lt;FQDN&gt;:&lt;PORT&gt;
$env:FTP_PROXY=http[s]://&lt;FQDN&gt;:&lt;PORT&gt;</code></pre><p>プロキシサーバーの設定を表示する</p><pre><code>echo $env:HTTP_PROXY
echo $env:HTTPS_PROXY
echo $env:FTP_PROXY</code></pre><p>プロキシサーバーの設定を削除する</p><pre><code>$env:HTTP_PROXY=
$env:HTTPS_PROXY=
$env:FTP_PROXY=</code></pre><h2 id="h709251dcf3">MacOS, Linux</h2><p>プロキシサーバーを設定する</p><pre><code>HTTP_PROXY=http[s]://&lt;FQDN&gt;:&lt;PORT&gt;
HTTPS_PROXY=http[s]://&lt;FQDN&gt;:&lt;PORT&gt;
FTP_PROXY=http[s]://&lt;FQDN&gt;:&lt;PORT&gt;</code></pre><p>プロキシサーバーの設定を表示する</p><pre><code>echo $HTTP_PROXY
echo $HTTPS_PROXY
echo $FTP_PROXY</code></pre><p>プロキシサーバーの設定を削除する</p><pre><code>HTTP_PROXY=
HTTPS_PROXY=
FTP_PROXY=</code></pre><h1 id="hdd06fafc64">個別設定</h1><h3 id="h8f037f7362">apt</h3><p><code>/etc/apt/apt.conf</code>に以下の設定を書き込む</p><pre><code>Acquire::http::proxy http[s]://&lt;FQDN&gt;:&lt;PORT&gt;;
Acquire::https::proxy http[s]://&lt;FQDN&gt;:&lt;PORT&gt;;
Acquire::ftp::proxy http[s]://&lt;FQDN&gt;:&lt;PORT&gt;;</code></pre><h3 id="h62db7f7988">yum</h3><p><code>/etc/yum.conf</code>に以下の設定を書き込む</p><pre><code>proxy=http[s]://&lt;FQDN&gt;:&lt;PORT&gt;</code></pre><h3 id="hd1ef0f582f">dnf</h3><p><code>/etc/dnf/dnf.conf</code>のmainセクションに以下の設定を書き込む</p><pre><code>proxy=http[s]://&lt;FQDN&gt;:&lt;PORT&gt;</code></pre><h3 id="h7343673ad5">Git</h3><p>プロキシサーバーを設定する</p><pre><code>git config --global http.proxy http[s]://&lt;FQDN&gt;:&lt;PORT&gt;
git config --global https.proxy http[s]://&lt;FQDN&gt;:&lt;PORT&gt;</code></pre><p>プロキシサーバーの設定を表示する</p><pre><code>git config --global --list</code></pre><p>プロキシサーバーの設定を削除する</p><pre><code>git config --global --unset http.proxy
git config --global --unset https.proxy</code></pre><p></p>
