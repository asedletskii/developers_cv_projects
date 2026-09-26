# Дополнительная активность

Список активностей: опенсорс, хакатоны, выступления, наставничество и т. п. Фронтенд: вкладка «Активности», записи сопоставляются **по `id`** (новые — `POST`, известные — `PUT`, пропавшие — `DELETE`).

| Метод и путь | Что делает | Ответ |
|---|---|---|
| `POST /lk/resume/:resumeId/activities` | добавить в конец | `201 {"id"}` |
| `GET /lk/resume/:resumeId/activities` | список | `200 [ … ]` |
| `GET /lk/resume/:resumeId/activities/:activityId` | одна запись | `200`; `404` |
| `PUT /lk/resume/:resumeId/activities/:activityId` | заменить запись | `204`; `404` |
| `DELETE /lk/resume/:resumeId/activities/:activityId` | удалить | `204`; `404` |

Порядок — по времени создания, переставить нельзя.

## Тело `POST` / `PUT`

| Поле | Тип | Обяз. | Ограничения | Дефолт при создании |
|---|---|---|---|---|
| `type` | string | да | `openSource hackathon speaker volunteering mentoring community` | — |
| `title` | string | да | непустое | — |
| `description` | string | нет | свободный текст | отсутствует |
| `date` | string | нет | `YYYY-MM` | отсутствует |
| `link` | string | нет | `http`/`https`-ссылка | отсутствует |

Ответ `GET` содержит те же поля плюс `id` и `resumeId`.
