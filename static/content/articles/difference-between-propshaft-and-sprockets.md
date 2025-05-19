---
title: 'PropshaftとSprocketsの違い: なぜPropshaftなのか'
image: images/Microsoft-Fluentui-Emoji-3d-Gem-Stone-3d.1024.png
publishedAt: 2023-08-26T00:00:00.000Z
updatedAt: 2024-05-01T00:00:00.000Z
---
# Sprockets

## コンセプト

HTTP/2以前を主なターゲットとしているアセットパイプラインです。

## 主な機能

-   アセット(静的ファイル)のパスの解決
-   アセットの連結
    -   Webページをレンダリングする際のリクエストの数を減らす
-   CoffeeScript, Sass/SCSSのトランスパイル
-   アセットファイルの縮小, 圧縮
-   ダイジェストの付与
-   開発サーバーの提供

# Propshaft

## コンセプト

HTTP/2とES6を利用することを想定したアセットパイプラインです。

HTTP/2の機能の1つに、単一コネクションにおけるレスポンスの並列化があります。これを利用することで、複数のアセットファイルを連結する必要がなくなります。

また、ES6に準拠したプログラムをそのまま利用するので、ES5に準拠したプログラムへのトランスパイルが不要になります。

## 主な機能

-   アセット(静的ファイル)のパスの解決
-   ダイジェストの付与
-   開発サーバーの提供
-   CSSファイルに書かれているアセットのURLのコンパイル
    -   アセットのURLにダイジェストを付与する

# Propshaftを使用する利点

ソフトウェアを構成するプログラムが少ないので、以下の利点があります。

-   Sprocketよりもアセットパイプラインのパフォーマンスが高い
-   バグによってソフトウェアが動作しなくなる可能性が比較的低い
-   開発者が使い方を理解しやすい

# Propshaftのリポジトリの中身を調べてみた

