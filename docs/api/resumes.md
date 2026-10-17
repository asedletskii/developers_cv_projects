# Резюме

Само резюме — контейнер с названием и языком. Содержимое (обзор, профиль, списки) живёт в своих блоках и редактируется отдельными ручками.
Фронтенд: `getResumes`, `createResume`, `getResume`, `updateTitle`, `deleteResume` → `LkResumesScreen`, `ResumeList`, `CreateResumeDialog`, `LkResumeEditScreen`.

## `GET /lk/resumes`

Список резюме пользователя. Тела нет.

**Ответ `200`** — массив:

| Поле | Тип | Описание |
|---|---|---|
| `id` | uuid | идентификатор резюме |
| `title` | string | название |
| `language` | string | язык резюме (см. справочник в [README](README.md#справочники-значений)) |

Пустой список — `[]`. Порядок не гарантирован.

## `POST /lk/resume`

Создать резюме. **Ответ `201 {"id": "<uuid>"}`.**

| Поле | Тип | Обяз. | Ограничения | Дефолт |
|---|---|---|---|---|
| `title` | string | да | после `trim` от 1 до 200 символов | — |
| `language` | string | да | одно из `ru en` | — |

**Ограничения уровня пользователя:** одно резюме на язык — повторный `language` у того же пользователя → `409`. Язык потом не меняется.

**Что создаётся вместе с резюме (дефолты):**

| Часть | Состояние сразу после создания |
|---|---|
| `createdAt`, `updatedAt` | текущее время |
| Обзор (`overview`) | создаётся: `greeting` — четыре пустые строки, `sections` — `[]`; `theme`, `availability` отсутствуют |
| Профиль, контакты, навыки | **не создаются**: в `GET /lk/resume/:id` полей `profile`, `contacts`, `skills` нет, пока их не создадут через `POST` |
| Языки, технологии, опыт, образование, курсы, активности, проекты | пустые массивы `[]` |

Пустые `sections` фронтенд трактует как «показать секции по умолчанию» (см. [overview.md](overview.md)).

## `GET /lk/resume/:resumeId`

Резюме владельца целиком — данные для редактора. Тела нет. **Ответ `200`:**

| Поле | Тип | Описание |
|---|---|---|
| `id`, `title`, `language` | | как в списке |
| `createdAt`, `updatedAt` | datetime | RFC 3339, UTC |
| `overview` | object | оформление, всегда есть → [overview.md](overview.md) (поля `theme`, `greeting`, `availability`, `sections` лежат прямо в нём, без вложенного `presentation`; `resumeLink` убран) |
| `profile` | object | есть только после создания → [personal.md](personal.md) |
| `contacts` | object | есть только после создания → [personal.md](personal.md) |
| `skills` | object | есть только после создания → [personal.md](personal.md) |
| `languages`, `technologies` | array | [languages-technologies.md](languages-technologies.md) |
| `experiences` | array | [experience.md](experience.md) |
| `educations`, `courses` | array | [education.md](education.md) |
| `projects` | array | [projects.md](projects.md) |
| `activities` | array | [activities.md](activities.md) |

Элементы массивов содержат собственные `id` и `resumeId` — по ним фронтенд делает `PUT`/`DELETE`. Пустые списки приходят как `[]`.

**Ошибки:** `400` — `resumeId` не UUID; `404` — не найдено или чужое.

## `PUT /lk/resume/:resumeId`

Переименовать. **Ответ `204`.**

| Поле | Тип | Обяз. | Ограничения | Дефолт |
|---|---|---|---|---|
| `title` | string | да | после `trim` от 1 до 200 символов | — |

Других полей нет (язык менять нельзя, лишние поля → `400`). Обновляет `updatedAt`.

## `DELETE /lk/resume/:resumeId`

Удалить резюме со всем содержимым (каскад: обзор, профиль, контакты, навыки, все списки). **Ответ `204`.** Необратимо; публичный адрес после этого отвечает `404`.
