# Что поправить на бэкенде под контракт фронта

Источник правды — фронт (типы в `frontend/src/features/curriculumVitae/types/`, описание в
[`backend-api.md`](./backend-api.md) и [`lk-api.md`](./lk-api.md)). Исключение — пути ручек:
они остаются как сейчас на бэке, фронт подстроится сам. Авторизацию (`X-Mock-User-ID`) пока не
трогаем.

Всё ниже касается тела `POST /api/v1/lk/resume` и `PUT /api/v1/lk/resume/{resumeId}` (оба сейчас
парсятся через `createResumeRequest` в `resume/transport/http/create.go`), плюс одно место в `GET`.

## 1. Тексты — обычные строки, без локализации

Сейчас бэк ждёт `{ "ru": "...", "en": "..." }` в полях:

- `greeting.first`, `greeting.second`, `greeting.secondHighlight`, `greeting.third`
- `availability.note`
- `contactsText.title`, `contactsText.subtitle`

Фронт шлёт обычные строки — в том же виде, в каком получил их из `GET`. Сейчас это всегда
`400 Invalid request body`.

**Нужно:**

- В этих полях принимать и отдавать `string`, а не объект.
- Локализация контента резюме (`?lang`, `ru`/`en` в хранилище) — рудимент: фронт `lang` больше
  нигде не передаёт (ни в `GET`, ни в `POST`/`PUT`). У каждого текстового поля одно значение.

Пример тела `PUT /api/v1/lk/resume/{id}`:

```json
{
  "greeting": { "first": "Hello!", "second": "My name is", "secondHighlight": "Anton", "third": "" },
  "availability": { "status": "open", "from": "2026-10-01", "note": "Remote only" },
  "contactsText": { "title": "", "subtitle": "Drop me a line" }
}
```

## 2. Пустые строки

Форма на фронте присылает `""` для незаполненных полей, а не опускает их. Сейчас бэк пытается
провалидировать `""` и отвечает 400 (например, пустой `endDate` падает в `ParseYearMonth("")`,
пустой `resumeLink` — в `NewAssetPath("")`).

**Особые поля — у `""` есть смысл, это не «нет значения»:**

| Поле                                   | `""`                     | Поле отсутствует / `null`   |
| -------------------------------------- | ------------------------ | ---------------------------- |
| `greeting.*` (все четыре, обязательные) | Строка скрыта на странице | — (поле обязательное)        |
| `contactsText.title`, `.subtitle`      | Элемент скрыт            | Текст по умолчанию           |

То есть для `greeting.*` пустая строка допустима (сейчас `NewLocalizedText` её не пропускает), а
для `contactsText.*` бэк должен хранить `""` отдельно от `null` и отдавать `""` в `GET` как есть.

**`availability.status: ""`** — у резюме нет плашки доступности: весь `availability` считать
`null` (фронт так передаёт вариант «—» в селекте).

**Все остальные необязательные поля: `""` = поле не задано (`null`).** Полный список:

- `patronymic`, `resumeLink`
- `availability.from`, `availability.note`
- `theme.primary`, `theme.secondary`
- `contacts.phone`, `contacts.telegram`, `contacts.email`
- `languages[].grade`
- `experience[].location`, `experience[].endDate`, `experience[].description`
- `education[].profileName`, `education[].endDate`
- `courses[].certificateLink`
- `activities[].description`, `activities[].date`, `activities[].link`
- `projects[].description`, `projects[].link`

## 3. `photoPath` и `resumeLink` — принимать и полный URL

Сейчас оба проходят через `NewAssetPath` и принимают только путь на сайте (`/...`, не `//...`).
По контракту фронта это ссылка: пользователь в форме может вставить и путь, и полный адрес
(`https://...`).

**Нужно:** принимать либо путь, начинающийся с `/` (но не `//`), либо абсолютный `http(s)` URL.

## 4. Проекты: `PUT` не должен затирать поля, которых нет у фронта

На фронте у проекта в резюме только `id`, `name`, `description`, `stack`, `link` — так их и отдаёт
`GET` (`ProjectSummary`). Бэк при этом хранит ещё `role`, `screenshots`, `gitLink`,
`contributors`, а `PUT` делает полную замену резюме. Каждое сохранение из редактора стирает эти
четыре поля у всех проектов.

**Нужно** одно из двух:

- при `PUT` для проекта с уже существующим `id` сохранять `role`, `screenshots`, `gitLink`,
  `contributors` как есть, если их нет в теле; новым проектам (нет `id` или `id` не найден) —
  пустые значения;
- или убрать эти поля из модели/API, если они не нужны (на фронте они сейчас нигде не
  используются).

В `POST` эти поля не обязательны: фронт их не шлёт.

## 5. Лишние поля в теле — игнорировать

Фронт отправляет резюме целиком, как получил его из `GET`, поэтому в теле могут быть `id`,
`title` и `updatedAt`. Их надо молча игнорировать: `id` берётся из пути (для `POST` генерируется
новый), `title` меняется только через `PATCH`, `updatedAt` ставит сервер. Сейчас так и работает
(в `createResumeRequest` этих полей нет) — важно не сломать.

## Что уже совпадает, менять не надо

- Пути ручек (фронт переедет на них сам): `GET /lk/profile`, `GET /lk/resumes`,
  `GET|PUT|PATCH /lk/resume/{id}`, `POST /lk/resume`, `GET /cv/{id}`.
- Форматы ответов `GET /lk/profile`, `GET /lk/resumes`, `PATCH` (`{ id, title }`).
- Название нового резюме по умолчанию = `position`; `PUT` не трогает `title`.
- Форматы дат: `YYYY-MM-DD` (`dateOfBirth`, `availability.from`), `YYYY-MM` (опыт, курсы,
  активности), `YYYY` (`education[].endDate`).
- Внешние ссылки (`link`, `gitLink`, `certificateLink`) — только `http(s)`; email; telegram
  (буквы, цифры, `_`, `@` в начале срезается); цвета `#RGB` / `#RRGGBB`.
- `id` вложенных записей — UUID: существующие фронт берёт из `GET`, новые генерирует сам через
  `crypto.randomUUID()`.
- У `Technology` только `name`, `grade` нет.
- Обязательные поля резюме: `name`, `surname`, `position`, `location`, `photoPath`,
  `dateOfBirth`.
