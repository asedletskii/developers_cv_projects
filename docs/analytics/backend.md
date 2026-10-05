# Аналитика ссылок: бэкенд

Общая картина, термины и этапы — в [README](README.md). Обязательные межслойные инварианты и правила реализации собраны в [implementation contract](implementation-contract.md). Код живёт в модуле `internal/modules/resume_link`. Отдельный модуль не нужен: аналитика не существует без ссылки и удаляется вместе с ней (`ON DELETE CASCADE`).

## Схема БД

### Этап 1: визиты

Новая таблица `resume_link_views`, по строке на каждую загрузку страницы. Поля, которые обновляет клиентский батч (этап 2), лежат здесь же: это состояние визита, а не агрегат, поэтому сырой принцип не нарушается.

```sql
CREATE TABLE resume_link_views (
    id              uuid PRIMARY KEY,                 -- viewId, генерирует сервер
    link_id         uuid NOT NULL REFERENCES resume_links(id) ON DELETE CASCADE,
    visitor_id      uuid NOT NULL,                    -- из cookie cv_vid или суточный fallback
    ip_hash         bytea NOT NULL,                   -- HMAC-SHA256(ip, суточная соль), 32 байта
    device          text NOT NULL,                    -- desktop | mobile | tablet | unknown
    browser         text NOT NULL,                    -- chrome | safari | firefox | edge | yandex | other
    os              text NOT NULL,                    -- windows | macos | linux | ios | android | other
    referrer_host   text,                             -- только хост: t.me, linkedin.com, hh.ru
    traffic         text NOT NULL DEFAULT 'human',    -- human | bot | owner | suspect
    traffic_reason  text,                             -- ua:googlebot, owner, burst, ...
    layout          jsonb NOT NULL,                   -- снимок страницы на момент визита, см. «Снимок раскладки»
    started_at      timestamptz NOT NULL,
    -- этап 2, заполняет батч
    confirmed_at    timestamptz,                      -- первый батч с active_ms >= 5000 → «прочтение»
    last_seen_at    timestamptz,
    active_ms       integer NOT NULL DEFAULT 0,       -- время видимой вкладки
    max_section     text,                             -- SectionId самой дальней секции; «дальше» считается по layout.order этого визита
    exit_section    text,                             -- секция в viewport в момент ухода
    ended_at        timestamptz,
    CONSTRAINT resume_link_views_device_check  CHECK (device  IN ('desktop','mobile','tablet','unknown')),
    CONSTRAINT resume_link_views_traffic_check CHECK (traffic IN ('human','bot','owner','suspect')),
    CONSTRAINT resume_link_views_active_check  CHECK (active_ms >= 0)
);

CREATE INDEX idx_resume_link_views_link_started ON resume_link_views (link_id, started_at DESC);
CREATE INDEX idx_resume_link_views_link_visitor ON resume_link_views (link_id, visitor_id);
```

#### Снимок раскладки

Владелец меняет страницу уже после того, как по ссылке ходили: переставляет секции, добавляет проекты. Поэтому всё, что зависит от устройства страницы, считается по **снимку на момент визита**, а не по текущему состоянию резюме. Снимок пишется в `layout` при `GET /cv/:linkId` (резюме в этот момент уже загружено, дополнительных запросов нет):

```json
{
  "order": ["about", "experience", "skills", "projects", "contacts"],
  "hashes": { "about": "9f2c41aa7d31b204", "experience": "03be77d1c2a54e19", "projects": "c4410be2a8f5d63c" }
}
```

- `order` — видимые секции в том порядке, в котором их увидел посетитель (`enabled` и не пустые, как в `getVisibleSections` на фронте).
- `hashes` — `layout_hash` каждой секции: первые 8 байт SHA-256 (16 hex-символов) от канонического JSON содержимого секции (поля, влияющие на вёрстку, без служебных `id`/`updatedAt`) плюс язык резюме. Перестановка секций, смена темы и правка *других* секций хеш не меняют; добавленный проект, правка текста или язык — меняют.

Зачем это нужно:

