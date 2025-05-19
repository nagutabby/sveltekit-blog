---
title: 'Windows, MacOS,  Linuxにおけるプロキシサーバーの設定'
image: images/Microsoft-Fluentui-Emoji-3d-Laptop-3d.1024.png
publishedAt: 2023-08-07T00:00:00.000Z
updatedAt: 2024-05-01T00:00:00.000Z
---
# 全体設定

## Windows

### コマンドプロンプト

プロキシサーバーを設定する

```
set HTTP_PROXY=http[s]://<FQDN>:<PORT>
set HTTPS_PROXY=http[s]://<FQDN>:<PORT>
set FTP_PROXY=http[s]://<FQDN>:<PORT>
```

プロキシサーバーの設定を表示する

```
echo %HTTP_PROXY%
echo %HTTPS_PROXY%
echo %FTP_PROXY%
```

プロキシサーバーの設定を削除する

```
set HTTP_PROXY=
set HTTPS_PROXY=
set FTP_PROXY=
```

### Powershell

```
$env:HTTP_PROXY=http[s]://<FQDN>:<PORT>
$env:HTTPS_PROXY=http[s]://<FQDN>:<PORT>
$env:FTP_PROXY=http[s]://<FQDN>:<PORT>
```

プロキシサーバーの設定を表示する

```
echo $env:HTTP_PROXY
echo $env:HTTPS_PROXY
echo $env:FTP_PROXY
```

プロキシサーバーの設定を削除する

```
$env:HTTP_PROXY=
$env:HTTPS_PROXY=
$env:FTP_PROXY=
```

## MacOS, Linux

プロキシサーバーを設定する

```
HTTP_PROXY=http[s]://<FQDN>:<PORT>
HTTPS_PROXY=http[s]://<FQDN>:<PORT>
FTP_PROXY=http[s]://<FQDN>:<PORT>
```

プロキシサーバーの設定を表示する

```
echo $HTTP_PROXY
echo $HTTPS_PROXY
echo $FTP_PROXY
```

プロキシサーバーの設定を削除する

```
HTTP_PROXY=
HTTPS_PROXY=
FTP_PROXY=
```

# 個別設定

### apt

`/etc/apt/apt.conf`に以下の設定を書き込む

```
Acquire::http::proxy http[s]://<FQDN>:<PORT>;
Acquire::https::proxy http[s]://<FQDN>:<PORT>;
Acquire::ftp::proxy http[s]://<FQDN>:<PORT>;
```

### yum

`/etc/yum.conf`に以下の設定を書き込む

```
proxy=http[s]://<FQDN>:<PORT>
```

### dnf

`/etc/dnf/dnf.conf`のmainセクションに以下の設定を書き込む

```
proxy=http[s]://<FQDN>:<PORT>
```

### Git

プロキシサーバーを設定する

```
git config --global http.proxy http[s]://<FQDN>:<PORT>
git config --global https.proxy http[s]://<FQDN>:<PORT>
```

プロキシサーバーの設定を表示する

```
git config --global --list
```

プロキシサーバーの設定を削除する

```
git config --global --unset http.proxy
git config --global --unset https.proxy
```