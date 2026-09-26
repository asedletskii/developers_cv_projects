# План: один вид ответа для публичной страницы и кабинета

## Итог

**Итоговая форма — вложенная.** Публичный `GET /cv/:resumeId` начинает возвращать ту же вложенную структуру `ResumeView`, что уже отдаёт кабинет `GET /lk/resume/:resumeId`. Плоская форма публичного ответа удаляется целиком. После работ в проекте одна модель данных резюме и один набор названий.

Ничего не остаётся «на выбор»: все решения ниже приняты, для каждого указано, что меняем на что.

## 1. Что меняем в контракте: было → стало

### 1.1 Переименования и структура

| # | Где | Было | Стало |
|---|---|---|---|
| 1 | Ответ `GET /cv/:id` | плоские `name`, `surname`, `patronymic`, `position`, `photoPath`, `dateOfBirth`, `location`, `relocation`, `workFormat` в корне | те же поля внутри объекта `profile` |
| 2 | Ответ `GET /cv/:id` | `contacts: {phone, telegram, email}`, `hiddenContacts`, `contactsText` в корне | всё внутри объекта `contacts`: `contacts.phone`, `contacts.telegram`, `contacts.email`, `contacts.hiddenContacts`, `contacts.contactsText` |
| 3 | Ответ `GET /cv/:id` | `hardSkills`, `softSkills` в корне | `skills.hardSkills`, `skills.softSkills` |
| 4 | Ответ `GET /cv/:id` | `theme`, `greeting`, `availability`, `resumeLink`, `sections` в корне | `overview.theme`, `overview.greeting`, `overview.availability`, `overview.resumeLink`, `overview.sections` |
| 5 | Ответ `GET /cv/:id` | `experience` | `experiences` |
| 6 | Ответ `GET /cv/:id` | `education` | `educations` |
| 7 | Ответ `GET /cv/:id`, элементы `languages` | `{ language, grade? }` | `{ id, language, grade?, position }` (как в кабинете) |
| 8 | Ответ `GET /cv/:id`, элементы `technologies` | `{ name }` | `{ id, name, grade?, position }` (как в кабинете) |
| 9 | Ответ `GET /lk/resume/:id`, поле `experiences[].technologies` | `technologies` | `stack` |
| 10 | Тело `POST`/`PUT /lk/resume/:id/experiences` | `technologies: [{name, grade?}]` | `stack: [{name, grade?}]` |
| 11 | `GET /lk/resume/:id/experiences[/:id]` | `technologies` | `stack` |
| 12 | Тело и ответ `PUT`/`GET /lk/resume/:id/overview` | `{ "presentation": { theme, greeting, availability, resumeLink, sections } }` | `{ theme, greeting, availability, resumeLink, sections }` — без обёртки `presentation` |

Больше **ничего не переименовываем**: `photoPath`, `resumeLink`, `hiddenContacts`, `stack` у проектов, адреса ручек (`/experiences`, `/educations`, …), имена полей запросов остальных ручек остаются как есть.

### 1.2 Правила для пустых значений

| # | Было | Стало |
|---|---|---|
| 13 | Публично: нет профиля → `name: ""`, `photoPath: ""` и т. п.; нет контактов → `contacts: {}`; нет навыков → `hardSkills: []`, `softSkills: []` | Публично, как и в кабинете: нет профиля → поля `profile` нет; нет контактов → поля `contacts` нет; нет навыков → поля `skills` нет |
| 14 | В части списков возможен `null` | Все списки (`languages`, `technologies`, `experiences`, `educations`, `courses`, `activities`, `projects`, а также вложенные `stack`, `achievements`, `sections`, `hiddenContacts`, `hardSkills`, `softSkills`, `workFormat`, `screenshots`, `contributors`) всегда приходят массивом, пустой — `[]`, никогда `null` |

### 1.3 Чем публичный ответ отличается от кабинетного

Один тип, одинаковые ключи и вложенность. Отличий ровно два, оба делает функция `toPublic(view)` на бэкенде:

| # | Часть | Кабинет | Публично |
|---|---|---|---|
| 15 | `projects[]` | полные проекты | у каждого проекта `role` и `gitLink` отсутствуют, `screenshots: []`, `contributors: []`; остальные поля (`id`, `resumeId`, `name`, `description`, `link`, `stack`) как в кабинете |
| 16 | `contacts` (в будущем, вместе с защитой от парсера) | значения контактов | значения скрытых контактов не отдаются, остаются ключи в `contacts.hiddenContacts`; **в рамках этого плана не делаем**, поведение контактов не меняется |

