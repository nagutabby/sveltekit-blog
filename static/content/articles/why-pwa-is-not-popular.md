---
title: PWA (Progressive Web Apps) が流行らない理由
image: images/Microsoft-Fluentui-Emoji-3d-End-Arrow-3d.1024.png
publishedAt: 2024-02-17
updatedAt: 2024-05-02
---
# PWA使ってますか？

PWA (Progressive Web Apps) とは、モバイルネイティブアプリの仕組みをWebアプリで利用できるようにする技術です。PWAによって利用できるようになる機能の例としては、オフラインアクセス、仮想ストレージ、ホーム画面への追加、ユーザーへのプッシュ通知があります。PWAのいいところは1種類のコードでAndroid、iOSという2つのプラットフォームに対応できる点です。PWAの現状を整理する前に、PWAを構成する技術を見ていきましょう。

# PWAの関連技術

## Service Worker

Webワーカーの一種であるService Workerは、ネットワークリクエストをインターセプトし、リソースのキャッシュを利用して任意のネットワークレスポンスを返すスクリプトです。Service Workerはページ読み込みを高速化するだけでなく、Webアプリのオフラインアクセスも提供します。

## Web Application Manifest

Webアプリにメタデータを提供するJSONファイルです。PWAをインストールする際に表示する説明文やスクリーンショット、PWAをホーム画面に追加したときのアプリのアイコンや名前などの情報を提供します。

## App Shell

WebアプリのUI (= シェル) から記事などのコンテンツを分離するためのアーキテクチャです。App Shellの導入により、シェルがキャッシュされるため、再アクセス時にページを瞬時に読み込めるようになります。

# PWAの悲しい現状

PWAはWebアプリ開発者にとって素晴らしい概念であるのは間違いありません。JavaScriptという馴染みのある言語でモバイルネイティブアプリと同等の機能を持つWebアプリを構築できるからです。しかし、それはPWAの生みの親であるGoogleとWebアプリ開発者が夢見る理想に過ぎません。

例えば、Appleは[EUでiOS 17.4を利用するユーザーはPWAが利用できなくなる](https://gigazine.net/news/20240216-ios-17-4-removes-web-apps-pwa/)と明らかにしました。

追記: iOS 17.4で予定されていた、[EUにおけるPWAの利用制限は撤回されました](https://gigazine.net/news/20240302-apple-pwa-support-not-remove/)。

他方では、Mozillaは[Standards Positions](https://github.com/mozilla/standards-positions)においてPWAに対するセキュリティ上の懸念を表明しています。次は、これらの課題を考察していきます。

## Appleの意見: PWAは既存のエコシステムを破壊する

Appleの意見を知るには、[WebKitのStandards Positions](https://github.com/WebKit/standards-positions)を読むのが手っ取り早いです。WebKitの開発チームの中にはAppleの社員もいるため、Standards Positionsでの議論を追うことで彼らの意見を把握できます。

Appleは、[Screen Wake Lock APIのサポートに消極的](https://github.com/WebKit/standards-positions/issues/19)です。Screen Wake Lock APIは、画面の明るさを変更したり、端末の画面をロックしたりする方法を提供するものです。Appleの開発者は、これは[ユーザーにとって煩わしい機能であると主張](https://github.com/WebKit/standards-positions/issues/19#issuecomment-1192089037)しています。しかし、多くの開発者からAppleの立場に対する反対意見が寄せられています。

別の例としては、[Web Application Manifestのscreenshotsメンバーのサポート](https://github.com/WebKit/standards-positions/issues/49)があります。screeenshotsメンバーは、PWAをインストールする際にアプリのスクリーンショットの情報を提供するものです。Appleの開発者は、これは[信頼できないアプリの情報を信頼できるUIに表示するため、ユーザーを誤解される恐れがあると主張](https://github.com/WebKit/standards-positions/issues/49#issuecomment-1232349214)しています。しかし、ユーザーはPWAをインストールする前に、必ずWebブラウザー上でそのアプリを読み込むため、PWAをインストールしようとしている時点でユーザーはそのアプリを信頼しています。したがってAppleの主張は不適切です。

このように、Appleは、ときどき尤もらしい理由を付けてPWAのサポートに反対しています。なぜAppleは説得力に欠ける方法でPWAのサポートに反対するのでしょうか？私は、「Appleが自社のモバイルネイティブプラットフォームのエコシステムを守ろうとしているから」だと思います。実際、AppleはPWAやサイドローディングに代表される自由なプラットフォームへの移行手段に反対しています。

## Mozillaの意見: PWAはセキュリティを担保できない

Mozillaの意見に関しても、同様に[MozillaのStandards Positions](https://github.com/mozilla/standards-positions)を読んで分析しました。

中でも、[Web Bluetooth APIのサポートに関する議論](https://github.com/mozilla/standards-positions/issues/95)は興味深いです。ある開発者は、[Web USB, Web Bluetooth, Web NFC, Web MIDIはいずれも自身の挙動をユーザーに分かりやすく伝えられないと主張](https://github.com/mozilla/standards-positions/issues/95#issuecomment-644962468)しています。PWAが行おうとしている処理の許可をユーザーに求める手段としては、許可プロンプトがあります。しかし、許可プロンプトは単にその処理を許可するか、あるいは拒否するという2つの選択肢しか提供できないため、ユーザーに処理内容を具体的に説明する方法はまだ用意されていません。しかも、Webアプリの公開にはApp StoreやGoogle Play Storeで行われているような審査の仕組みがないため、仮にWebアプリが行おうとしている処理を説明するための仕組みが提供されたとしても、Webアプリの開発者がその説明を偽装する可能性は残ります。

つまり、モバイルネイティブアプリで利用できる機能のうち、正確で詳細な説明を提供する必要があるものをPWAでサポートするには、セキュリティ上の懸念があり、仮にモバイルネイティブアプリのあらゆる機能をPWAに取り入れてしまうとセキュリティを担保できないというのがMozillaの主張です。

# まとめ

PWAはクロスプラットフォーム開発における1つの選択肢ですが、Safariを開発しているAppleと、Firefoxを開発しているMozillaはPWAのサポートに慎重であり、PWAのサポート状況はWebブラウザー間で隔たりがあります。モバイルネイティブアプリと同等の機能を持つWebアプリを提供するというWebアプリ開発者の夢はしばらく叶いそうにありません。

# 余談: あの人が反応してた

逆張りおじさんこと[DHHがiOS 17.4でのPWAサポートの一部廃止に反応](https://twitter.com/dhh/status/1755974833375211752)していました。DHHがAppleに批判的だとは知りませんでした。以前に、HEYの収益のうち15%から30%をAppleに支払えと言われたことがあるようで、[当時のポスト](https://x.com/dhh/status/1272968382329942017)を見る限り、Appleを相当嫌っているようです。
