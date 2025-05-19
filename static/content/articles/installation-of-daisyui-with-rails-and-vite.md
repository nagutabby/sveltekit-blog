---
title: RailsとViteを使ったdaisyUIのインストール
image: images/Microsoft-Fluentui-Emoji-3d-Blossom-3d.1024.png
publishedAt: 2023-08-28T00:00:00.000Z
updatedAt: 2024-05-04T00:00:00.000Z
---

<h1 id="h4f59cac49a">このアプリの構成</h1><table><tbody><tr><td colspan="1" rowspan="1"><p>ビルドツール</p></td><td colspan="1" rowspan="1"><p>Vite</p></td></tr><tr><td colspan="1" rowspan="1"><p>データベース</p></td><td colspan="1" rowspan="1"><p>PostgreSQL</p></td></tr><tr><td colspan="1" rowspan="1"><p>CSSフレームワーク</p></td><td colspan="1" rowspan="1"><p>Tailwind CSS</p></td></tr><tr><td colspan="1" rowspan="1"><p>UIライブラリ</p></td><td colspan="1" rowspan="1"><p>daisyUI</p></td></tr><tr><td colspan="1" rowspan="1"><p>仮想環境</p></td><td colspan="1" rowspan="1"><p>Docker</p></td></tr><tr><td colspan="1" rowspan="1"><p>Node.jsのパッケージマネージャー</p></td><td colspan="1" rowspan="1"><p>pnpm</p></td></tr></tbody></table><h1 id="h8818f98a46">プロジェクトの準備</h1><h2 id="h903d18fb31">Dockerのインストール</h2><p><a href="https://www.docker.com/" target="_blank" rel="noopener noreferrer nofollow">Dockerの公式サイト</a>からDocker Desktopをダウンロードしてインストールします。</p><h2 id="hcc910410a6">プロジェクトディレクトリの作成</h2><pre><code class="language-bash">mkdir ~/rails-project
cd ~/rails-project</code></pre><h2 id="hf1701c5a32">Dockerfileの作成</h2><p>Railsのイメージをビルドするために<code>Dockerfile</code>を作ります。</p><pre><code class="language-dockerfile">FROM ruby:3.2.2

RUN apt-get update
RUN curl -sL http s://deb.nodesource.com/setup_18.x | bash -
RUN apt-get install -y nodejs npm
RUN npm install -g pnpm

ENV PROJECT_DIR /root/rails-project

RUN mkdir $PROJECT_DIR
WORKDIR $PROJECT_DIR

COPY Gemfile $PROJECT_DIR/Gemfile
COPY Gemfile.lock $PROJECT_DIR/Gemfile.lock

ENV BUNDLER_VERSION 2.4.19

RUN gem update --system
RUN gem install bundler -v $BUNDLER_VERSION
RUN bundle
RUN bundle update --source bundler --local

COPY . $PROJECT_DIR

EXPOSE 3000</code></pre><h2 id="ha8b2ce2ca8">compose.yamlの作成</h2><p>PostgreSQLとRailsのコンテナを構築するために<code>compose.yaml</code>を作ります。</p><pre><code class="language-yaml">services:
  web:
    image: rails:latest
    build: .
    volumes:
      - .:/root/rails-project
    command: bash -c &quot;rm -f tmp/pids/server.pid &amp;&amp; rails s -b 0.0.0.0&quot;
    environment:
      POSTGRES_DEFAULT_USER: postgres
      POSTGRES_DEFAULT_PASSWORD: password
    ports:
      - &quot;3000:3000&quot;
    depends_on:
      - db

  db:
    image: postgres:latest
    volumes:
      - postgres-volume:/var/lib/postgresql/data
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password

volumes:
  postgres-volume:</code></pre><h2 id="h818d04f982">GemfileとGemfile.lockの作成</h2><p><code>Gemfile</code>を作り、以下の内容を記述します。</p><pre><code>source &quot;https://rubygems.org&quot;