Скриншоты и `gitLink` в публичном ответе **не отдаём** (как сейчас). `createdAt`, `updatedAt`, `resumeId` у записей публично отдаём — они безвредны, а лишние различия — лишний код.

## 2. Бэкенд: что делаем

Пути — от `backend/internal/modules/resume`. Порядок соответствует нумерации в таблице выше.

| Шаг | Действие | Файлы | Покрывает |
|---|---|---|---|
| B1 | Переименовать JSON-тег `technologies` → `stack` у опыта: вид ответа и тело запроса (`ExperienceRequest`, `TechnologyRequest` в опыте) | `application/queries/experience.go`, `transport/http/experiences/handler.go`, их тесты | 9, 10, 11 |
| B2 | Убрать обёртку `presentation`: `OverviewRequest`/`ToCommand` принимают поля прямо в корне тела, `ResumeOverviewView` отдаёт их прямо; вложенное поле `overview` в `ResumeView` уже без обёртки — оставить как есть | `transport/http/overview/request.go`, `transport/http/overview/handler.go`, `application/queries/overview.go`, тесты | 12 |
| B3 | Публичный ответ отдаёт `ResumeView`: удалить `PublicResumeView`, все вложенные публичные типы (`Contacts`, `ContactsText`, `Greeting`, `Availability`, `ThemeSettings`, `SectionSettings`, `ResumeLanguage`, `Technology`, `Experience`, `Education`, `Course`, `Activity`, `ProjectSummary`) и функцию `projectPublicResume`; `PublicResumeQueryHandler.Handle` возвращает `ResumeView`, пропущенный через `toPublic` | `application/queries/public_resume.go`, `transport/http/public/handler.go` | 1–8, 13, 15 |
| B4 | Написать `toPublic(view ResumeView) ResumeView`: для каждого проекта обнулить `Role`, `GitLink`, поставить `Screenshots = []`, `Contributors = []`; больше ничего не менять | `application/queries/public_resume.go` | 15 |
| B5 | Гарантировать `[]` вместо `null`: в конструировании `ResumeView` (reader `GetResume` и `GetPublic`) инициализировать все списки непустым срезом `make(T, 0)`; добавить тест на пустое резюме | `infrastructure/persistence/readers/resume.go`, `application/queries/resume_test.go` | 14 |
| B6 | Обновить тесты публичной ручки под новую форму | `transport/http/public/handler_test.go`, `application/queries/public_resume_test.go` | — |
| B7 | Обновить документацию | `docs/api/resumes.md`, `docs/api/public.md`, `docs/api/overview.md`, `docs/api/experience.md` | — |

Итог по коду: удаляется публичный DTO и проекция, добавляется `toPublic`.

## 3. Фронтенд: что делаем

Пути — от `frontend/src/features/curriculumVitae`. Тип `Person` (плоский) удаляется, вместо него единственный тип `Resume` — это переименованный `ResumeView` из `types/lk.ts`, перенесённый в `types/resume.ts`.

### 3.1 Доступ к данным: было → стало

Для страницы, предпросмотра и форм редактора.

| Было (`Person`) | Стало (`Resume`) |
|---|---|
| `name`, `surname`, `patronymic`, `position`, `photoPath`, `dateOfBirth`, `location`, `relocation`, `workFormat` | `profile.name`, `profile.surname`, `profile.patronymic`, `profile.position`, `profile.photoPath`, `profile.dateOfBirth`, `profile.location`, `profile.relocation`, `profile.workFormat` |
| `contacts.phone`, `contacts.telegram`, `contacts.email` | без изменений (`contacts.phone`, …) |
| `hiddenContacts` | `contacts.hiddenContacts` |
| `contactsText.title`, `contactsText.subtitle` | `contacts.contactsText.title`, `contacts.contactsText.subtitle` |
| `hardSkills`, `softSkills` | `skills.hardSkills`, `skills.softSkills` |
| `theme.primary`, `theme.secondary` | `overview.theme.primary`, `overview.theme.secondary` |
| `greeting.first`, `.second`, `.secondHighlight`, `.third` | `overview.greeting.first`, … |
| `availability.status`, `.from`, `.note` | `overview.availability.status`, … |
| `resumeLink` | `overview.resumeLink` |
| `sections` | `overview.sections` |
| `experience` | `experiences` (у записи `stack` — без изменений) |
| `education` | `educations` |
| `languages`, `technologies`, `courses`, `activities`, `projects` | без изменений (в элементах добавились `id`, `position` у языков и технологий) |
| `updatedAt` | без изменений |

