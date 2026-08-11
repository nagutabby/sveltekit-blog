FROM node:24.19.0-alpine AS base
RUN corepack enable
WORKDIR /app

FROM base AS deps
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

FROM deps AS build
# SSG化によりコンテンツ取得(backendへのConnect RPC)はビルド時にしか
# 行われないため、BACKEND_URLはここでのみ必要。railway.internalの
# private networkingはビルドステージからは到達できないため、backendの
# public URL(wrangler.tomlのBACKEND_ORIGINと同じ値)を渡すこと。
ARG BACKEND_URL
ENV BACKEND_URL=$BACKEND_URL
COPY . .
RUN pnpm build

FROM base AS prod-deps
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile --prod --ignore-scripts

FROM base AS runtime
ENV NODE_ENV=production
COPY --from=prod-deps /app/node_modules ./node_modules
COPY --from=build /app/build ./build
COPY --from=build /app/package.json ./package.json

EXPOSE 3000
CMD ["node", "build/index.js"]
