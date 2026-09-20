# HTTP API

Источник истины — код: `backend/internal/modules/*/transport/http` (маршруты и тела запросов), миграции в `backend/db/migrations`. Клиент фронтенда — `frontend/src/features/curriculumVitae/api/curriculumVitae.api.ts`, тела PUT собирает `utils/saveResumeSection/buildSectionRequest.ts`.

Машиночитаемая версия — [openapi.yaml](openapi.yaml) (OpenAPI 3.0, проверена линтером Redocly). При расхождении прав код бэкенда.

## Общие правила

**Адрес.** Всё под `/api/v1`. В проде тот же origin, что и у сайта (Caddy проксирует `/api`), локально — через прокси Vite.

**Авторизация.** `/api/v1/lk/*` требует сессии (cookie `__Host-session`); без неё `401 UNAUTHORIZED`. Публичные ручки (`/cv/*`) открыты. Чужой или несуществующий ресурс — `404`, а не `403`.

**Идентификаторы.** UUID. Невалидный UUID в пути — `400`. Идентификаторы записей создаёт бэкенд; в тело списков `id` передаётся только для записей, которые бэкенд уже знает (у новых `id` нет, иначе `400`).

**Тело запроса.** JSON UTF-8, не больше 1 МиБ (`413 REQUEST_TOO_LARGE`). Разбор строгий: неизвестное поле — `400`. Строки обрезаются по краям. **`PUT` — полная замена** секции: что не передано, сбрасывается (необязательное) или даёт `400` (обязательное).

**Версия.** Резюме и ссылка защищены оптимистической блокировкой: каждый `PUT`/`DELETE` несёт `version`, полученный из последнего `GET`. Должна быть положительной (иначе `400 Version must be positive`); устаревшая — `412 PRECONDITION_FAILED` (резюме уже изменили в другой вкладке). Версия резюме общая на все его секции.

**Ответы на запись.** `POST` → `201 {"id": "<uuid>"}`; `PUT` и `DELETE` → `204` без тела. Актуальное состояние фронтенд перечитывает `GET`-ом.

**Ошибки.**

```json
{ "error": { "code": "BAD_REQUEST", "message": "Version must be positive", "requestId": "3f6e9da8-…" } }
```

| Статус | `code` | Когда |
|---|---|---|
| 400 | `BAD_REQUEST` | неверное тело, неизвестное поле, нарушено ограничение |
| 401 | `UNAUTHORIZED` | нет или недействительна сессия |
| 403 | `FORBIDDEN` | CSRF-защита: изменяющий запрос (не GET/HEAD/OPTIONS) с неразрешённым `Origin` |
| 404 | `NOT_FOUND` | не найдено или чужое |
| 409 | `CONFLICT` | дубликат (например повтор языка в списке `languages`) |
| 410 | `GONE` | публичная ссылка отключена |
| 412 | `PRECONDITION_FAILED` | устаревшая `version` |
| 413 | `REQUEST_TOO_LARGE` | тело больше лимита |
| 429 | `TOO_MANY_REQUESTS` | лимит запросов |
| 500 / 503 | `INTERNAL_ERROR` / `SERVICE_UNAVAILABLE` | ошибка сервера; искать по `requestId` (он же заголовок `X-Request-ID`) |

**Форматы.** Дата `YYYY-MM-DD`; месяц `YYYY-MM` (`startDate`, `endDate` опыта, `endDate` курса, `date` активности); год `YYYY` (`endDate` образования); ссылки только `http`/`https`; цвет `#RGB` или `#RRGGBB`; время RFC 3339 UTC.

## Справочники

| Что | Значения |
|---|---|
| Язык резюме | `ru en fr de kz zh es` (как `RESUME_LANGUAGES` на фронте); через API после создания не меняется |
| Языки, которыми владеет автор | `ru en fr de kz zh es`; уровень `A1 A2 B1 B2 C1 C2` |
| Статус резюме | `draft` → `opened` ⇄ `closed` (в `draft` вернуть нельзя) |
| Статус ссылки | `active`, `inactive` |
| Формат работы | `fullTime partTime office remote hybrid` |
| Готовность к предложениям | `open considering busy` |
| Секции | `about techStack experience skills education courses projects activities contacts` |
| Вид секции (`style`) | `first second third`, по умолчанию `first` |
| Контакты | `telegram email phone` |
| Вид активности | `openSource hackathon speaker volunteering mentoring community` |
| Уровень образования | `0` среднее, `1` бакалавриат, `2` специалитет, `3` магистратура, `4` аспирантура, `5` кандидат наук, `6` доктор наук |

