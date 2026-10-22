# developers-cv.ru

Конструктор резюме для разработчиков: человек собирает резюме в личном кабинете, публикует его по именованной ссылке (`/cv/<linkId>`) и видит, кто и как его читал.

- **Фронтенд** — React 19, TypeScript, Vite, MUI 7, Redux Toolkit, react-hook-form, i18next (ru/en), GSAP.
- **Бэкенд** — Go, Gin, GORM, PostgreSQL; модульный монолит (`auth`, `user`, `resume`, `resume_link`, `analytics`).
- **Вход** — Telegram OIDC, сессия в cookie `__Host-session`.
- **Деплой** — Docker Compose на одном VPS; HTTPS и маршрутизацию по доменам делает серверный Caddy из отдельного репозитория прокси.

## Что есть в продукте

| Часть | Что делает |
|---|---|
| Публичная страница `/cv/:linkId` | резюме без авторизации; контакты скрыты за Cloudflare Turnstile; кнопка «Скачать» собирает файл на странице |
| Личный кабинет `/cv/lk` | список резюме, редактор по секциям, ссылки с разными именами (HH, Telegram…), аналитика по ссылке |
| Секции и виды | 9 секций, порядок и видимость настраиваются; у каждой секции три вида (`first`, `second` — «красивый», `third` — «компактный») |
| Аналитика | визиты, дочитывание, клики, интерес к технологиям; без cookie посетителя и без хранения IP в открытом виде |

## Структура репозитория

```
backend/    Go API, миграции (db/migrations), docker-compose для локальной БД
frontend/   SPA, Dockerfile, Caddyfile статики
docker-compose.prod.yml   прод-стек
docker-compose.stage.yml  стейдж-стек
.github/workflows/        ci.yml (main → прод), stage.yml (ветка stage → стейдж)
docs/       документация
```

## Быстрый старт

Фронтенд (Node из `frontend/.nvmrc`, Yarn):

```bash
cd frontend
yarn install
yarn start
```

Бэкенд (Go из `backend/.go-version`, Docker для PostgreSQL):

```bash
cd backend
cp .env.example .env
make compose-up        # PostgreSQL
make migrate-up        # миграции
make run-mock-auth     # API на :8080 с mock-входом (только dev-сборка)
```

Vite проксирует `/api` на `http://localhost:8080`. Другой бэкенд (например стейдж) задаётся через `VITE_PROXY_TARGET`; подробности в [docs/architecture.md](docs/architecture.md#локальная-разработка).

## Проверки перед коммитом

```bash
# frontend/
yarn lint && yarn typecheck && yarn test && yarn build

# backend/
make check             # нужен Docker (интеграционные тесты)
```

CI гоняет то же самое на каждый push и PR.

## Документация

| Файл | О чём |
|---|---|
| [docs/architecture.md](docs/architecture.md) | устройство фронта и бэка, соглашения по коду, секции и их виды, локальная разработка |
| [docs/api.md](docs/api.md) | HTTP-контракт: правила, ручки, тела запросов, ошибки, аналитика |
| [docs/openapi.yaml](docs/openapi.yaml) | то же в виде OpenAPI 3.0 (для Swagger/Redoc и генерации клиентов) |
| [docs/deploy.md](docs/deploy.md) | окружения, CI/CD, секреты, переменные, миграции, откат |
| [frontend/STYLEGUIDE.md](frontend/STYLEGUIDE.md) | стайлгайд кода фронтенда |
