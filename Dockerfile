FROM node:20-slim

RUN corepack enable
RUN corepack prepare pnpm@9.12.0 --activate

WORKDIR /

COPY package.json pnpm-lock.yaml ./

RUN pnpm install --frozen-lockfile

COPY . .

RUN pnpm build

EXPOSE 3000

CMD ["node", "build"]