## Ручки

### Авторизация

| Метод и путь | Что делает |
|---|---|
| `GET /auth/telegram` | начало входа (редирект в Telegram) |
| `GET /auth/telegram/callback` | завершение входа, ставит cookie сессии |
| `GET /auth/me` | текущая сессия: `{ userId, authenticated }` |
| `POST /auth/logout` | выход |

### Кабинет: владелец и резюме

| Метод и путь | Тело | Ответ |
|---|---|---|
| `GET /lk/profile` | — | `{ userId, ru, en, photoPath? }` |
| `PUT /lk/profile` | `{ ru, en, photoPath? }` (`ru`/`en`: `name surname patronymic position location`) | 204 |
| `GET /lk/resumes` | — | массив `{ id, title, language, status, version }` |
| `POST /lk/resume` | `{ title, language }` | 201 `{ id }`; на один язык можно создать несколько резюме |
| `GET /lk/resume/:resumeId` | — | резюме целиком (ниже) |
| `PUT /lk/resume/:resumeId` | `{ version, title }` | 204 |
| `PUT /lk/resume/:resumeId/status` | `{ version, status }` | 204 |
| `DELETE /lk/resume/:resumeId` | `{ version }` | 204 |

`GET /lk/resume/:resumeId` возвращает `id, title, language, status, version, createdAt, updatedAt` и блоки `overview`, `profile?`, `contacts?`, `skills?`, `languages[]`, `technologies[]`, `experiences[]`, `educations[]`, `courses[]`, `activities[]`, `projects[]`. `profile`, `contacts`, `skills` отсутствуют, пока их не сохранили.

### Кабинет: секции резюме (`PUT /lk/resume/:resumeId/<секция>`)

Все — полная замена, `204`. У одиночных секций поля лежат в корне рядом с `version`, у списков — в `items`.

| Секция | Тело |
|---|---|
| `overview` | `{ version, theme?: { primary?, secondary? }, greeting: { first, second, secondHighlight, third }, availability?: { status, from?, note? }, sections: [{ id, enabled, collapsed, style? }] }` |
| `profile` | `{ version, id?, name, surname, patronymic?, position, photoPath, dateOfBirth, location, relocation, workFormat[] }` |
| `contacts` | `{ version, id?, phone?, telegram?, email?, hiddenContacts[], contactsText?: { title?, subtitle? } }` |
| `skills` | `{ version, id?, hardSkills[], softSkills[] }` (свободный текст; дубли в списке запрещены) |
| `languages` | `{ version, items: [{ id?, language, grade?, position }] }` |
| `technologies` | `{ version, items: [{ id?, name, position }] }` |
| `experiences` | `{ version, items: [{ id?, company, position, location?, startDate, endDate?, description?, achievements[], tasks[], stack: [{ name }] }] }` |
| `educations` | `{ version, items: [{ id?, location, university, grade, programName, profileName?, faculty, course?, endDate? }] }` |
| `courses` | `{ version, items: [{ id?, name, organization, endDate, hours?, certificateLink? }] }` |
| `activities` | `{ version, items: [{ id?, type, title, description?, date?, link? }] }` |
| `projects` | `{ version, items: [{ id?, name, description?, role?, screenshots[], link?, gitLink?, stack: [{ name }], contributors: [{ name, surname, role, contacts: { phone, telegram?, email? } }] }] }` |

Особенности `overview`:
- `theme`, `availability` необязательны; если `availability` есть, `status` обязателен.
- `sections` — порядок массива и есть порядок на странице; допустим любой поднабор, `id` уникальны. Пустой массив фронтенд читает как «все секции по умолчанию». «Обо мне» всегда первая и включена (это обеспечивает фронтенд).
- `sections[].style` — вид блока (`first | second | third`). **Бэкенд пока не принимает это поле** (строгий разбор → `400`), поэтому фронтенд шлёт его, только если выбран не `first`. Что нужно на бэке: принять и сохранить строку (проверять только формат `^[a-z]{1,16}$`), хранить в колонке `style text NOT NULL DEFAULT 'first'`, отдавать в `GET /lk/resume/:id` и `GET /cv/:linkId`; пустое трактовать как `first`.