参照: [https://github.com/rails/propshaft](https://github.com/rails/propshaft)

## lib/propshaft/compiler/css\_asset\_urls.rb

CSSに書かれているアセットのURLをコンパイルします。

```ruby
# frozen_string_literal: true

require "propshaft/compiler"

class Propshaft::Compiler::CssAssetUrls < Propshaft::Compiler
  ASSET_URL_PATTERN = /url\(\s*["']?(?!(?:\#| %23 |data| http |\/\/))([^"'\s?#)]+)([#?][^"')]+)?\s*["']?\)/

  def compile(logical_path, input)
    input.gsub(ASSET_URL_PATTERN) { asset_url resolve_path(logical_path.dirname, $1), logical_path, $2, $1 }
  end

  private
    def resolve_path(directory, filename)
      if filename.start_with?("../")
        Pathname.new(directory + filename).relative_path_from("").to_s
      elsif filename.start_with?("/")
        filename.delete_prefix("/").to_s
      else
        (directory + filename.delete_prefix("./")).to_s
      end
    end

    def asset_url(resolved_path, logical_path, fingerprint, pattern)
      if asset = assembly.load_path.find(resolved_path)
        %[url("#{url_prefix}/#{asset.digested_path}#{fingerprint}")]
      else
        Propshaft.logger.warn "Unable to resolve '#{pattern}' for missing asset '#{resolved_path}' in #{logical_path}"
        %[url("#{pattern}")]
      end
    end
end
```

## lib/propshaft/compiler/source\_mapping\_urls.rb

ソースマップのURLをコメントとして付与します。

```ruby
# frozen_string_literal: true

require "propshaft/compiler"

class Propshaft::Compiler::SourceMappingUrls < Propshaft::Compiler
  SOURCE_MAPPING_PATTERN = %r{^(//|/\*)# sourceMappingURL=(.+\.map)}

  def compile(logical_path, input)
    input.gsub(SOURCE_MAPPING_PATTERN) { source_mapping_url(asset_path($2, logical_path), $1) }
  end

  private
    def asset_path(source_mapping_url, logical_path)
      if logical_path.dirname.to_s == "."
        source_mapping_url
      else
        logical_path.dirname.join(source_mapping_url).to_s
      end
    end

    def source_mapping_url(resolved_path, comment)
      if asset = assembly.load_path.find(resolved_path)
        "#{comment}# sourceMappingURL=#{url_prefix}/#{asset.digested_path}"
      else
        Propshaft.logger.warn "Removed sourceMappingURL comment for missing asset '#{resolved_path}' from #{resolved_path}"
        comment
      end
    end
end
```

## lib/propshaft/railties/assets.rake

rakeタスクを定義します。

## lib/propshaft/resolver/dynamic.rb

load\_pathからアセットを検索して、ファイルにダイジェストを付与します。

```ruby
module Propshaft::Resolver
  class Dynamic
    attr_reader :load_path, :prefix

    def initialize(load_path:, prefix:)
      @load_path, @prefix = load_path, prefix
    end

    def resolve(logical_path)
      if asset = load_path.find(logical_path)
        File.join prefix, asset.digested_path
      end
    end

    def read(logical_path)
      if asset = load_path.find(logical_path)
        asset.content
      end
    end
  end
end
```

## lib/propshaft/resolver/static.rb

manifestを元にアセットを検索します。

```ruby
module Propshaft::Resolver
  class Static
    attr_reader :manifest_path, :prefix

    def initialize(manifest_path:, prefix:)
      @manifest_path, @prefix = manifest_path, prefix
    end

    def resolve(logical_path)
      if asset_path = parsed_manifest[logical_path]
        File.join prefix, asset_path
      end
    end

    def read(logical_path)
      if asset_path = parsed_manifest[logical_path]
        manifest_path.dirname.join(asset_path).read
      end
    end

    private
      def parsed_manifest
        @parsed_manifest ||= JSON.parse(manifest_path.read, symbolize_names: false)
      end
  end
end
```

## lib/propshaft/assembly.rb

### resolver

manifestが存在していればそれを元にアセットを検索し、そうでなければload\_pathからアセットを検索します。

```ruby
def resolver
  @resolver ||= if manifest_path.exist?
    Propshaft::Resolver::Static.new manifest_path: manifest_path, prefix: config.prefix
  else
    Propshaft::Resolver::Dynamic.new load_path: load_path, prefix: config.prefix
  end
end
```

### reveal

パスの種類に応じてそれぞれのアセットに対してメソッドを呼び出します。

```ruby
def reveal(path_type = :logical_path)
  path_type = path_type.presence_in(%i[ logical_path path ]) || raise(ArgumentError, "Unknown path_type: #{path_type}")

  load_path.assets.collect do | asset |
    asset.send(path_type)
  end
end
```

## lib/propshaft/asset.rb

### digest

SHA1アルゴリズムを使用して16進数のダイジェストを付与します。

```ruby
def digest
  @digest ||= Digest::SHA1.hexdigest("#{content}#{version}")
end
```

## lib/propshaft/compiler.rb

コンパイラが行うことを定義します。compileメソッドを定義することで、それぞれのコンパイラがコンパイルを行うことを示しています。

```ruby
# Override this in a specific compiler
def compile(logical_path, input)
  raise NotImplementedError
end
```

## lib/propshaft/compilers.rb

### compilable?

特定の種類のアセットが存在しているかどうかを判定します。

```ruby
def compilable?(asset)
  registrations[asset.content_type.to_s].present?
end
```

### compile

アセットの内容を入力として受け取り、コンパイルします。

```ruby
def compile(asset)
  if relevant_registrations = registrations[asset.content_type.to_s]
    asset.content.dup.tap do | input |
      relevant_registrations.each do | compiler |
        input.replace compiler.new(assembly).compile(asset.logical_path, input)
      end
    end
  else
    asset.content
  end
end
```

## lib/propshaft/helper.rb

### compute\_asset\_path

アセットのパスを返します。

```ruby
def compute_asset_path(path, options = {})
  Rails.application.assets.resolver.resolve(path) || raise(MissingAssetError.new(path))
end
```

## lib/propshaft/load\_path.rb

### cache\_sweeper

JavaScriptファイルの変更を監視するオブジェクトを返します。

```ruby
def cache_sweeper
  @cache_sweeper ||= begin
    exts_to_watch  = Mime::EXTENSION_LOOKUP.map(&:first)
    files_to_watch = Array(paths).collect { | dir | [ dir.to_s, exts_to_watch ] }.to_h

    Rails.application.config.file_watcher.new([], files_to_watch) do
      clear_cache
    end
  end
end
```

### assets\_by\_path

アセットのパスとアセットの値をペアとするハッシュを返します。

```ruby
def assets_by_path
  @cached_assets_by_path || = Hash.new.tap do  |mapped|
    paths.each do | path |
      without_dotfiles(all_files_from_tree(path)).each do | file |
        logical_path = file.relative_path_from(path)
        mapped[logical_path.to_s] ||= Propshaft::Asset.new(file, logical_path: logical_path, version: version)
      end if path.exist?
    end
  end
end
```

## lib/propshaft/output\_path.rb

### files

それぞれのファイルの値として、ファイルのパス、ダイジェスト、最終更新日時を設定します。

```ruby
def files
  Hash.new.tap do | files |
    all_files_from_tree(path).each do | file |
      digested_path = file.relative_path_from(path)
      logical_path, digest = extract_path_and_digest(digested_path)

      files[digested_path.to_s] = {
        logical_path: logical_path.to_s,
        digest: digest,
        mtime: File.mtime(file)
      }
    end
  end
end
```

### fresh\_version\_within\_limit

ファイルの最終更新日時と更新期限を比較して、必要に応じてファイルを更新します。

```ruby
def fresh_version_within_limit(mtime, count, expires_at:, limit:)
  modified_at = [ 0, Time.now - mtime ].max
  modified_at < expires_at || limit < count
end
```

## lib/propshaft/processor.rb

### process

出力先のパスが存在するかどうかを確認してからmanifestを出力して、最後にアセットを出力します。

```ruby
def process
  ensure_output_path_exists
  write_manifest
  output_assets
end
```

### output\_assets

コンパイルできる場合はそのアセットをコンパイルし、コンパイルできない場合はそのアセットを出力先パスにコピーします。

```ruby
def output_assets
  load_path.assets.each do | asset |
    unless output_path.join(asset.digested_path).exist?
      Propshaft.logger.info "Writing #{asset.digested_path}"
      FileUtils.mkdir_p output_path.join(asset.digested_path.parent)
      output_asset(asset)
    end
  end
end
```

# まとめ

Propshaftはトランスパイルをしないので、ソフトウェアを構成するプログラムの量やファイルの数が少ないだろうとは思っていましたが、想像以上に少なかったです。有名なソフトウェアだけあってプログラムの振る舞いが分かりやすく記述されていたので、プログラムを眺めるだけでも学びがあるな～と改めて感じました。
