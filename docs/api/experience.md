# Опыт работы

Список мест работы. Фронтенд: вкладка «Редактирование → Опыт» (`ExperienceEditor`), записи сопоставляются с серверными **по `id`**: новые (клиентский id, которого нет на сервере) создаются `POST`, известные обновляются `PUT`, пропавшие удаляются `DELETE`.

| Метод и путь | Что делает | Ответ |
|---|---|---|
| `POST /lk/resume/:resumeId/experiences` | добавить в конец | `201 {"id"}` |
| `GET /lk/resume/:resumeId/experiences` | список | `200 [ … ]` |
| `GET /lk/resume/:resumeId/experiences/:experienceId` | одна запись | `200`; `404` |
| `PUT /lk/resume/:resumeId/experiences/:experienceId` | заменить запись | `204`; `404` |
| `DELETE /lk/resume/:resumeId/experiences/:experienceId` | удалить | `204`; `404` |

## Тело `POST` / `PUT`

| Поле | Тип | Обяз. | Ограничения | Дефолт при создании |
|---|---|---|---|---|
| `company` | string | да | непустое | — |
| `position` | string | да | непустое | — |
| `location` | string | нет | город или формат («Удалённо») | отсутствует |
| `startDate` | string | да | `YYYY-MM` | — |
| `endDate` | string | нет | `YYYY-MM`, не раньше `startDate`; нет = «по настоящее время» | отсутствует |
| `description` | string | нет | свободный текст | отсутствует |
| `achievements` | string[] | нет | каждый элемент непустой (пустые строки → `400`, фронтенд их отфильтровывает) | `[]` |
| `technologies` | object[] | нет | см. ниже | `[]` |
| `technologies[].name` | string | да | непустое | — |
| `technologies[].grade` | string | нет | `intern junior middle senior` | отсутствует |

**Порядок.** Записи идут в порядке создания (внутренний `sort_order`: новая = максимум + 1). Переставить нельзя.

**Ответ `GET`** содержит те же поля плюс `id` и `resumeId`; даты в форматах, как при записи.