### Кабинет: ссылки

| Метод и путь | Тело | Ответ |
|---|---|---|
| `GET /lk/links?resumeId=…` | — | массив `{ id, name, status, version, createdAt, updatedAt }` |
| `POST /lk/link` | `{ resumeId, name }` | 201 `{ id }` |
| `PUT /lk/link/:linkId` | `{ version, name }` | 204 (переименовать) |
| `PUT /lk/link/:linkId/status` | `{ version, status }` | 204 |
| `DELETE /lk/link/:linkId` | — | 204 |

Адрес публичного резюме — `/cv/<linkId>`: каждая ссылка считает свои визиты отдельно.

### Публичные

`GET /cv/:linkId` — публичная проекция резюме (без авторизации).

- Отдаётся только опубликованное резюме (`opened`) по активной ссылке: неизвестная ссылка или неопубликованное резюме — `404`, отключённая ссылка — `410`, не UUID — `400`.
- В теле нет `id`, `status`, `version`; скрытые контакты не приходят, их ключи перечислены в `contacts.hiddenContacts`. У проектов нет `role`, `gitLink`, `contributors`.
- Заголовки: `Cache-Control: no-store`; `X-View-Id` — идентификатор визита (если удалось его создать; без него трекер не запускается); при необходимости ставится cookie посетителя.
- Заголовок запроса `X-Referrer-Host` — источник визита (хост страницы-источника); `Referer` у запроса всегда адрес самого сайта.

`POST /cv/:linkId/contact-views` — раскрыть скрытый контакт или получить все для файла.

```json
{ "initiator": "contact", "contact": "email", "viewId": "…", "turnstileToken": "…" }
```
Ответ: `{ "contact": "email", "value": "…" }`.

```json
{ "initiator": "download", "viewId": "…", "turnstileToken": "…" }
```
Ответ: `{ "contacts": [{ "contact": "email", "value": "…" }] }` — только непустые. Нужен одноразовый токен Cloudflare Turnstile (`action: contact_view`), если проверка включена. Значения контактов в аналитику не попадают.

## Аналитика

`POST /cv/:linkId/events` — пачка событий посетителя; ответ `204`.

```json
{
  "viewId": "…", "seq": 3, "final": false,
  "state": { "activeMs": 41000, "maxSection": "projects", "exitSection": "skills" },
  "events": [{ "id": "…", "type": "section_view", "at": 1730000000000 }],
  "clicks": [{ "section": "projects", "x": 4120, "y": 1800, "vw": 1280, "target": "link", "at": 1730000000000 }]
}
```

Типы событий: `section_view`, `download`, `tech_filter`, `external_click`, `print`, `contact_copy`, `contact_click`. Фронтенд шлёт пачки по таймеру и при уходе со страницы (`final: true`); посетитель анонимный, трекер уважает DNT/GPC, если включено `RESPECT_DO_NOT_TRACK`.

Отчёты для владельца (все требуют сессии). Периодные отчёты принимают `from`, `to` (даты) и `tz` (часовой пояс); фронтенд предлагает 7, 30 и 90 дней.

| Путь | Что отдаёт |
|---|---|
| `GET /lk/link/:linkId/statistics` | сводка: просмотры, прочтения, уникальные и вернувшиеся, раскрытия контактов, скачивания, устройства, источники, отфильтрованный трафик (период) |
| `GET /lk/link/:linkId/statistics/daily` | точки по дням (период) |
| `GET /lk/link/:linkId/statistics/sections` | сколько дошло до каждой секции и сколько ушло с неё (период) |
| `GET /lk/link/:linkId/statistics/technologies` | интерес к технологиям: выборы в фильтре (период) |
| `GET /lk/link/:linkId/views` | лента визитов, `cursor`, `limit`, `includeFiltered` |
| `GET /lk/link/:linkId/clicks` | тепловая карта кликов: `section`, `breakpoint` (`xs md lg`), `includeStale` |
| `GET /lk/resume/:resumeId/statistics` | сводка по резюме со списком ссылок (период) |

Бот-трафик, визиты владельца и подозрительные всплески помечаются и по умолчанию в сводках не учитываются.
