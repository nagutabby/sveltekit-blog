---
title: 'PropshaftとSprocketsの違い: なぜPropshaftなのか'
image: images/Microsoft-Fluentui-Emoji-3d-Gem-Stone-3d.1024.png
publishedAt: 2023-08-26T00:00:00.000Z
updatedAt: 2024-05-01T00:00:00.000Z
---

<h1 id="hca38cfafb0">Sprockets</h1><h2 id="h0544ccba1b">コンセプト</h2><p>HTTP/2以前を主なターゲットとしているアセットパイプラインです。</p><h2 id="he3bf764d36">主な機能</h2><ul><li>アセット(静的ファイル)のパスの解決</li><li>アセットの連結<ul><li>Webページをレンダリングする際のリクエストの数を減らす</li></ul></li><li>CoffeeScript, Sass/SCSSのトランスパイル</li><li>アセットファイルの縮小, 圧縮</li><li>ダイジェストの付与</li><li>開発サーバーの提供</li></ul><h1 id="h7caff2d77f">Propshaft</h1><h2 id="h0544ccba1b">コンセプト</h2><p>HTTP/2とES6を利用することを想定したアセットパイプラインです。</p><p>HTTP/2の機能の1つに、単一コネクションにおけるレスポンスの並列化があります。これを利用することで、複数のアセットファイルを連結する必要がなくなります。</p><p>また、ES6に準拠したプログラムをそのまま利用するので、ES5に準拠したプログラムへのトランスパイルが不要になります。</p><h2 id="he3bf764d36">主な機能</h2><ul><li>アセット(静的ファイル)のパスの解決</li><li>ダイジェストの付与</li><li>開発サーバーの提供</li><li>CSSファイルに書かれているアセットのURLのコンパイル<ul><li>アセットのURLにダイジェストを付与する</li></ul></li></ul><h1 id="h7f8e05668e">Propshaftを使用する利点</h1><p>ソフトウェアを構成するプログラムが少ないので、以下の利点があります。</p><ul><li>Sprocketよりもアセットパイプラインのパフォーマンスが高い</li><li>バグによってソフトウェアが動作しなくなる可能性が比較的低い</li><li>開発者が使い方を理解しやすい</li></ul><h1 id="hd16f02a732">Propshaftのリポジトリの中身を調べてみた</h1><p>参照: <a href="https://github.com/rails/propshaft">https://github.com/rails/propshaft</a></p><h2 id="h4effb31774">lib/propshaft/compiler/css_asset_urls.rb</h2><p>CSSに書かれているアセットのURLをコンパイルします。</p><pre><code># frozen_string_literal: true

require &quot;propshaft/compiler&quot;