1. **Порядок.** «Дошёл до секции N» нельзя сравнивать по текущему порядку: после перестановки старые визиты окажутся «дальше» или «ближе», чем были. Поэтому `max_section` хранится как id секции и сравнивается по индексу в `layout.order` *этого же* визита. Отчёт «до какой секции дочитывают» ([ниже](#get-apiv1lklinklinkidstatisticssections-этап-2)) вообще не зависит от порядка: считает по множеству событий `section_view`.
2. **Карта кликов.** Клик хранится в процентах внутри секции. Если секция после этого изменилась (добавили проект, высота выросла), те же проценты указывают уже на другое место. Клик получает `layout_hash` секции из снимка своего визита, а карта по умолчанию строится только по кликам с **текущим** хешем секции: то, что на экране, соответствует тому, где кликали.


### Этап 1–2: события

Расширяем существующую `resume_link_events`, а не заводим новую:

```sql
ALTER TABLE resume_link_events
    ADD COLUMN view_id uuid REFERENCES resume_link_views(id) ON DELETE CASCADE,  -- NULL у старых записей
    ADD COLUMN client_event_id uuid,                                            -- идемпотентность батча
    ADD COLUMN payload jsonb NOT NULL DEFAULT '{}';

ALTER TABLE resume_link_events DROP CONSTRAINT resume_link_events_type_check;
ALTER TABLE resume_link_events ADD CONSTRAINT resume_link_events_type_check CHECK (event_type IN (
    'view',            -- сервер, этап 1 (дублирует строку в views; оставляем для обратной совместимости отчётов)
    'contact_view',    -- сервер, payload: {"contact":"telegram"}
    'section_view',    -- клиент, payload: {"section":"experience"}
    'download',        -- клиент
    'tech_filter',     -- клиент, payload: {"technology":"Go","action":"select"|"clear"}
    'external_click',  -- клиент, payload: {"kind":"project"|"git"|"contact"|"other","target":"github.com"}
    'print',           -- клиент: печать или «сохранить как PDF» (beforeprint)
    'contact_copy',    -- клиент, payload: {"contact":"telegram"} — скопировали раскрытый контакт
    'contact_click'    -- клиент, payload: {"contact":"telegram"} — перешли по контакту (t.me, mailto, tel)
));

CREATE UNIQUE INDEX ux_resume_link_events_client_id
    ON resume_link_events (view_id, client_event_id) WHERE client_event_id IS NOT NULL;
CREATE INDEX idx_resume_link_events_view_id ON resume_link_events (view_id);
```

`payload` — `jsonb`, но с жёсткой схемой на каждый тип в домене (`domain/event.go`): неизвестные ключи отбрасываются, строки обрезаются до 64 символов, `section` и `contact` проверяются по существующим enum (`SectionId`, `ContactKey`).

`print` — отдельный тип без payload. Клиент не шлёт его чаще раза в 2 секунды (браузеры иногда стреляют `beforeprint` дважды), но сервер всё равно считает одну печать за визит с интервалом меньше 2 секунд дублем. `contact_copy` и `contact_click` несут `contact` из `ContactKey` и отличаются от серверного `contact_view` (раскрытие) тем, что описывают поведение *после* раскрытия: скопировал ли человек значение и перешёл ли по нему.

Время активности, скролл и точка ухода **не** пишутся отдельными событиями: они обновляют поля визита (`active_ms`, `max_section`, `exit_section`). Отдельное событие на каждую секцию (`section_view`) пишется один раз за визит — для отчёта «до какой секции дочитывают».

### Источник истины для просмотров

`resume_link_views` — единственный источник истины для метрики открытия/визита. Строка `resume_link_events(event_type = 'view')` является только исторической совместимостью со старой схемой и не участвует в новых расчётах `views`, `reads`, `uniqueVisitors`, `returningVisitors` и `daily`. Новый `GET /cv/:linkId` создаёт одну строку в `resume_link_views` и, если для совместимости это требуется текущему коду, может дополнительно записать событие `view` в `resume_link_events` в той же транзакции. В дальнейшем отчёты не должны считать просмотр одновременно из обеих таблиц.

### Этап 4: клики

Отдельная таблица: клики — самый объёмный поток, им не место в общей таблице событий.

```sql
CREATE TABLE resume_link_clicks (
    id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    view_id      uuid NOT NULL REFERENCES resume_link_views(id) ON DELETE CASCADE,
    link_id      uuid NOT NULL,                  -- денормализация для отчёта без JOIN
    section      text NOT NULL,                  -- SectionId
    x_pct        smallint NOT NULL,              -- 0..10000, сотые доли процента ширины секции
    y_pct        smallint NOT NULL,              -- 0..10000, от верха секции
    viewport_w   smallint NOT NULL,              -- для разбивки heatmap по брейкпоинтам
    layout_hash  text NOT NULL,                  -- хеш секции из views.layout на момент визита; сервер берёт сам, клиент не шлёт
    target       text,                           -- data-analytics-id элемента, если есть
    occurred_at  timestamptz NOT NULL
);
CREATE INDEX idx_resume_link_clicks_link_section ON resume_link_clicks (link_id, section, layout_hash, occurred_at);
```

Замечание к TODO: координаты нормализуем **относительно секции**, а не viewport. Нормализация по viewport не даёт свести клики между визитами — вертикаль зависит от прокрутки, а высота страницы разная на разных экранах. Точка «секция + процент внутри неё» сравнима между устройствами одного брейкпоинта, **пока секция не менялась**: за это отвечает `layout_hash` (см. «Снимок раскладки»).

### Уведомления (этап 3)

```sql
CREATE TABLE analytics_outbox (
    id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id      uuid NOT NULL,
    kind         text NOT NULL,          -- link_opened | link_read | contact_revealed
    payload      jsonb NOT NULL,
    dedup_key    text NOT NULL UNIQUE,   -- см. «Уведомления»
    created_at   timestamptz NOT NULL DEFAULT now(),
    sent_at      timestamptz,
    attempts     smallint NOT NULL DEFAULT 0
);
CREATE INDEX idx_analytics_outbox_pending ON analytics_outbox (created_at) WHERE sent_at IS NULL;

CREATE TABLE user_notification_settings (
    user_id            uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    telegram_chat_id   bigint,
    link_token         text UNIQUE,      -- одноразовый токен для /start в боте
    on_open            boolean NOT NULL DEFAULT false,
    on_read            boolean NOT NULL DEFAULT true,
    on_contact         boolean NOT NULL DEFAULT true,
    quiet_hours        int4range         -- часы по tz пользователя, опционально
);
```

## Идентификация посетителя

- **Cookie посетителя**: случайный UUIDv4. Ставит сервер в ответе `GET /api/v1/cv/:linkId`, если её нет. Параметры: `Path=/api/v1/cv`, `HttpOnly`, `Secure`, `SameSite=Lax`, `Max-Age=31536000`. В cookie нет ничего, кроме случайного ID, и она first-party.
  - Имя — по образцу cookie авторизации: в проде `__Secure-cv_vid`, локально без `Secure` — `cv_vid` (как `__Host-auth-flow` / `auth-flow` в `auth/transport/http/cookie.go`). Префикс `__Host-` не подходит: он требует `Path=/`, а cookie посетителя нужна только публичным ручкам.
  - С сессией (`__Host-session`) cookie посетителя не связана и не заменяет её: одна идентифицирует владельца кабинета, другая — анонимный браузер посетителя.
- **Fallback без cookie** (например, cookie запрещены браузером или это первый запрос в инкогнито): `visitor_id = UUIDv5(namespace, ip_hash || ua_hash || date)`. Такой посетитель считается тем же анонимным посетителем только в пределах одних суток; на следующий день при отсутствии cookie получается новый `visitor_id`. Это сознательная граница приватности, а не попытка постоянной идентификации.
- **`viewId`**: сервер генерирует его при каждом успешном `GET` опубликованной активной ссылки и отдаёт в заголовке ответа `X-View-Id` (плюс `Access-Control-Expose-Headers`, если фронт и API на разных origin). Клиент не может создать визит сам, он может только дописать события к выданному `viewId`. Если запись визита не удалась, `X-View-Id` не отдаётся и клиентский трекер работает как no-op для этого открытия.

Уникальные посетители = `COUNT(DISTINCT visitor_id)`. Пара «один рекрутер перечитал пять раз» и «пять разных людей» различается ровно этим числом рядом с числом открытий.

## Классификация трафика

Выполняется синхронно при `GET` и записывается в `traffic` / `traffic_reason`. Порядок проверок:

1. **Владелец** → `owner`. Сессия — cookie `__Host-session` с `Path=/` и `SameSite=Lax`, поэтому браузер сам отправляет её и на публичный `GET /api/v1/cv/:linkId` (тот же origin). На фронте ничего делать не нужно. Если запрос несёт действующую сессию пользователя, которому принадлежит резюме ссылки, — `owner`.
   - **Нельзя вешать на публичный роут `middleware.Authentication` или вызывать `authsession.Authenticator` как есть.** На запрос без cookie он пишет WARN-лог `auth.session.rejected` и метрику отказа сессии, а на публичной странице без cookie приходит почти каждый посетитель. Метрики авторизации превратятся в шум, а алерты на отказы сессий станут бесполезны.
   - Нужен отдельный «необязательный» путь: `OptionalAuthenticate(ctx, request) (userID, ok, error)` в модуле `auth`. Нет cookie — сразу `ok = false` без обращения к БД и без событий наблюдаемости. Есть cookie, но сессия недействительна (истекла, отозвана) — тоже `ok = false`, без WARN и без метрики отказа. Ошибка инфраструктуры (БД недоступна) — визит пишется как обычный, ошибка уходит в лог; выдачу резюме она не ломает.
   - Владение проверяется в `resume_link`: пользователь из сессии сравнивается с владельцем резюме ссылки. Модуль `resume_link` получает только `userID`, про cookie и сессии он не знает.
   - `CurrentSessionHandler` раз в `AUTH_SESSION_TOUCH_INTERVAL` (5 мин) обновляет `last_seen_at` сессии. Просмотр своей ссылки продлит сессию владельца — это допустимо, отдельный «только чтение» режим не нужен.
   - В режиме `AUTH_MODE=mock` (только dev) владелец определяется по `X-Mock-User-ID`. Dev-сборка фронта добавляет этот заголовок ко всем запросам, если задан `VITE_MOCK_USER_ID`, поэтому пометка `owner` проверяется локально без Telegram.
   - Предпросмотр в кабинете (`/cv/lk/:cvId/preview`) грузит резюме через `GET /lk/resume/:resumeId` и визитов не создаёт вовсе.
2. **UA-бот** → `bot`, `traffic_reason = ua:<name>`. Список паттернов — в `infrastructure/analytics/useragent` (готовая библиотека вроде `github.com/x-way/crawlerdetect` или свой список из ~50 регулярок, покрытый тестами). Сюда же пустой UA и `HeadlessChrome`.
3. **Всплеск** → `suspect`, `traffic_reason = burst`. Один `ip_hash` за 60 с открыл > 5 разных ссылок (любых пользователей) — сканер. Счётчик — in-memory LRU по `ip_hash` рядом с существующим `rate_limit`, без БД.
4. Иначе `human`.

**Прочтение** — второй, поведенческий фильтр поверх `human`. Визит становится прочтением, когда приходит батч с `active_ms >= 5000` (сервер ставит `confirmed_at`). Бот без JS или headless-сканер, закрывающийся через секунду, прочтением не станет. В отчётах `views` и `reads` считаются только по `traffic = human`; `bot`, `owner` и `suspect` остаются в сырых данных и показываются отдельно только там, где это явно предусмотрено (`filtered` или `includeFiltered=true`).

Раскрытие контакта (`contact_view`) без `viewId` или с `viewId` визита типа `bot`/`suspect` сохраняется, но в отчёт не идёт.

## Ручки

### Публичные

#### `GET /api/v1/cv/:linkId` (изменяется)

Тело ответа без изменений. Добавляется:

- заголовок `X-View-Id: <uuid>`;
- `Set-Cookie: cv_vid=...`, если cookie не было;
- `Cache-Control: no-store` (иначе промежуточный кеш съест визиты).

Запись визита — синхронный `INSERT` с таймаутом, как сейчас в `recordLinkEvent`. Ошибка записи **не** ломает выдачу резюме: логируем, заголовок `X-View-Id` не отдаём, фронт тогда просто не шлёт батчи.

Выключенная ссылка отвечает `410 Gone` (`Resume link is no longer active`), неопубликованное резюме — `404`; визит в обоих случаях не создаётся. Опционально можно считать обращения к выключенной ссылке отдельным счётчиком без визита (`inactive_hits` у ссылки): «по выключенной ссылке „Яндекс“ заходили 3 раза» — повод включить её обратно или написать адресату.

#### `POST /api/v1/cv/:linkId/contact-views` (изменяется)

В тело добавляется необязательное поле `viewId`. В событие пишется `view_id` и `payload.contact`. Всё остальное без изменений.

#### `POST /api/v1/cv/:linkId/events` (новая, этап 2)

Единая ручка для всего клиентского потока. Сюда же в этапе 4 идут клики — отдельная ручка не нужна.

```json
{
  "viewId": "8c0e…",
  "seq": 3,
  "final": false,
  "state": {
    "activeMs": 41250,
    "maxSection": "experience",
    "exitSection": "projects"
  },
  "events": [
    { "id": "b1f…", "type": "section_view", "at": 1727700000123, "section": "experience" },
    { "id": "c2a…", "type": "tech_filter",  "at": 1727700003456, "technology": "Go", "action": "select" },
    { "id": "d3b…", "type": "download",     "at": 1727700009000 }
  ],
  "clicks": [
    { "section": "projects", "x": 4210, "y": 1733, "vw": 1440, "target": "project-link", "at": 1727700010000 }
  ]
}
```

Правила:

- **Content-Type:** принимать `application/json` **и** `text/plain`. `navigator.sendBeacon` со строкой шлёт `text/plain`, а Blob с `application/json` на части браузеров отправляется с ошибкой. Парсим тело как JSON независимо от заголовка.
- **Ответ:** `204` без тела. `404`, если ссылки нет или резюме не опубликовано; `410`, если ссылка выключена — как у `GET /cv/:linkId`. `400` — только на невалидный JSON. Ручка не требует сессии и не читает её. Отдельные битые события молча отбрасываются: клиент не ретраит, так что отвечать ему ошибкой бесполезно.
- **Проверки:** `viewId` принадлежит `linkId`, визит начат не более 12 ч назад, не больше 100 событий и 200 кликов, тело не больше 32 КБ (`http.MaxBytesReader`). Рейт-лимит по IP через существующий `middleware/rate_limit`.
- **`state`** применяется монотонно: `active_ms = GREATEST(active_ms, $1)`, `max_section` — «дальше» по индексу в `layout.order` визита, `exit_section` и `last_seen_at` — последним значением по `seq`. `final: true` ставит `ended_at`. Поэтому повтор или перестановка батчей ничего не ломают.
- **Идемпотентность событий:** `INSERT ... ON CONFLICT (view_id, client_event_id) DO NOTHING`. Периодический флаш и `pagehide` могут прислать одно и то же — дубликатов не будет.
- **`at`** с клиента не доверяем: пишем `occurred_at = now()` сервера, а `at` используем только для порядка внутри батча.
- **`section_view`** пишется один раз за визит на секцию (частичный уникальный индекс или проверка в запросе).
- Батч от визита с `traffic = bot` принимаем (`204`), но не пишем.

Слои: `transport/http/public/events.go` → `application/commands/record_events.go` (валидация, применение `state`, фильтрация) → `ports.ResumeLinkViewRepository`. Команда, а не query: здесь запись, в отличие от текущего `recordLinkEvent`, прикрученного к query.

### Кабинет

Все ручки `lk` проверяют владение так же, как текущий `GetResumeLinkStatisticsQueryHandler`. Во всех отчётах, кроме ленты визитов с `includeFiltered=true`, учитывается только `traffic = 'human'`. Период задаётся параметрами `from`/`to` (ISO date, по умолчанию последние 30 дней) и `tz` (IANA, по умолчанию `UTC`).

#### `GET /api/v1/lk/link/:linkId/statistics` (расширяется)

Старые поля `views` и `contactViews` сохраняются, чтобы не ломать фронт; `views` начинает означать открытия-людей.

```json
{
  "views": 12,
  "reads": 7,
  "uniqueVisitors": 4,
  "returningVisitors": 1,
  "contactViews": 2,
  "downloads": 1,
  "firstViewAt": "2026-09-12T08:14:03Z",
  "lastViewAt": "2026-09-29T16:40:11Z",
  "medianActiveSeconds": 94,
  "devices": { "desktop": 9, "mobile": 3, "tablet": 0, "unknown": 0 },
  "referrers": [ { "host": "t.me", "views": 8 }, { "host": null, "views": 4 } ],
  "filtered": { "bot": 5, "owner": 3, "suspect": 0 }
}
```

`firstViewAt`/`lastViewAt` считаются за всё время, не за период. `filtered` показываем владельцу, чтобы было видно, что фильтр работает.

#### `GET /api/v1/lk/link/:linkId/statistics/daily?from&to&tz`

```json
{ "days": [ { "date": "2026-09-28", "views": 3, "reads": 2, "uniqueVisitors": 2 } ] }
```

Дни без визитов заполняются нулями на сервере (`generate_series`). Группировка — `date_trunc('day', started_at AT TIME ZONE $tz)`.

#### `GET /api/v1/lk/link/:linkId/views?cursor&limit=20` — лента визитов

Главный экран аналитики. Пагинация курсором по `(started_at, id)`.

```json
{
  "items": [
    {
      "id": "8c0e…",
      "startedAt": "2026-09-29T16:40:11Z",
      "visitor": { "ordinal": 2, "isReturning": true },
      "device": "mobile", "browser": "safari", "os": "ios",
      "referrerHost": "t.me",
      "isRead": true,
      "activeSeconds": 131,
      "maxSection": "experience",
      "exitSection": "projects",
      "events": [
        { "type": "contact_view", "at": "2026-09-29T16:41:02Z", "contact": "telegram" },
        { "type": "tech_filter",  "at": "2026-09-29T16:41:30Z", "technology": "Go" }
      ],
      "traffic": "human"
    }
  ],
  "nextCursor": "…"
}
```

`visitor.ordinal` — порядковый номер посетителя в рамках ссылки («Посетитель №2»). Сам `visitor_id` наружу не отдаём. Параметр `includeFiltered=true` добавляет ботов и владельца (с полем `traffic`) — для отладки фильтра.

#### `GET /api/v1/lk/link/:linkId/statistics/sections` (этап 2)

```json
{
  "reads": 7,
  "sections": [
    { "section": "about",      "reached": 7, "exits": 0 },
    { "section": "experience", "reached": 6, "exits": 2 },
    { "section": "projects",   "reached": 3, "exits": 3 }
  ]
}
```

`reached` — прочтения, в которых был `section_view`; `exits` — прочтения с таким `exit_section`. Считается только по прочтениям: у открытий без JS нет данных о секциях. Значения `reached`/`exits` считаются по идентификаторам секций из событий и не зависят от текущего порядка страницы. Порядок выдачи для UI — текущий порядок секций; если секция уже удалена, она может быть возвращена последней или скрыта клиентом. Историческая глубина визита всегда интерпретируется по `layout.order` конкретного визита.

#### `GET /api/v1/lk/link/:linkId/statistics/technologies` (этап 2)

```json
{ "technologies": [ { "name": "Go", "selections": 3, "visitors": 2 } ] }
```

#### `GET /api/v1/lk/resume/:resumeId/statistics`

Та же структура, что у статистики ссылки, плюс разбивка `links: [{ id, name, views, reads, lastViewAt }]`. Нужна для сводки на «Главной» резюме и в списке резюме.

#### `GET /api/v1/lk/link/:linkId/clicks?section&breakpoint&includeStale` (этап 4)

Сервер отдаёт сетку, а не сырые клики: бакеты 20×20 по секции, отфильтрованные по брейкпоинту (`xs` < 600, `md` < 1200, `lg`). Агрегация при чтении, как и задумано в TODO.

```json
{ "section": "projects", "grid": 20, "cells": [ { "x": 4, "y": 7, "count": 5 } ], "total": 31, "staleTotal": 12 }
```

По умолчанию в сетку попадают только клики с **текущим** `layout_hash` секции (сервер считает его по резюме в момент запроса). `staleTotal` — сколько кликов пришло на прошлые версии секции и в сетку не вошло. `includeStale=true` добавляет и их: проценты прошлых версий могут указывать мимо, поэтому на фронте это явный переключатель с оговоркой, а не поведение по умолчанию.

#### Уведомления (этап 3)

- `GET /api/v1/lk/notifications/settings` → `{ telegramConnected, onOpen, onRead, onContact }`.
- `PUT /api/v1/lk/notifications/settings` — те же булевы поля.
- `POST /api/v1/lk/notifications/telegram/link` → `{ url: "https://t.me/<bot>?start=<token>" }`. Токен одноразовый, живёт 15 минут.
- `DELETE /api/v1/lk/notifications/telegram` — отвязать бота.

## Уведомления

Поток: команда записи визита/батча/контакта → в той же транзакции `INSERT` в `analytics_outbox` → воркер-горутина в процессе бэкенда раз в 5 с забирает строки `FOR UPDATE SKIP LOCKED` → адаптер Telegram → `sent_at = now()`.

- **Когда:** `link_read` (визит стал прочтением — главное уведомление), `contact_revealed`, опционально `link_opened`. Только `traffic = 'human'`.
- **Анти-спам через `dedup_key`:** `link_read:<link_id>:<visitor_id>:<date>` — одно уведомление на посетителя в сутки; `contact_revealed:<view_id>:<contact>`. `ON CONFLICT DO NOTHING`.
- **Текст:** «Резюме по ссылке „Яндекс, Go-разработчик“ читают: мобильный, из Telegram, 2 мин, дошёл до «Опыта». Посетитель уже заходил 1 раз». Со ссылкой на ленту визитов в кабинете.
- **Ретраи:** до 5 попыток с экспоненциальной паузой, потом строка остаётся с `attempts = 5` и попадает в лог/метрику.
- **Привязка бота:** пользователь нажимает «Подключить Telegram» → открывает `t.me/<bot>?start=<token>` → вебхук бота (`POST /api/v1/integrations/telegram/webhook` с секретом в заголовке `X-Telegram-Bot-Api-Secret-Token`) находит токен и сохраняет `chat_id`.

Брокера сообщений в стеке нет (NATS и MinIO убраны из `docker-compose.prod.yml`), и для уведомлений он не нужен: outbox в Postgres проще и транзакционно надёжнее.

**Telegram-вход и бот — разные вещи.** Пользователь входит через Telegram OIDC (`oauth.telegram.org`, scope `openid profile`). Это даёт `sub` в `user_identities`, но не даёт боту права писать пользователю: Telegram разрешает боту первое сообщение только после того, как пользователь сам начал диалог. Поэтому привязка через `/start <token>` остаётся обязательной. Упростить её можно так: кнопка «Подключить Telegram» видна только вошедшим через Telegram, а после `/start` сервер может сверить Telegram ID отправителя с `provider_subject` — если OIDC отдаёт в `sub` именно Telegram ID (проверить на реальном ответе провайдера).

## Обязательные межслойные инварианты

Следующие правила являются частью контракта аналитики и должны соблюдаться backend и frontend одновременно:

1. **Один `GET` → один `viewId` → один `resume_link_views`.** Повторный `GET`, reload и новый вход после `pagehide` создают новый визит; возврат из `bfcache` продолжает существующий `viewId`.
2. **`viewId` выдаёт только backend.** Frontend не генерирует и не принимает внешний `viewId` как основание для создания визита. Все события и клики принимаются только для существующего `viewId`, принадлежащего той же ссылке.
3. **`visitor_id` не приходит из JSON событий.** Его определяет backend по cookie или описанному суточному fallback. Frontend не читает и не отправляет его.
4. **Исторический layout принадлежит визиту.** `layout.order` и `layout_hash` фиксируются при `GET`; изменение резюме после визита не переписывает историю.
5. **Координаты click-map относительны секции.** `viewport_w` используется только для выбора breakpoint; он не меняет систему координат `x_pct/y_pct`.
6. **Сырые временные метки клиента не являются источником истины.** Поле `at` используется только для порядка внутри батча; серверное `occurred_at` определяет время события в аналитике.
7. **Клиентские события вторичны по отношению к визиту.** Потерянный/невыданный `viewId` означает отсутствие клиентской аналитики, но не должен ломать публичную выдачу резюме.
8. **Все пользовательские отчёты по умолчанию работают только с `traffic = human`.** Фильтрованные классы сохраняются для диагностики и прозрачности, но не смешиваются с основными метриками.

## Приватность

- **IP в базе:** не хранится. `ip_hash = HMAC-SHA256(ip, salt_of_day)`. Соль генерируется случайно раз в сутки и держится в памяти (при нескольких инстансах — в таблице с TTL). Через сутки хеш необратим даже для нас. IP берётся из `c.ClientIP()` — нужно, чтобы `HTTP_TRUSTED_PROXIES` на проде указывал на Caddy, иначе все визиты будут с IP контейнера.
- **IP в логах:** Caddy пишет JSON-журнал доступа в stdout (глобальный `log` в `Caddyfile`), а в нём — IP клиента и полный путь `/cv/<linkId>`. По этому журналу можно восстановить, кто открывал какую ссылку, в обход обезличенной аналитики. Нужно либо ограничить ретенцию логов (30 дней, см. NFR-PRV-07 в `docs_new`), либо маскировать IP в журнале Caddy (фильтр `ip_mask`).
- **Сессии:** `auth_sessions` хранит `ip_address` и `user_agent` владельца. В аналитику они не попадают и с визитами не связываются: по сессии ставится только пометка `owner`.
- **User-Agent:** разбирается на лету в `device/browser/os`, сырая строка не сохраняется. Для отладки классификатора можно писать сырой UA только у `traffic = bot` с отдельной ретенцией 30 дней.
- **Referrer:** только хост, без пути и query.
- **Ретенция:** ночная задача удаляет `resume_link_views` старше 12 месяцев (события уйдут каскадом) и `resume_link_clicks` старше 90 дней.
- **Раскрытие:** в футер публичной страницы — строка «Владелец резюме видит статистику просмотров без персональных данных посетителя». Это техническая формулировка; юридическую формулировку для 152-ФЗ согласовать отдельно.

## Производительность

При текущих объёмах всё считается запросами по индексу `(link_id, started_at)`. Рубеж для пересмотра — сотни тысяч визитов на одну ссылку, что для резюме нереально. Если понадобится, первым шагом будет дневная материализация `resume_link_daily_stats`, а не смена хранилища.

Запись визита стоит один `INSERT` на `GET`; батч — один `UPDATE` визита и один многострочный `INSERT` событий в транзакции.

## Тестирование

- Домен: валидация payload по типам, монотонное применение `state`.
- UA-классификатор: табличный тест на реальных UA (Chrome/Safari/Firefox desktop и mobile, iPad, Googlebot, TelegramBot, HeadlessChrome, пустой).
- Интеграционные (по образцу `resume_link/.../integration_test.go`): дубликат батча не дублирует события; батч с чужим `viewId` → `404`, с выключенной ссылкой → `410`; визит с действующей сессией владельца помечен `owner` и не попадает в `statistics`; визит с сессией другого пользователя — `human`; визит без cookie и с протухшей сессией не увеличивает метрику отказов сессий и не пишет WARN; дни без визитов в `daily` — нули; `tz` сдвигает границу дня.
- Роутер: новые пути — в `router_test.go`.