gem &quot;rails&quot;</code></pre><p>空の<code>Gemfile.lock</code>も作っておきます。</p><h2 id="h1d589f1d35">Railsのイメージの作成</h2><pre><code class="language-bash">docker compose build</code></pre><h2 id="hbe10bea12e">プロジェクトの作成</h2><pre><code class="language-bash">docker compose run web rails new . --minimal --skip-asset-pipeline -d postgresql --force</code></pre><ul><li>--minimal: Railsの追加機能であるAction Cable, Action Mailbox, Action Textをセットアップしない</li><li>--skip-asset-pipeline: アセットパイプラインをセットアップしない</li><li>-d: データベースを指定する</li><li>--force: ファイルを上書きする</li></ul><p><code>Gemfile.lock</code>が更新されたので、改めてRailsのイメージをビルドします。</p><pre><code class="language-bash">docker compose build</code></pre><h2 id="h0a5645141d">データベースと通信できるようにする</h2><p><code>config/database.yml</code> のdefaultセクションにhost, username, passwordを追加します。</p><pre><code class="language-yaml">default: &amp;default
  adapter: postgresql
  host: db
  username: &lt;%= ENV[&quot;POSTGRES_DEFAULT_USER&quot;] %&gt;
  password: &lt;%= ENV[&quot;POSTGRES_DEFAULT_PASSWORD&quot;] %&gt;
  encoding: unicode
  # For details on connection pooling, see Rails configuration guide
  # https://guides.rubyonrails.org/configuring.html#database-pooling
  pool: &lt;%= ENV.fetch(&quot;RAILS_MAX_THREADS&quot;) { 5 } %&gt;

...</code></pre><p>データベースを作成してマイグレーションをします。</p><pre><code class="language-bash">docker compose run web rails db:create db:migrate</code></pre><p>デフォルトのページを表示してみましょう。以下のコマンドを実行して、<a href="http://localhost:3000" target="_blank" rel="noopener noreferrer nofollow">http://localhost:3000</a>にアクセスします。</p><pre><code class="language-bash">docker compose up</code></pre><figure><img src="https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/62826e0184c84bfca34393223a442a9a/%E3%82%B9%E3%82%AF%E3%83%AA%E3%83%BC%E3%83%B3%E3%82%B7%E3%83%A7%E3%83%83%E3%83%88%202023-08-28%2011.54.02.png" alt="" width="3024" height="1890"></figure><h2 id="h446a12fbcc">Makefileの作成</h2><p>Docker Composeのコマンドは長いので、エイリアスを設定するために<code>Makefile</code>を作ります。</p><pre><code class="language-makefile"># イメージをビルドする
.PHONY: build
build:
	docker compose build

# イメージをビルドしてコンテナをフォアグラウンドで実行する
.PHONY: up
up:
	docker compose up --build --remove-orphans || true

# コンテナを停止する
.PHONY: stop
stop:
	docker compose stop

# コンテナを停止して削除する
.PHONY: down
down:
	docker compose down --remove-orphans

# バックグラウンドでコンテナを実行し、webコンテナに入る
.PHONY: exec-web
exec-web:
	docker compose up --build --remove-orphans -d
	docker compose exec web bash

# バックグラウンドでコンテナを実行し、dbコンテナに入る
.PHONY: exec-db
exec-db:
	docker compose up --build --remove-orphans -d
	docker compose exec db bash

# rails db:createを実行する
.PHONY: create
create:
	docker compose run web rails db:create

# rails db:migrateを実行する
.PHONY: migrate
migrate:
	docker compose run web rails db:migrate

# rails db:seedを実行する
.PHONY: seed
seed:
	docker compose run web rails db:seed

# rails db:purgeを実行する
.PHONY: purge
purge:
	docker compose run web rails db:purge

# rails db:dropを実行する
.PHONY: drop
drop:
	docker compose run web rails db:drop</code></pre><h2 id="hbadc803201">GNU Makeのインストール</h2><p>Makefileを扱うにはGNU Makeが必要です。UbuntuとMacOSにおけるインストールコマンドを示します。</p><h3 id="hc578189bb9">Ubuntu</h3><pre><code class="language-bash">sudo apt install -y build-essential</code></pre><h3 id="h5b4fb634f9">MacOS</h3><pre><code class="language-bash">brew install make
