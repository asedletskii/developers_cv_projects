# Публичная проекция

Плоское представление резюме для страницы `/cv/:resumeId`. Авторизации нет — резюме открыто любому, кто знает адрес.
Фронтенд: `getPublicResume` → `CurriculumVitaeScreen` и предпросмотр в кабинете (`LkResumeScreen`).

## `GET /cv/:resumeId`

Тела нет. **Ответ `200`** — собранная из всех секций плоская структура (тип `Person` на фронтенде). Записать через эту ручку ничего нельзя.

| Поле | Тип | Откуда |
|---|---|---|
| `id`, `title`, `language` | | резюме |
| `updatedAt` | datetime | резюме (может отсутствовать) |
| `name`, `surname`, `patronymic`, `position`, `photoPath`, `dateOfBirth`, `location`, `relocation`, `workFormat` | | [профиль](personal.md) |
| `contacts` (`phone`, `telegram`, `email`), `hiddenContacts`, `contactsText` | | [контакты](personal.md) |
| `hardSkills`, `softSkills` | string[] | [навыки](personal.md) |
| `theme`, `greeting`, `availability`, `sections` | | [обзор](overview.md) (без обёртки `presentation`); `resumeLink` убрано — фронтенд его не использует |
| `languages` | `{ language, grade? }[]` | [языки](languages-technologies.md), по порядку |
| `technologies` | `{ name }[]` | [стэк](languages-technologies.md), по порядку |
| `experience` | | [опыт](experience.md); у записи `stack` (`{name}[]`) вместо `technologies` |
| `education`, `courses`, `activities` | | [образование](education.md), [активности](activities.md) |
| `projects` | `{ id, name, description?, stack, link? }[]` | [проекты](projects.md) — только краткая часть |

**Отличия от ответа `GET /lk/resume/:id`:**
- имена: `experience`/`education` вместо `experiences`/`educations`, у опыта `stack` вместо `technologies`;
- нет `resumeId` у записей, у технологий и стэков нет `grade`, у проектов нет `role`, `screenshots`, `gitLink`, `contributors`;
- профиль, контакты и навыки развёрнуты на верхний уровень; если их не создавали, поля пусты (`""`, `[]`, `{}`), а не отсутствуют;
- массивы приходят как `[]`, но фронтенд дополнительно превращает `null` в `[]` (`normalizePerson`).

**Ошибки:** `400` — `resumeId` не UUID; `404` — резюме нет.

**Аналитика.** Этот же запрос создаёт визит: в ответе должны быть заголовок `X-View-Id` и cookie посетителя — см. [analytics.md](analytics.md#get-cvlinkid--выдача-резюме-и-визит).

**Скрытые контакты.** Сейчас `contacts` приходят в открытом виде, а `hiddenContacts` лишь подсказывает интерфейсу не показывать их сразу — парсер прочтёт всё из JSON. Планируется отдавать скрытые значения только отдельным запросом с проверкой на бота (см. обсуждение защиты от парсинга).

## `POST /cv/:linkId/contact-views` — раскрытие контактов и Turnstile

Для раскрытия одного скрытого контакта тело содержит `initiator: "contact"` и ключ контакта:

```json
{
  "initiator": "contact",
  "contact": "email",
  "viewId": "8c0e…",
  "turnstileToken": "..."
}
```

Ответ: `{ "contact": "email", "value": "person@example.com" }`.

Для скачивания файл собирается на фронтенде из публичного резюме и всех контактов. Запрос отличается тем, что `contact` не передаётся:

```json
{
  "initiator": "download",
  "viewId": "8c0e…",
  "turnstileToken": "..."
}
```

Ответ содержит только заполненные контакты, включая скрытые в публичном представлении:
`{ "contacts": { "phone": "...", "telegram": "...", "email": "..." } }`.
Успешная выдача при корректном `viewId`, относящемся к этой ссылке, также записывает событие `download`; статистика скачиваний уже учитывает этот тип события.

`initiator` принимает `contact` или `download`. Для совместимости с прежним клиентом поле можно не передавать, если указан `contact`; запрос без `initiator` и без контакта некорректен. При `initiator: "contact"` контакт обязателен, а при `initiator: "download"` он должен отсутствовать.

Оба сценария требуют активную ссылку на опубликованное резюме и проходят общую проверку Turnstile, если она включена на сервере. Фронтенд перед каждым действием запрашивает у Cloudflare новый токен (режим `execute`, `interaction-only`: задача показывается, только если Cloudflare сочтёт её нужной) и отправляет его в `turnstileToken`. Токен **одноразовый** и живёт 300 секунд.
