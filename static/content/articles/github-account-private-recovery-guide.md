---
title: GitHubのアカウントが非公開になった場合の対処法
image: images/Microsoft-Fluentui-Emoji-3d-Zzz-3d.1024.png
publishedAt: 2024-10-06
updatedAt: 2024-10-06
---
# 初めに

GitHubのアカウントが突然非公開になることがあります。[GitHub Transparency Centerが公開しているデータ](https://transparencycenter.github.com/appeals/#abuse-related-violations)によると、2023年にアカウントを非公開にする処罰のみを受けたアカウントの数は2.5万個以上、アカウントアクセスを制限する処罰のみを受けたアカウントの数は1万個以上、両方の処罰を受けたアカウントの数は1万個以上にもなります。このことから、アカウントに何らかの問題が生じるリスクは低いとはいえ、誰もがアカウントの問題に直面する可能性があると言えます。では、なぜこのような問題が起きるのでしょうか？

# アカウントが非公開になる原因

大抵の場合は、「GitHubの[利用規約](https://docs.github.com/en/site-policy/github-terms/github-terms-of-service)や[コミュニティガイドライン](https://docs.github.com/en/site-policy/github-terms/github-community-guidelines)に違反している」と**疑われている**ことが原因です。例えば、営利目的のWebサイトに誘導するために、自分のアカウントにそのWebサイトのURLを関連付けると、それはスパム行為であると見なされます。ここで重要なのは、GitHubのモデレーターが手動でアカウントを非公開にしたり、アカウントアクセスを制限したりしている訳ではないということです。GitHubサポートによれば、「スパムアカウントを見つけて処罰するための多くのプログラムが自動で実行されている」そうです。そのため、スパム検出プログラムが誤ってアカウントを非公開にすることがあります。実際に届いた返信を以下に示します。

> Hello,
>  
> Thank you for writing in to GitHub Support and apologies for any inconvenience.
>  
> We use a lot of scripts to find (and hide) spammy accounts. It looks like one of these scripts caused your account to be flagged by mistake.
>  
> I've removed that flag from your account now, so things ought to be back to normal the next time you sign in.
>  
> If there's anything else we might be able to do to help, don't hesitate to let us know.
>  
> Best regards,
> Felix

# アカウントの復旧方法

[Appeal and Reinstatement](https://support.github.com/contact/reinstatement)というページにある問い合わせフォームに情報を入力し、フォームを送信する。フォームを送信するとサポートチケットが作られるので、GitHubサポートからの返信が来るまで待ちます。返信までは3営業日ほどかかります。アカウントが非公開になった理由によっては、追加の対応が必要になる場合があります。ユーザーによる対応が必要な場合は、何をすべきかについてGitHubサポートから連絡があるので、その内容に従いましょう。対応が完了したらその旨を連絡し、GitHubサポートがアカウントを公開にするまで待ちます。

# 今後の対策

GitHubサポートから指摘された問題が今後も発生するのを防ぐために、利用規約やコミュニティガイドラインを熟読しましょう。また、[GitHubアカウントの問題について議論しているGitHub Communityのページ](https://github.com/orgs/community/discussions/27294)を参考にすると違反事例を把握できます。例えグレーゾーンであっても制限の対象になるので、GitHubのガイドラインを理解するように努めることが大切です。