echo &apos;export PATH=&quot;$(brew --prefix)/opt/make/libexec/gnubin:$PATH&quot;&apos; &gt;&gt; ~/.zshrc
source ~/.zshrc</code></pre><h1 id="h5c0da2f09e">Viteのインストール</h1><p>package.jsonを作ります。</p><pre><code class="language-json">{
  &quot;devDependencies&quot;: {
  }
}</code></pre><p>空の<code>pnpm-lock.yaml</code>も作っておきます。</p><p><code>.pnpm-store</code>を<code>.gitignore</code>に追加します。</p><pre><code># See https://help.github.com/articles/ignoring-files for more about ignoring files.
#
# If you find yourself ignoring temporary files generated by your text editor
# or operating system, you probably want to add a global ignore instead:
#   git config --global core.excludesfile &apos;~/.gitignore_global&apos;

# Ignore bundler config.
/.bundle

# Ignore all logfiles and tempfiles.
/log/*
/tmp/*
!/log/.keep
!/tmp/.keep

# Ignore pidfiles, but keep the directory.
/tmp/pids/*
!/tmp/pids/
!/tmp/pids/.keep


/public/assets

# Ignore master key for decrypting credentials and more.
/config/master.key

# Vite Ruby
/public/vite*
node_modules
# Vite uses dotenv and suggests to ignore local-only env files. See
# https://vitejs.dev/guide/env-and-mode.html#env-files
*.local

.pnpm-store</code></pre><p>以下のコマンドでViteをインストールします。</p><pre><code class="language-bash">make exec-web

bundle add vite_rails
bundle exec vite install

exit</code></pre><p><code>bin/dev</code>を作り、以下の内容を記述します。</p><pre><code>#!/usr/bin/env sh

if ! gem list foreman -i --silent; then
  echo &quot;Installing foreman...&quot;
  gem install foreman
fi

exec foreman start -f Procfile.dev &quot;$@&quot;</code></pre><p><code>bin/dev</code>に実行権限を付与します。</p><pre><code class="language-bash">sudo chmod u+x bin/dev</code></pre><p><code>Procfile.dev</code>を変更して、ホストOSからコンテナの3000番ポートにアクセスできるようにします。</p><pre><code>vite: bin/vite dev
web: bin/rails s -p 3000 -b 0.0.0.0</code></pre><p><code>config/vite.json</code>のdevelopmentセクションにhostセクションを追加して、ホストOSからコンテナの3036番ポートにアクセスできるようにします。</p><pre><code class="language-json">{
  &quot;all&quot;: {
    &quot;sourceCodeDir&quot;: &quot;app/frontend&quot;,
    &quot;watchAdditionalPaths&quot;: []
  },
  &quot;development&quot;: {
    &quot;autoBuild&quot;: true,
    &quot;publicOutputDir&quot;: &quot;vite-dev&quot;,
    &quot;host&quot;: &quot;0.0.0.0&quot;,
    &quot;port&quot;: 3036
  },
  &quot;test&quot;: {
    &quot;autoBuild&quot;: true,
    &quot;publicOutputDir&quot;: &quot;vite-test&quot;,
    &quot;port&quot;: 3037
  }
}</code></pre><p><code>Dockerfile</code>にEXPOSEセクションを追加して、ホストOSからコンテナの3036番ポートにアクセスできるようにします。</p><pre><code class="language-dockerfile">FROM ruby:3.2.2

RUN apt-get update
RUN curl -sL http s://deb.nodesource.com/setup_18.x | bash -
RUN apt-get install -y nodejs npm
RUN npm install -g pnpm

ENV PROJECT_DIR /root/rails-project

RUN mkdir $PROJECT_DIR
WORKDIR $PROJECT_DIR

COPY Gemfile $PROJECT_DIR/Gemfile
COPY Gemfile.lock $PROJECT_DIR/Gemfile.lock

ENV BUNDLER_VERSION 2.4.19

RUN gem update --system
RUN gem install bundler -v $BUNDLER_VERSION
RUN bundle
RUN bundle update --source bundler --local

COPY . $PROJECT_DIR

EXPOSE 3000
EXPOSE 3036</code></pre><p><code>compose.yaml</code>を変更して、webコンテナの起動時にbin/devを実行するようにします。また、ホストOSの3036番ポートからwebコンテナの3036番ポートにアクセスするようにします。</p><pre><code class="language-yaml">services:
  web:
    image: rails:latest
    build: .
    volumes:
      - .:/root/rails-project
    command: bash -c &quot;rm -f tmp/pids/server.pid &amp;&amp; bin/dev&quot;
    environment:
      POSTGRES_DEFAULT_USER: postgres
      POSTGRES_DEFAULT_PASSWORD: password
    ports:
      - &quot;3000:3000&quot;
      - &quot;3036:3036&quot;
    depends_on:
      - db

  db:
    image: postgres:latest
    volumes:
      - postgres-volume:/var/lib/postgresql/data
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password

volumes:
  postgres-volume:</code></pre><h1 id="h4d2e82f7e0">homeコントローラーの作成</h1><p>indexアクションを持つhomeコントローラーを作ります。</p><pre><code class="language-bash">make exec-web

rails g controller home index

exit</code></pre><p><code>config/routes.rb</code>を変更して、rootにアクセスした際にhome#indexのビューが表示されるようにします。</p><pre><code class="language-ruby">Rails.application.routes.draw do
  root &quot;home#index&quot;
  # Define your application routes per the DSL in https://guides.rubyonrails.org/routing.html

  # Defines the root path route (&quot;/&quot;)
  # root &quot;articles#index&quot;
end</code></pre><p><code>app/assets/stylesheets</code>を削除し、<code>app/frontend/entrypoints/application.css</code>を作ります。</p><p><code>app/views/layouts/application.html.erb</code>の<code>stylesheet_link_tag</code>を<code>vite_stylesheet_tag</code>に置き換えます。</p><pre><code class="language-erb">&lt;!DOCTYPE html&gt;
&lt;html&gt;
  &lt;head&gt;
    &lt;title&gt;RailsProject&lt;/title&gt;
    &lt;meta name=&quot;viewport&quot; content=&quot;width=device-width,initial-scale=1&quot;&gt;
    &lt;%= csrf_meta_tags %&gt;
    &lt;%= csp_meta_tag %&gt;

    &lt;%= vite_stylesheet_tag &apos;application&apos; %&gt;
    &lt;%= vite_client_tag %&gt;
    &lt;%= vite_javascript_tag &apos;application&apos; %&gt;
    &lt;!--
      If using a TypeScript entrypoint file:
        vite_typescript_tag &apos;application&apos;

      If using a .jsx or .tsx entrypoint, add the extension:
        vite_javascript_tag &apos;application.jsx&apos;

      Visit the guide for more information: https://vite-ruby.netlify.app/guide/rails
    --&gt;

  &lt;/head&gt;

  &lt;body&gt;
    &lt;%= yield %&gt;
  &lt;/body&gt;
&lt;/html&gt;</code></pre><p><code>make up</code>を実行してから、<a href="http://localhost:3000" target="_blank" rel="noopener noreferrer nofollow">http://localhost:3000</a>にアクセスすると、コンソールに<code>Vite ⚡️ Rails</code>と表示されます。</p><figure><img src="https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/d01e110a1f35428b97733c6dc681436e/%E3%82%B9%E3%82%AF%E3%83%AA%E3%83%BC%E3%83%B3%E3%82%B7%E3%83%A7%E3%83%83%E3%83%88%202023-08-28%2012.40.24.png" alt="" width="3024" height="1890"></figure><h1 id="h24f46f17b4">config/routes.rbとerbファイルに対してHMRを有効にする</h1><p>HMR用のプラグインをインストールします。</p><pre><code class="language-bash">make exec-web

pnpm add -D vite-plugin-full-reload

exit</code></pre><p><code>vite.config.ts</code>に<code>vite-plugin-full-reload</code>を追加します。</p><pre><code class="language-typescript">import { defineConfig } from &apos;vite&apos;
import RubyPlugin from &apos;vite-plugin-ruby&apos;
import FullReload from &apos;vite-plugin-full-reload&apos;

export default defineConfig({
  plugins: [
    RubyPlugin(),
    FullReload([&apos;config/routes.rb&apos;, &apos;app/views/**/*&apos;], { delay: 100 }),
  ],
})</code></pre><p><code>make up</code>を実行してから<code>app/views/home/index.html.erb</code>を変更してみましょう。HMRが動作していることを確認できます。</p><figure><img src="https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/f51d0d1af58e425e8916df40ab8d1ad8/%E7%94%BB%E9%9D%A2%E5%8F%8E%E9%8C%B2-2023-08-28-12.49.49.gif" alt="" width="3024" height="1964"></figure><h1 id="h6a003ead1f">Tailwindのインストール</h1><pre><code class="language-bash">make exec-web

