# API резюме: контракт для бэкенда

Источник правды по типам: [`src/features/curriculumVitae/types/`](../src/features/curriculumVitae/types/)
(`curriculumVitae.ts` — поля, `enums.ts` — перечисления). Пример ответа: мок
[`curriculumVitae.mock.ts`](../src/features/curriculumVitae/api/curriculumVitae.mock.ts).

## Эндпоинт

`GET /cv/{cvId}?lang={ru|en}` → `Person`

`cvId` is a UUID.

| Код   | Результат                                         |
| ----- | ------------------------------------------------- |
| `200` | Резюме                                            |
| `404` | Нет такого резюме (страница «не найдено»)         |
| иное  | Страница «не удалось загрузить»                   |

- `Person.id` = `cvId` из URL.
- Тексты (свободные поля) сервер отдаёт на языке `lang`. Enum-значения — только коды, подписи на фронте.
- Публичный эндпоинт не должен отвечать `401` (клиент на 401 дёргает `/auth/refresh`).
- CORS: конкретный origin + `Access-Control-Allow-Credentials: true` (axios с `withCredentials`).

## Person

`✓` обязательное, `—` необязательное.

| Поле             | Тип                                 |  | Примечание                                                                  |
| ---------------- | ----------------------------------- | - | --------------------------------------------------------------------------- |
| `id`             | string                              | ✓ |                                                                             |
| `name`           | string                              | ✓ |                                                                             |
| `surname`        | string                              | ✓ |                                                                             |
| `patronymic`     | string                              | — |                                                                             |
| `position`       | string                              | ✓ |                                                                             |
| `photoPath`      | string (URL)                        | ✓ | Не загрузится — фронт покажет инициалы                                      |
| `dateOfBirth`    | `YYYY-MM-DD`                        | ✓ | На странице показывается только возраст                                     |
| `location`       | string                              | ✓ |                                                                             |
| `relocation`     | boolean                             | ✓ |                                                                             |
| `workFormat`     | `WorkFormat[]`                      | ✓ |                                                                             |
| `languages`      | `{ language: LanguageType; grade?: LanguageGrade }[]` | ✓ |                                                           |
| `contacts`       | `Contacts`                          | ✓ |                                                                             |
| `hiddenContacts` | `ContactKey[]`                      | — | Контакты, раскрываемые по клику                                             |
| `contactsText`   | `ContactsText`                      | — |                                                                             |
| `greeting`       | `Greeting`                          | ✓ |                                                                             |
| `availability`   | `Availability`                      | — | Плашка в hero                                                               |
| `updatedAt`      | ISO 8601                            | — | Показывается только вместе с `availability`                                 |
| `resumeLink`     | string (URL)                        | — | Файл резюме. Нет поля — кнопок скачивания нет                               |
| `theme`          | `ThemeSettings`                     | — |                                                                             |
| `sections`       | `SectionSettings[]`                 | ✓ |                                                                             |
| `technologies`   | `Technology[]`                      | ✓ | Порядок = порядок в ленте                                                   |
| `hardSkills`     | `HardSkill[]`                       | ✓ |                                                                             |
| `softSkills`     | `SoftSkill[]`                       | ✓ |                                                                             |
| `experience`     | `Experience[]`                      | ✓ | Новые → старые                                                              |
| `education`      | `Education[]`                       | ✓ |                                                                             |
| `courses`        | `Course[]`                          | ✓ |                                                                             |
| `activities`     | `Activity[]`                        | ✓ |                                                                             |
| `projects`       | `Project[]`                         | ✓ | Краткая версия, поля ниже                                                   |

Массивы всегда отдавать, пусть пустыми. Пустая по данным секция на странице не показывается.

### Contacts

| Поле       | Тип    |  | Примечание                              |
| ---------- | ------ | - | --------------------------------------- |
| `telegram` | string | — | Ник без `@`                             |
| `email`    | string | — |                                         |
| `phone`    | string | — | Формат `+79999999999` (идёт в `tel:`)   |

### ContactsText

