FROM node:22-bookworm-slim AS builder

WORKDIR /src

COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile

COPY . .

ARG VITE_TURNSTILE_SITE_KEY
RUN yarn build

FROM caddy:2-alpine

COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=builder /src/build /srv

EXPOSE 80
