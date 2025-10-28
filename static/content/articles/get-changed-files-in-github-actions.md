---
title: GitHub Actionsを用いて変更されたファイルを取得する
image: images/Microsoft-Fluentui-Emoji-3d-Oden-3d.1024.png
publishedAt: 2025-10-28
updatedAt: 2025-10-28
---
## はじめに
GitHub Actionsでワークフローを実行する際、どのファイルが変更されたかを検出し、その情報に基づいて処理を分岐させたい場合があります。例えば、ドキュメントファイルのみが変更された場合はテストをスキップする、特定のディレクトリに変更があった場合のみデプロイを実行する、といったケースです。
この記事では、GitHub Actionsをを用いてコミット間で変更されたファイルの一覧を取得する方法を紹介します。
## プラグインの利用
[changed-files](https://github.com/tj-actions/changed-files)というActionを使うことで、変更されたファイルやディレクトリを取得できます。このActionは、GitHubのREST APIやGitのdiffコマンドを使用して、変更されたファイルを検出しているようです。
### 実装例
プルリクエストにおいてマークダウンファイルが変更されたときに、どのファイルが更新されたのかを確認したい場合があります。具体的には、ドキュメントの更新確認、翻訳、デプロイといった作業がこれに当たります。この状況を想定したワークフローを以下に示します。

```yaml
name: Track Markdown Changes

on:
  pull_request:
    types: [opened, synchronize]
    branches:
      - main

jobs:
  detect-changes:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v5

      - name: Get changed markdown files
        id: changed-files
        uses: tj-actions/changed-files@v47
        with:
          files: |
            **.md

      - name: Display changed files
        if: steps.changed-files.outputs.any_changed == 'true'
        env:
          CHANGED_FILES: ${{ steps.changed-files.outputs.all_changed_files }}
        run: |
          echo "変更されたMarkdownファイル:"
          for file in ${CHANGED_FILES}; do
            echo "  - $file"
          done

      - name: No changes detected
        if: steps.changed-files.outputs.any_changed == 'false'
        run: echo "Markdownファイルの変更はありませんでした"
```
デモ: https://github.com/nagutabby/changed-files-demo/pull/1

このワークフローの処理の流れを以下に示します。
1. baseブランチがmainであるプルリクエストにおいて、プルリクエストが作成されたとき、またはプルリクエストのheadブランチに新しいコミットがプッシュされたときにワークフローが実行されます。
2. `actions/checkout`を使用して、リポジトリのコードを取得します。
3. `tj-actions/changed-files`を使用して、プルリクエストのbaseブランチの最新のコミットとheadブランチの最新のコミットを比較し、`.md`で終わるファイルの変更を検出します。
4. 変更が検出された場合は、変更された全てのマークダウンファイル名を表示します。変更が検出されなかった場合は、「Markdownファイルの変更はありませんでした」というメッセージを表示します。