| Поле       | Тип    |  | Примечание                                                          |
| ---------- | ------ | - | ------------------------------------------------------------------- |
| `title`    | string | — | Нет поля — текст по умолчанию; `""` — заголовок скрыт               |
| `subtitle` | string | — | Так же                                                              |

### Greeting

Все четыре поля обязательны; пустая строка скрывает строку (вторая скрывается, если пусты оба поля).

| Поле              | Тип    | Примечание                        |
| ----------------- | ------ | --------------------------------- |
| `first`           | string | 1-я строка                        |
| `second`          | string | 2-я строка, обычный текст         |
| `secondHighlight` | string | 2-я строка, выделенная часть      |
| `third`           | string | 3-я строка                        |

### Availability

| Поле     | Тип                  |  | Примечание                     |
| -------- | -------------------- | - | ------------------------------ |
| `status` | `AvailabilityStatus` | ✓ |                                |
| `from`   | `YYYY-MM-DD`         | — | С какой даты можно выйти       |
| `note`   | string               | — |                                |

### ThemeSettings

Цвета только `#RGB` / `#RRGGBB`, иначе игнорируются.

| Поле        | Тип    |  |
| ----------- | ------ | - |
| `primary`   | string | — |
| `secondary` | string | — |

### SectionSettings

Порядок массива = порядок на странице. Секции, которых нет в массиве, не показываются. Пустой или
отсутствующий массив — все секции в порядке `SectionId`.

| Поле        | Тип         |  | Примечание                                                                        |
| ----------- | ----------- | - | --------------------------------------------------------------------------------- |
| `id`        | `SectionId` | ✓ |                                                                                   |
| `enabled`   | boolean     | ✓ |                                                                                   |
| `collapsed` | boolean     | — | Только для `skills`, `education`, `courses`, `projects`, `activities`             |

### Technology

Названия должны совпадать в `technologies`, `experience[].stack`, `projects[].stack` (по ним работает
фильтр; регистр не важен).

| Поле    | Тип                | |
| ------- | ------------------ | - |
| `name`  | string             | ✓ |
| `grade` | `TechnologyGrade`  | — |

### Experience

| Поле           | Тип            |  | Примечание                              |
| -------------- | -------------- | - | --------------------------------------- |
| `id`           | string         | ✓ |                                         |
| `company`      | string         | ✓ |                                         |
| `position`     | string         | ✓ |                                         |
| `location`     | string         | — |                                         |
| `startDate`    | `YYYY-MM`      | ✓ |                                         |
| `endDate`      | `YYYY-MM`      | — | Нет поля — «по настоящее время»         |
| `description`  | string         | — |                                         |
| `achievements` | `string[]`     | ✓ |                                         |
| `stack`        | `Technology[]` | ✓ |                                         |

### Education

| Поле          | Тип              |  | Примечание                              |
| ------------- | ---------------- | - | --------------------------------------- |
| `id`          | string           | ✓ |                                         |
| `university`  | string           | ✓ |                                         |
| `location`    | string           | ✓ |                                         |
| `grade`       | `EducationLevel` | ✓ | Число                                   |
| `programName` | string           | ✓ |                                         |
| `profileName` | string           | — |                                         |
| `faculty`     | string           | ✓ |                                         |
| `endDate`     | `YYYY`           | — | Год окончания                           |
| `course`      | number           | — | Показывается, если нет `endDate`        |

### Course

| Поле              | Тип          |  |
| ----------------- | ------------ | - |
| `id`              | string       | ✓ |
| `name`            | string       | ✓ |
| `organization`    | string       | ✓ |
| `endDate`         | `YYYY-MM`    | ✓ |
| `hours`           | number       | — |
| `certificateLink` | string (URL) | — |

### Activity

| Поле          | Тип            |  |
| ------------- | -------------- | - |
| `id`          | string         | ✓ |
| `type`        | `ActivityType` | ✓ |
| `title`       | string         | ✓ |
| `description` | string         | — |
| `date`        | `YYYY-MM`      | — |
| `link`        | string (URL)   | — |

### Project

