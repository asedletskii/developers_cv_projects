# Проекты

Список проектов резюме. Фронтенд: вкладка «Редактирование → Проекты» (`ProjectsEditor`), записи сопоставляются **по `id`** (новые — `POST`, известные — `PUT`, пропавшие — `DELETE`).

**Важно для клиентов.** `PUT` заменяет проект целиком, включая `role`, `screenshots`, `gitLink` и `contributors`. Форма фронтенда правит только название, описание, ссылку и стэк, поэтому перед `PUT` подмешивает остальные поля из последнего `GET /lk/resume/:id` — иначе они бы стёрлись.

| Метод и путь | Что делает | Ответ |
|---|---|---|
| `POST /lk/resume/:resumeId/projects` | добавить в конец | `201 {"id"}` |
| `GET /lk/resume/:resumeId/projects` | список | `200 [ … ]` |
| `GET /lk/resume/:resumeId/projects/:projectId` | один проект | `200`; `404` |
| `PUT /lk/resume/:resumeId/projects/:projectId` | заменить проект | `204`; `404` |
| `DELETE /lk/resume/:resumeId/projects/:projectId` | удалить | `204`; `404` |

Порядок — по времени создания, переставить нельзя.

## Тело `POST` / `PUT`

| Поле | Тип | Обяз. | Ограничения | Дефолт при создании |
|---|---|---|---|---|
| `name` | string | да | непустое | — |
| `description` | string | нет | свободный текст | отсутствует |
| `role` | string | нет | роль в проекте | отсутствует |
| `screenshots` | string[] | нет | каждый — `http`/`https`-ссылка | `[]` |
| `link` | string | нет | `http`/`https`-ссылка «где посмотреть» | отсутствует |
| `gitLink` | string | нет | `http`/`https`-ссылка на репозиторий | отсутствует |
| `stack` | object[] | нет | технологии проекта | `[]` |
| `stack[].name` | string | да | непустое | — |
| `stack[].grade` | string | нет | `intern junior middle senior` | отсутствует |
| `contributors` | object[] | нет | участники | `[]` |
| `contributors[].name` | string | да | непустое | — |
| `contributors[].surname` | string | да | непустое | — |
| `contributors[].role` | string | да | непустое | — |
| `contributors[].contacts` | object | да | нужен хотя бы один контакт из трёх, иначе `400` | — |
| `contacts.phone` | string | нет | свободная строка; пустая допустима, если задан `telegram` или `email` | отсутствует |
| `contacts.telegram` | string | нет | `a-z A-Z 0-9 _`, ведущий `@` отбрасывается | отсутствует |
| `contacts.email` | string | нет | корректный адрес | отсутствует |

**Ответ `GET`** содержит те же поля плюс `id` и `resumeId`; у участников в `contacts` все три поля необязательны (`omitempty`).

Фронтенд при `PUT` подставляет `phone: ""`, если у участника телефона нет (ключ в этом DTO — строка, а не необязательное поле).