class Propshaft::Compiler::CssAssetUrls &lt; Propshaft::Compiler
  ASSET_URL_PATTERN = /url\(\s*[&quot;&apos;]?(?!(?:\#|%23|data|http|\/\/))([^&quot;&apos;\s?#)]+)([#?][^&quot;&apos;)]+)?\s*[&quot;&apos;]?\)/

  def compile(logical_path, input)
    input.gsub(ASSET_URL_PATTERN) { asset_url resolve_path(logical_path.dirname, $1), logical_path, $2, $1 }
  end

  private
    def resolve_path(directory, filename)
      if filename.start_with?(&quot;../&quot;)
        Pathname.new(directory + filename).relative_path_from(&quot;&quot;).to_s
      elsif filename.start_with?(&quot;/&quot;)
        filename.delete_prefix(&quot;/&quot;).to_s
      else
        (directory + filename.delete_prefix(&quot;./&quot;)).to_s
      end
    end

    def asset_url(resolved_path, logical_path, fingerprint, pattern)
      if asset = assembly.load_path.find(resolved_path)
        %[url(&quot;#{url_prefix}/#{asset.digested_path}#{fingerprint}&quot;)]
      else
        Propshaft.logger.warn &quot;Unable to resolve &apos;#{pattern}&apos; for missing asset &apos;#{resolved_path}&apos; in #{logical_path}&quot;
        %[url(&quot;#{pattern}&quot;)]
      end
    end
end</code></pre><h2 id="h10259288f1">lib/propshaft/compiler/source_mapping_urls.rb</h2><p>ソースマップのURLをコメントとして付与します。</p><pre><code># frozen_string_literal: true

require &quot;propshaft/compiler&quot;

class Propshaft::Compiler::SourceMappingUrls &lt; Propshaft::Compiler
  SOURCE_MAPPING_PATTERN = %r{^(//|/\*)# sourceMappingURL=(.+\.map)}

  def compile(logical_path, input)
    input.gsub(SOURCE_MAPPING_PATTERN) { source_mapping_url(asset_path($2, logical_path), $1) }
  end

  private
    def asset_path(source_mapping_url, logical_path)
      if logical_path.dirname.to_s == &quot;.&quot;
        source_mapping_url
      else
        logical_path.dirname.join(source_mapping_url).to_s
      end
    end

    def source_mapping_url(resolved_path, comment)
      if asset = assembly.load_path.find(resolved_path)
        &quot;#{comment}# sourceMappingURL=#{url_prefix}/#{asset.digested_path}&quot;
      else
        Propshaft.logger.warn &quot;Removed sourceMappingURL comment for missing asset &apos;#{resolved_path}&apos; from #{resolved_path}&quot;
        comment
      end
    end
end</code></pre><h2 id="hc3f3c18909">lib/propshaft/railties/assets.rake</h2><p>rakeタスクを定義します。</p><h2 id="hd336a0774b">lib/propshaft/resolver/dynamic.rb</h2><p>load_pathからアセットを検索して、ファイルにダイジェストを付与します。</p><pre><code>module Propshaft::Resolver
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
end</code></pre><h2 id="h450f404220">lib/propshaft/resolver/static.rb</h2><p>manifestを元にアセットを検索します。</p><pre><code>module Propshaft::Resolver
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
end</code></pre><h2 id="h0f559948c4">lib/propshaft/assembly.rb</h2><h3 id="hd3306b3325">resolver</h3><p>manifestが存在していればそれを元にアセットを検索し、そうでなければload_pathからアセットを検索します。</p><pre><code>def resolver
  @resolver ||= if manifest_path.exist?
    Propshaft::Resolver::Static.new manifest_path: manifest_path, prefix: config.prefix
  else
    Propshaft::Resolver::Dynamic.new load_path: load_path, prefix: config.prefix
  end
end</code></pre><h3 id="h577b6bbd2e">reveal</h3><p>パスの種類に応じてそれぞれのアセットに対してメソッドを呼び出します。</p><pre><code>def reveal(path_type = :logical_path)
  path_type = path_type.presence_in(%i[ logical_path path ]) || raise(ArgumentError, &quot;Unknown path_type: #{path_type}&quot;)
    
  load_path.assets.collect do |asset|
    asset.send(path_type)
  end
end</code></pre><h2 id="h2af09b3801">lib/propshaft/asset.rb</h2><h3 id="hb7dd64435d">digest</h3><p>SHA1アルゴリズムを使用して16進数のダイジェストを付与します。</p><pre><code>def digest
  @digest ||= Digest::SHA1.hexdigest(&quot;#{content}#{version}&quot;)
end</code></pre><h2 id="h7ec20ff65e">lib/propshaft/compiler.rb</h2><p>コンパイラが行うことを定義します。compileメソッドを定義することで、それぞれのコンパイラがコンパイルを行うことを示しています。</p><pre><code># Override this in a specific compiler
def compile(logical_path, input)
  raise NotImplementedError
end</code></pre><h2 id="hb2caa82440">lib/propshaft/compilers.rb</h2><h3 id="hf3c19c9154">compilable?</h3><p>特定の種類のアセットが存在しているかどうかを判定します。</p><pre><code>def compilable?(asset)
  registrations[asset.content_type.to_s].present?
end</code></pre><h3 id="hd1a08bfaa6">compile</h3><p>アセットの内容を入力として受け取り、コンパイルします。</p><pre><code>def compile(asset)
  if relevant_registrations = registrations[asset.content_type.to_s]
    asset.content.dup.tap do |input|
      relevant_registrations.each do |compiler|
        input.replace compiler.new(assembly).compile(asset.logical_path, input)
      end
    end
  else
    asset.content
  end
end</code></pre><h2 id="hb6535fa800">lib/propshaft/helper.rb</h2><h3 id="h8f2cd7cdca">compute_asset_path</h3><p>アセットのパスを返します。</p><pre><code>def compute_asset_path(path, options = {})
  Rails.application.assets.resolver.resolve(path) || raise(MissingAssetError.new(path))
end</code></pre><h2 id="h5f0742a71b">lib/propshaft/load_path.rb</h2><h3 id="haffa70685e">cache_sweeper</h3><p>JavaScriptファイルの変更を監視するオブジェクトを返します。</p><pre><code>def cache_sweeper
  @cache_sweeper ||= begin
    exts_to_watch  = Mime::EXTENSION_LOOKUP.map(&amp;:first)
    files_to_watch = Array(paths).collect { |dir| [ dir.to_s, exts_to_watch ] }.to_h

    Rails.application.config.file_watcher.new([], files_to_watch) do
      clear_cache
    end
  end
end</code></pre><h3 id="hfaeae3c7cf">assets_by_path</h3><p>アセットのパスとアセットの値をペアとするハッシュを返します。</p><pre><code>def assets_by_path
  @cached_assets_by_path ||= Hash.new.tap do |mapped|
    paths.each do |path|
      without_dotfiles(all_files_from_tree(path)).each do |file|
        logical_path = file.relative_path_from(path)
        mapped[logical_path.to_s] ||= Propshaft::Asset.new(file, logical_path: logical_path, version: version)
      end if path.exist?
    end
  end
end</code></pre><h2 id="h9c748402fa">lib/propshaft/output_path.rb</h2><h3 id="h7bb3554771">files </h3><p>それぞれのファイルの値として、ファイルのパス、ダイジェスト、最終更新日時を設定します。</p><pre><code>def files
  Hash.new.tap do |files|
    all_files_from_tree(path).each do |file|
      digested_path = file.relative_path_from(path)
      logical_path, digest = extract_path_and_digest(digested_path)

      files[digested_path.to_s] = {
        logical_path: logical_path.to_s,
        digest: digest,
        mtime: File.mtime(file)
      }
    end
  end
end</code></pre><h3 id="hb3373c81c6">fresh_version_within_limit</h3><p>ファイルの最終更新日時と更新期限を比較して、必要に応じてファイルを更新します。</p><pre><code>def fresh_version_within_limit(mtime, count, expires_at:, limit:)
  modified_at = [ 0, Time.now - mtime ].max
  modified_at &lt; expires_at || limit &lt; count
end</code></pre><h2 id="hf352d86a66">lib/propshaft/processor.rb</h2><h3 id="h21296380ce">process</h3><p>出力先のパスが存在するかどうかを確認してからmanifestを出力して、最後にアセットを出力します。</p><pre><code>def process
  ensure_output_path_exists
  write_manifest
  output_assets
end</code></pre><h3 id="h8ec1c3ae3b">output_assets</h3><p>コンパイルできる場合はそのアセットをコンパイルし、コンパイルできない場合はそのアセットを出力先パスにコピーします。</p><pre><code>def output_assets
  load_path.assets.each do |asset|
    unless output_path.join(asset.digested_path).exist?
      Propshaft.logger.info &quot;Writing #{asset.digested_path}&quot;
      FileUtils.mkdir_p output_path.join(asset.digested_path.parent)
      output_asset(asset)
    end
  end
end</code></pre><h1 id="ha214098e44">まとめ</h1><p>Propshaftはトランスパイルをしないので、ソフトウェアを構成するプログラムの量やファイルの数が少ないだろうとは思っていましたが、想像以上に少なかったです。有名なソフトウェアだけあってプログラムの振る舞いが分かりやすく記述されていたので、プログラムを眺めるだけでも学びがあるな～と改めて感じました。</p>
