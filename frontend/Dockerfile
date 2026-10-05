FROM node:22-bookworm-slim AS builder

WORKDIR /src

COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile

COPY . .
RUN yarn build

FROM caddy:2-alpine

COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=builder /src/build /srv
COPY docker-entrypoint.sh /docker-entrypoint.sh

EXPOSE 80

# Entrypoint пишет /srv/config.js из переменных окружения и запускает Caddy. CMD задаём заново:
# при смене ENTRYPOINT Docker сбрасывает CMD базового образа
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