pnpm add -D tailwindcss postcss autoprefixer
pnpm tailwindcss init -p

exit</code></pre><p><code>tailwind.config.js</code>を変更して、Tailwindがコンパイルされるようにします。</p><pre><code class="language-javascript">/** @type {import(&apos;tailwindcss&apos;).Config} */
module.exports = {
  content: [
    &apos;./app/views/**/*.html.erb&apos;,
    &apos;./app/frontend/**/*.{js,ts,jsx,tsx,vue,svelte}&apos;,
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}</code></pre><p><code>app/frontend/entrypoints/application.css</code>にTailwindを追加します。</p><pre><code class="language-css">@import &quot;tailwindcss/base&quot;;
@import &quot;tailwindcss/components&quot;;
@import &quot;tailwindcss/utilities&quot;;</code></pre><p><code>app/views/home/index.html.erb</code>にTailwindのクラスを追加して動作を確認します。</p><pre><code class="language-erb">&lt;h1 class=&quot;text-3xl font-bold underline&quot;&gt;Home#index&lt;/h1&gt;
&lt;p&gt;Find me in app/views/home/index.html.erb&lt;/p&gt;</code></pre><figure><img src="https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/141834ec7c9643db9a2e0b84f14a8825/%E3%82%B9%E3%82%AF%E3%83%AA%E3%83%BC%E3%83%B3%E3%82%B7%E3%83%A7%E3%83%83%E3%83%88%202023-08-28%2014.31.06.png" alt="" width="3024" height="1890"></figure><h1 id="h5ad05c16ef">daisyUIのインストール</h1><pre><code class="language-bash">make exec-web

pnpm add -D daisyui

exit</code></pre><p>daisyUIを<code>tailwind.config.js</code>に追加します。</p><pre><code class="language-javascript">/** @type {import(&apos;tailwindcss&apos;).Config} */
module.exports = {
  content: [
    &apos;./app/views/**/*.html.erb&apos;,
    &apos;./app/frontend/**/*.{js,ts,jsx,tsx,vue,svelte}&apos;,
  ],
  theme: {
    extend: {},
  },
  plugins: [require(&quot;daisyui&quot;)],
}</code></pre><p><code>app/views/home/index.html.erb</code>にdaisyUIのクラスを追加して動作を確認します。</p><pre><code class="language-erb">&lt;h1 class=&quot;text-3xl font-bold underline text-accent&quot;&gt;Home#index&lt;/h1&gt;
&lt;p&gt;Find me in app/views/home/index.html.erb&lt;/p&gt;</code></pre><figure><img src="https://images.microcms-assets.io/assets/99c53a99ae2b4682938f6c435d83e3d9/1b015b5c9b4d47e7b59bc33936d5310e/%E3%82%B9%E3%82%AF%E3%83%AA%E3%83%BC%E3%83%B3%E3%82%B7%E3%83%A7%E3%83%83%E3%83%88%202023-08-28%2015.13.01.png" alt="" width="3024" height="1890"></figure>