| Поле          | Тип            |  |
| ------------- | -------------- | - |
| `id`          | string         | ✓ |
| `name`        | string         | ✓ |
| `description` | string         | — |
| `stack`       | `Technology[]` | ✓ |
| `link`        | string (URL)   | — |

## Enums

Сервер отдаёт значения (правая колонка). Неизвестное значение фронт не отобразит.

| Enum                 | Значения                                                                                      |
| -------------------- | --------------------------------------------------------------------------------------------- |
| `WorkFormat`         | `fullTime` `partTime` `office` `remote` `hybrid`                                              |
| `LanguageType`       | `ru` `en` `fr` `de` `kz` `zh` `es`                                                            |
| `LanguageGrade`      | `A1` `A2` `B1` `B2` `C1` `C2`                                                                 |
| `HardSkill`          | `architecture` `testing` `cicd` `performance` `accessibility` `codeReview` `uiDesign` `api`   |
| `SoftSkill`          | `communication` `teamwork` `problemSolving` `timeManagement` `mentoring` `adaptability` `ownership` |
| `TechnologyGrade`    | `intern` `junior` `middle` `senior`                                                           |
| `EducationLevel`     | `0` среднее, `1` бакалавр, `2` специалист, `3` магистр, `4` аспирант, `5` кандидат наук, `6` доктор наук |
| `ContactKey`         | `telegram` `email` `phone`                                                                    |
| `ActivityType`       | `openSource` `hackathon` `speaker` `volunteering` `mentoring` `community`                     |
| `AvailabilityStatus` | `open` `considering` `busy`                                                                   |
| `SectionId`          | `about` `techStack` `experience` `skills` `education` `courses` `projects` `activities` `contacts` |

## Прочее

- `resumeLink` на другом домене: чтобы файл скачивался, а не открывался, нужен `Content-Disposition: attachment`.
- Внешние ссылки (`link`, `certificateLink`, `photoPath`, `resumeLink`) фронт не валидирует: проверять схему `http(s)` на сервере.
- Фронт пока работает на моках. Включить реальный API — фронт заменит вызов в thunk `getPersonData` (`curriculumVitae.thunks.ts`); базовый URL зашит в `src/shared/api/axios.config.ts`.

## Деплой (CD)

CI (`.github/workflows/ci.yml`) при пуше в `main`: собирает и публикует образы `frontend` и `backend`
в GHCR, затем по SSH обновляет стек на сервере (`docker-compose.prod.yml` из корня репозитория).

**Топология:** Caddy (образ `frontend`) — единственная точка входа, отдаёт статику фронта и
проксирует `/api/*` на `backend:8080` внутри docker-сети. `backend`, `postgres`, `minio`, `nats`
наружу не торчат.

**Одноразовая настройка сервера:**

1. Установить Docker и Docker Compose.
2. Создать директорию деплоя (например `/opt/developers-cv`) и положить туда `.env` на основе
   [`.env.example`](../.env.example) из корня репозитория — с реальными паролями и, если есть,
   доменом в `FRONTEND_ADDRESS` (тогда Caddy сам поднимет HTTPS).
3. Если пакеты в GHCR приватные — авторизовать Docker на сервере одной командой:
   `docker login ghcr.io -u <username> -p <personal access token со scope read:packages>`.
   Проще сделать пакеты `frontend`/`backend` публичными в настройках GHCR и пропустить этот шаг.

**Секреты в GitHub (Settings → Secrets and variables → Actions), нужны для job `deploy`:**

| Секрет           | Значение                                                    |
| ---------------- | ------------------------------------------------------------ |
| `DEPLOY_HOST`    | IP или домен сервера                                          |
| `DEPLOY_USER`    | Пользователь SSH                                               |
| `DEPLOY_SSH_KEY` | Приватный ключ SSH (без пароля)                                |
| `DEPLOY_PORT`    | Порт SSH, необязательно (по умолчанию 22)                      |
| `DEPLOY_PATH`    | Директория на сервере с `.env` из шага 2 (например `/opt/developers-cv`) |

Без этих секретов джобы `docker`/`frontend-docker` (сборка и паблиш образов) всё равно отработают,
упадёт только `deploy`.