**Отсутствующие секции.** `profile`, `contacts`, `skills` могут отсутствовать (правило 13). В одном месте — новом файле `constants/emptyResume.ts` — задаются `EMPTY_PROFILE`, `EMPTY_CONTACTS`, `EMPTY_SKILLS`; компоненты читают `resume.profile ?? EMPTY_PROFILE` и т. д. Своих значений по умолчанию в компонентах не заводим. Начальное состояние хранилища — резюме с пустым `overview` и без `profile`/`contacts`/`skills`.

### 3.2 Шаги

**Этап A — чтение.** Делается вместе с бэкендовыми шагами B3–B6.

| Шаг | Действие | Файлы |
|---|---|---|
| F1 | Ввести тип `Resume` (= прежний `ResumeView`, `experiences[].stack`); удалить `Person` | `types/lk.ts` → `types/resume.ts`, `types/curriculumVitae.ts`, `types/index.ts` |
| F2 | `getPublicResume` и `getResume` возвращают `Resume`; удалить `normalizePerson` и его тест | `api/curriculumVitae.api.ts`, `utils/normalizePerson/*`, `utils/index.ts` |
| F3 | Компоненты страницы читают по таблице 3.1 | `sections/sectionDefinitions.tsx`, `components/AboutSection/*`, `ContactsSection/*`, `SkillsSection/*`, `ExperienceSection/*`, `ProjectsSection/*`, `TechStackSlider/*`, `FilterBar/*`, `TechnologyFilterProvider/*`, `CurriculumVitae/*`, `HeroSection/*` (приветствие и доступность приходят из `overview`) |
| F4 | Хранилище и его состояние по `Resume` | `store/curriculumVitae/*`, `types/curriculumVitae.slice.ts` |
| F5 | Ввести константы пустых секций | `constants/emptyResume.ts`, `constants/index.ts` |

**Этап B — редактор.** После этапа A.

| Шаг | Действие | Файлы |
|---|---|---|
| F6 | Форма каждой вкладки — это `Resume`, её поля адресуются по таблице 3.1 (`register('overview.greeting.first')`, `register('profile.name')`, `register('experiences.0.company')`) | все редакторы в `screens/LkResumeEditScreen/`: `AboutEditor`, `HeroEditor`, `ThemeEditor`, `ColorField`, `ContactsEditor`, `SkillsEditor`, `LanguagesEditor`, `TechnologyListField`, `ExperienceEditor`, `EducationEditor`, `CoursesEditor`, `ActivitiesEditor`, `ProjectsEditor`, `SectionsEditor`, `DateField`, `SelectField`, `HeroPreview`, `SectionPreview`, `ThemePreviewDialog` |
| F7 | Удалить `viewToPerson`; `saveSection` читает значения из вложенных путей формы | `screens/LkResumeEditScreen/resumeMapping.ts`, `ResumeEditor.tsx`, `SectionForm.tsx`, `LkResumeEditScreen.tsx` |
| F8 | Тело `PUT /experiences` шлёт `stack` вместо `technologies`; `PUT /overview` шлёт поля без обёртки `presentation` | `screens/LkResumeEditScreen/resumeMapping.ts` |
| F9 | Обновить документацию, если что-то разошлось | `docs/api/*` |

## 4. Проверка

- [ ] `GET /cv/:id` и `GET /lk/resume/:id` возвращают один и тот же тип; различаются только проектами по пункту 15
- [ ] У пустого резюме нет `null`: все списки `[]`, `profile`, `contacts`, `skills` отсутствуют
- [ ] Страница и предпросмотр открываются у резюме без профиля (нет имени и фото), без ошибок
- [ ] Редактор создаёт и сохраняет каждую из 11 вкладок; после сохранения значения не теряются
- [ ] `PUT /experiences` принимает `stack`; `PUT /overview` принимает тело без `presentation`
- [ ] В коде нет `Person`, `viewToPerson`, `normalizePerson`, `PublicResumeView`, `projectPublicResume`
- [ ] `docs/api/*` описывает ответы так, как их отдаёт бэкенд; тесты бэкенда и фронтенда проходят

## 5. Порядок работ

1. Бэкенд B1, B2, B5 и соответствующие шаги фронтенда F8 (имена в запросах и обёртка) — без изменения публичного ответа.
2. Бэкенд B3, B4, B6 вместе с фронтендом F1–F5 (публичный ответ и чтение).
3. Фронтенд F6, F7 (формы редактора), затем F9.
