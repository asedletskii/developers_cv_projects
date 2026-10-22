# Деплой

## Окружения

| | Прод | Стейдж |
|---|---|---|
| Домен | `cv.developers-cv.ru` | `dev.cv.developers-cv.ru` |
| Ветка | `main` | `stage` |
| Compose | `docker-compose.prod.yml` | `docker-compose.stage.yml` (проект `cv-stage`) |
| Папка на сервере | `secrets.DEPLOY_PATH` | `secrets.DEPLOY_PATH_STAGE` |
| Теги образов | `latest` и `sha-<коротко>` | `stage` |
| Ключ Turnstile (сборка фронта) | `vars.TURNSTILE_SITE_KEY` | `vars.TURNSTILE_SITE_KEY_STAGE` |

Оба стека живут на одном VPS. Наружу порты не публикуются: единственная точка входа — серверный Caddy (отдельный репозиторий прокси), он ходит во фронтенд по общей сети `edge` по алиасу (`cv-prod`, `cv-stage`).

Образы собираются в CI и лежат в GHCR (`ghcr.io/<owner>/<repo>/{frontend,backend}`). На сервере ничего не собирается.

## Сети прод-стека

| Сеть | Кто в ней | Зачем |
|---|---|---|
| `edge` (external) | серверный Caddy, frontend | вход снаружи |
| `proxy` | frontend, backend | `/api` проксируется на `backend:8080`; frontend закреплён за `CADDY_PROXY_IP` (по умолчанию `172.30.240.2`), это единственный доверенный прокси бэкенда (`HTTP_TRUSTED_PROXIES`) |
| `data` (internal) | backend, migrations, postgres | база недоступна снаружи |

Подсеть `proxy` задаётся `PROXY_NETWORK_SUBNET`; динамические адреса выдаются из `ip_range` `172.30.240.128/25`, чтобы не столкнуться с закреплённым адресом frontend.

## CI/CD

`.github/workflows/ci.yml` (push в `main` и PR):

1. **backend**: gofmt, golangci-lint, `go vet`, unit-тесты (в том числе с тегом `devauth`), интеграционные тесты (Testcontainers), проверка Wire, `govulncheck`, сборка, валидация compose.
2. **frontend**: `yarn audit` (prod-зависимости), lint, typecheck, тесты, сборка; проверка, что dev-код (`X-Mock-User-ID`) не попал в бандл.
3. **backend-docker / frontend-docker** (только `main`): сборка и push образов; для backend — скан Trivy (CRITICAL, HIGH).
4. **deploy** (только `main`): `scp` файла `docker-compose.prod.yml` и `backend/db/migrations` на сервер, затем по SSH `pull` → `up -d --remove-orphans` → ожидание healthcheck бэкенда (до 60 с). Если healthcheck не прошёл, backend откатывается на предыдущий образ, в лог выводятся логи backend и migrations, джоба падает.

`.github/workflows/stage.yml` (push в `stage`): vet и unit-тесты, сборка образов с тегом `stage`, деплой тем же способом без отката.

`scp` не удаляет лишние файлы: убранная из репозитория миграция остаётся на сервере и продолжит применяться, её нужно удалить руками.

### Секреты и переменные GitHub

| Имя | Тип | Что это |
|---|---|---|
| `DEPLOY_HOST`, `DEPLOY_USER`, `DEPLOY_PORT` | secret | SSH-доступ к серверу (`DEPLOY_PORT` по умолчанию 22) |
| `DEPLOY_SSH_KEY` | secret | приватный ключ для SSH |
| `DEPLOY_PATH`, `DEPLOY_PATH_STAGE` | secret | папки прод- и стейдж-стеков на сервере |
| `TURNSTILE_SITE_KEY`, `TURNSTILE_SITE_KEY_STAGE` | variable | публичный ключ сайта Turnstile, вшивается в сборку фронтенда |

## Конфигурация на сервере

В папке стека один раз вручную создаётся `.env` (в git не попадает, CI его не трогает). За основу берётся корневой [.env.example](../.env.example): он служит и для подстановок `${VAR}` в compose, и (как `env_file`) для backend-контейнера. Описание переменных бэкенда — `backend/.env.example`.

Что обязательно заполнить:

| Группа | Переменные |
|---|---|
| Образы | `FRONTEND_IMAGE`, `BACKEND_IMAGE`, `IMAGE_TAG` |
| БД | `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD` (compose падает без них) |
| Вход | `AUTH_MODE=telegram`, `AUTH_TELEGRAM_CLIENT_ID`, `AUTH_TELEGRAM_CLIENT_SECRET`, `AUTH_TELEGRAM_REDIRECT_URL`, `AUTH_COOKIE_SECURE=true` |
| Turnstile | `TURNSTILE_ENABLED`, `TURNSTILE_SECRET_KEY`, `TURNSTILE_EXPECTED_HOSTNAME`, `TURNSTILE_EXPECTED_ACTION=contact_view` |
| Аналитика | `ANALYTICS_VISITOR_SECRET` (случайная секретная строка; заменить значение-заглушку из `.env.example`) |
| Фронт | `RESPECT_DO_NOT_TRACK` (применяется при перезапуске контейнера, пересборка не нужна) |

Сервис `migrations` подставляет пароль в URL подключения как есть, поэтому пароль БД делайте без `@ : / ? #` и пробелов. Параметры подключения backend собираются в `backend/internal/platform/config/postgres.go`.

## Миграции

SQL-файлы лежат в `backend/db/migrations`. Сервис `migrations` (`migrate/migrate`) монтирует их с диска сервера и выполняет `up` на каждом деплое до старта backend. Нумерация сквозная: не переименовывать и не переиспользовать номера уже применённых миграций (иначе `migrate` падает на несовпадении истории). Новая схема — новые файлы `NNNNNN_name.up.sql` и `.down.sql`.

Локально: `make migrate-up`, `make migrate-down`, `make migrate-version` (в `backend/`).

## Первый запуск стека

1. На сервере создать папку стека и положить `.env` по образцу.
2. Убедиться, что существует внешняя сеть `edge` и прокси подключён к ней.
3. Запушить в `main` (или `stage`): CI доставит compose-файл и миграции, поднимет стек.
4. Проверить: `docker compose -f docker-compose.prod.yml ps` и логи `backend`, `migrations`.

## Откат и диагностика

```bash
cd <папка стека>
docker compose -f docker-compose.prod.yml logs --tail=100 backend migrations
docker compose -f docker-compose.prod.yml exec -T backend /app/developers-cv-api healthcheck
```

Автоматический откат срабатывает только для backend при провале healthcheck. Откат фронтенда или конкретной версии — задать `IMAGE_TAG=sha-<коротко>` в `.env` и выполнить `pull` + `up -d`. Откат миграции — только вручную (`migrate down`) и с осторожностью: данные из удаляемых колонок теряются.

## Локальный стек

`backend/docker-compose.yml` поднимает PostgreSQL и (при необходимости) backend с миграциями; наружу публикует `:8080`. Подробности запуска — в [README](../README.md#быстрый-старт).
