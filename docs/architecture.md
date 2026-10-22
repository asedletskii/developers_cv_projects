# Архитектура

## Общая картина

```
Браузер ──HTTPS──▶ серверный Caddy ──▶ frontend (Caddy со статикой)
                 (репозиторий прокси)      └─ /api/*  ──▶ backend :8080 ──▶ PostgreSQL
```

- Один origin: браузер ходит на `/api/v1/...` того же домена, поэтому CORS в проде не участвует.
- Публичные домены: `cv.developers-cv.ru` (прод) и `dev.cv.developers-cv.ru` (стейдж). Серверный Caddy и его доступ по IP живут в отдельном репозитории прокси.
- Бэкенд наружу порт не публикует, база сидит во внутренней сети `data`.

## Бэкенд (`backend/`)

Модульный монолит, внутри модуля слои `domain → application (commands/queries) → infrastructure → transport/http`. Сборка зависимостей через Wire (`make generate`).

| Модуль | Отвечает за |
|---|---|
| `auth` | вход через Telegram OIDC, сессии, `GET /auth/me`, выход |
| `user` | профиль владельца кабинета (`ru`/`en`, фото) |
| `resume` | резюме и все его секции, публичная проекция |
| `resume_link` | именованные ссылки, статус ссылки, раскрытие скрытых контактов |
| `analytics` | приём событий посетителя, отчёты |

Платформа (`internal/platform`): конфиг из переменных окружения, HTTP-роутер и ошибки, логгер, Postgres, телеметрия (OpenTelemetry, Prometheus), healthcheck.

Ключевые правила:
- **Оптимистическая блокировка.** Каждое изменение резюме и ссылки несёт `version`; при несовпадении ответ `PRECONDITION_FAILED`.
- **Строгий разбор JSON.** Неизвестное поле в теле — `400`.
- **Миграции** лежат в `backend/db/migrations` (`NNNNNN_name.up/down.sql`), применяются сервисом `migrations` на каждом деплое.
- **Mock-вход** собирается только с тегом `devauth` (`make run-mock-auth`); prod-образ без него.

## Фронтенд (`frontend/src`)

```
app/                 корень приложения, роутер, провайдеры, ErrorBoundary, PageLoader
features/
  auth/              экран ошибки входа
  main/              страница «О проекте»
  curriculumVitae/   всё про резюме
    api/             клиент, один класс CurriculumVitaeApi
    components/cv/   публичная страница и блоки резюме
    components/lk/   личный кабинет
    components/editor/  редактор секций
    sections/        реестр секций и их видов
    store/ hooks/ utils/ locale/ types/ constants/ screens/ styles/ animations/
shared/              api (axios), компоненты, хуки, словари, утилиты, тема
```

### Маршруты

| Путь | Экран |
|---|---|
| `/` | О проекте |
| `/cv/:linkId` | публичное резюме |
| `/cv/lk` | профиль и список резюме |
| `/cv/lk/:cvId` | редактор резюме (`?section=edit&tab=…`) |
| `/cv/lk/:cvId/preview` | предпросмотр резюме владельцем |
| `/cv/auth/error` | ошибка входа |

`/main`, `/about` и `/cv` перенаправляют на `/`; `/cv/lk/:legacyUserId/*` — старые адреса кабинета, перенаправляются на `/cv/lk/...`.
| `*` | 404 |

### Секции и виды

Реестр — `sections/sectionDefinitions.tsx`: для каждой из 9 секций (`about`, `techStack`, `experience`, `skills`, `education`, `courses`, `projects`, `activities`, `contacts`) там описано, как понять, что есть что показывать (`hasContent`), как отрисовать (`render(resume, style)`), попадает ли секция в навигацию и сворачивается ли.

**Виды блоков.** Вид хранится в `overview.sections[].style`:

| Значение | Вид |
|---|---|
| `first` | основной (по умолчанию) |
| `second` | «красивый»: карточки `richCard` (стекло, акцентная полоса, подъём при наведении), иконки, плитки |
| `third` | «компактный»: строки и текст без карточек |

- Незнакомое значение показывается как `first` (`resolveSectionStyle`).
- Вид выбирают в редакторе секции переключателем «Вид блока»; выбор сразу сохраняется PUT-ом всего `overview`. Для «Обо мне» переключатель только на вкладке «Профиль», вкладка «Языки» его не показывает.
- `first` в запросе не отправляется (`buildSectionRequest`), так что старый бэкенд без поля продолжает работать, пока никто не выбрал другой вид.
- Новый вид или секция: добавить компонент, подключить в `sectionDefinitions.tsx`, при необходимости расширить `SECTION_STYLES`.

### Соглашения по коду

- Каждый компонент в своей папке, наружу отдаётся через `index.ts`. Подкомпоненты лежат плоскими файлами рядом с родителем, без вложенных `components/`.
- Имя ребёнка не повторяет имя родителя (`ExperienceSection` → `Job`) и не совпадает с компонентами библиотек (MUI, react-router).
- Компонент, нужный нескольким родителям, поднимается на уровень выше (`TechnologyList`, `IconTile`, `SparkleField`).
- Конечные наборы значений — `enum`; значение совпадает с кодом из API. Доменные enum лежат в `features/curriculumVitae/types/enums.ts`, общие в `shared/types`.
- Чистые функции — в `utils`, каждая в своей папке с тестом рядом (`formatAge/formatAge.ts`, `formatAge.test.ts`, `index.ts`).
- Не компоненты (реестр секций, маршруты, помощники) лежат отдельными модулями.
- Остальные правила форматирования и стиля — [frontend/STYLEGUIDE.md](../frontend/STYLEGUIDE.md).

### Тема и стили

- Тема MUI: берётся из сохранённого выбора, иначе из системной (`prefers-color-scheme`); акцент жёлтый, в светлой теме текст и контурные кнопки синие.
- Общие стили кода: `styles/glassCard`, `styles/richCard`, `styles/accent` (`iconTile`, `badge`, `secondaryButton`, `dialogPaper`).
- Появление блоков — классы `.reveal`, `.reveal-group`, `.reveal-item` (GSAP).

### Данные и сохранение

- Редактор держит форму в react-hook-form; сохранение по секции: `buildSectionRequest` собирает тело, `saveResumeSection` шлёт PUT, затем резюме перечитывается.
- `normalizeResume` приводит ответ бэка к безопасному виду (пустые массивы вместо `null`).
- Публичное резюме и кабинет читают одну и ту же форму `Resume`; публичная проекция беднее (нет `role`, `gitLink`, `contributors`).

## Локальная разработка

| Что | Как |
|---|---|
| Свой бэкенд | `make run-mock-auth` в `backend/`, фронт ходит на `localhost:8080` |
| Удалённый бэкенд (стейдж) | `VITE_PROXY_TARGET=https://dev.cv.developers-cv.ru yarn start` |
| Сессия удалённого бэкенда | войти там в браузере, значение cookie `__Host-session` положить в `PROXY_SESSION` (без префикса `VITE_`, в бандл не попадает) |
| Трекинг и DNT | `RESPECT_DO_NOT_TRACK=false` отключает уважение сигнала «Не отслеживать» |
| Turnstile | ключ сайта — `VITE_TURNSTILE_SITE_KEY` (на сборке) |

Переменные Vite читаются из `frontend/.env` или окружения. Файлы `.env*` в git не попадают (кроме `.env.example`).
