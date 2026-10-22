# Матрица трассировки

Связь пользовательских историй с требованиями, кодом и тестами. Колонки **FR** и **NFR** берутся из историй; колонки кода и тестов заполняются в PR, который реализует историю.

Пути бэкенда — от `backend/internal/modules/` (или `backend/internal/platform/` для `middleware`), фронтенда — от `frontend/src/features/curriculumVitae/` (или `frontend/src/` для `app/` и `shared/`). Тесты — имена файлов рядом с кодом.

Требования без истории (технические, эксплуатационные) трассируются через NFR и [FR-API](03-functional/README.md#fr-api-общие-правила-http-api).

## EP-ACC. Доступ к кабинету и профиль пользователя

| История | Название | Приоритет | Статус | FR | NFR | Бэкенд | Фронтенд | Тесты |
|---|---|---|---|---|---|---|---|---|
| [US-ACC-01](02-user-stories/EP-ACC.md#us-acc-01-регистрация-и-вход) | Регистрация и вход | Must | План | FR-ACC-001, FR-ACC-002, FR-ACC-003 | NFR-SEC-01, NFR-SEC-02, NFR-PRV-01 | — | — | — |
| [US-ACC-02](02-user-stories/EP-ACC.md#us-acc-02-открыть-кабинет) | Открыть кабинет | Must | Реализовано | FR-ACC-001, FR-ACC-010, FR-ACC-011 | — | `user/transport/http/profile`, `platform/http/middleware/auth.go` | `screens/LkProfileScreen`, `components/ProfileCard` | `user/transport/http/profile/handler_test.go`, `middleware/auth_test.go` |
| [US-ACC-03](02-user-stories/EP-ACC.md#us-acc-03-редактировать-профиль-на-двух-языках) | Редактировать профиль на двух языках | Should | Реализовано | FR-ACC-010 – FR-ACC-015 | — | `user/transport/http/profile`, `user/domain/profile.go` | `components/EditProfileDialog`, `utils/getLocalizedProfile` | `user/domain/profile_test.go`, `profile_repository_integration_test.go`, `getLocalizedProfile.test.ts` |
| [US-ACC-04](02-user-stories/EP-ACC.md#us-acc-04-изоляция-данных) | Изоляция данных | Must | Частично | FR-ACC-003, FR-RES-004, FR-LNK-004 | NFR-SEC-03 | `resume/infrastructure/persistence`, `resume_link/application/queries/statistics.go` | — | `ownership_write_integration_test.go`, `resume_link/transport/http/integration_test.go` |
| [US-ACC-05](02-user-stories/EP-ACC.md#us-acc-05-удалить-аккаунт) | Удалить аккаунт | Must | План | FR-ACC-020 | NFR-PRV-04 | — | — | — |

## EP-RES. Управление резюме

| История | Название | Приоритет | Статус | FR | NFR | Бэкенд | Фронтенд | Тесты |
|---|---|---|---|---|---|---|---|---|
| [US-RES-01](02-user-stories/EP-RES.md#us-res-01-список-резюме) | Список резюме | Must | Реализовано | FR-RES-001, FR-RES-060 | — | `resume/application/queries/list_resumes.go` | `screens/LkResumesScreen`, `components/ResumeList` | `list_resumes_test.go` |
| [US-RES-02](02-user-stories/EP-RES.md#us-res-02-создать-резюме) | Создать резюме | Must | Реализовано | FR-RES-002, FR-RES-003, FR-RES-005 | — | `resume/application/commands/resume.go` | `components/CreateResumeDialog` | `commands/resume_test.go`, `resume_handler_test.go` |
| [US-RES-03](02-user-stories/EP-RES.md#us-res-03-переименовать-резюме) | Переименовать резюме | Should | Реализовано | FR-RES-007, FR-RES-010 | — | `resume/transport/http/resume_handler.go` | `components/EditResumeTitleDialog`, `store/lk/lk.thunks.ts` | `resume_handler_test.go` |
| [US-RES-04](02-user-stories/EP-RES.md#us-res-04-опубликовать-резюме) | Опубликовать резюме | Must | Реализовано | FR-RES-020 – FR-RES-023 | — | `resume/domain/resume.go` (`Publish`) | `components/ResumeStatusChip` | `domain/resume_test.go`, `public_private_integration_test.go` |
| [US-RES-05](02-user-stories/EP-RES.md#us-res-05-снять-с-публикации) | Снять с публикации | Must | Реализовано | FR-RES-020 – FR-RES-023 | — | `resume/domain/resume.go` (`Close`) | `components/ResumeStatusChip` | `domain/resume_test.go` |
| [US-RES-06](02-user-stories/EP-RES.md#us-res-06-удалить-резюме) | Удалить резюме | Must | Реализовано | FR-RES-030, FR-RES-031 | — | `resume/application/commands/resume.go` | `components/ResumeList` | `commands/resume_test.go`, `persistence_integration_test.go` |
| [US-RES-07](02-user-stories/EP-RES.md#us-res-07-не-потерять-изменения-при-работе-в-двух-вкладках) | Не потерять изменения при работе в двух вкладках | Must | Частично | FR-RES-010, FR-RES-011, FR-RES-062 | — | `resume/infrastructure/persistence/repository/resume.go` | `utils/saveResumeSection` | `persistence_integration_test.go` |
| [US-RES-08](02-user-stories/EP-RES.md#us-res-08-видеть-заполненность-резюме) | Видеть заполненность резюме | Should | Реализовано | FR-RES-061 | — | — | `components/MainTab`, `utils/getResumeChecklist` | `getResumeChecklist.test.ts` |
| [US-RES-09](02-user-stories/EP-RES.md#us-res-09-скопировать-резюме) | Скопировать резюме | Could | Идея | FR-RES-070 | — | — | — | — |
| [US-RES-10](02-user-stories/EP-RES.md#us-res-10-история-изменений) | История изменений | Could | Идея | FR-RES-071 | — | — | — | — |

## EP-EDT. Редактор содержимого резюме

| История | Название | Приоритет | Статус | FR | NFR | Бэкенд | Фронтенд | Тесты |
|---|---|---|---|---|---|---|---|---|
| [US-EDT-01](02-user-stories/EP-EDT.md#us-edt-01-заполнить-обо-мне) | Заполнить «Обо мне» | Must | Реализовано | FR-EDT-001 – FR-EDT-006, FR-EDT-010 – FR-EDT-013 | — | `resume/transport/http/profile`, `resume/domain/profile.go` | `ResumeEditor/editors/AboutEditor` | `profile/handler_test.go`, `profile_replacement_test.go` |
| [US-EDT-02](02-user-stories/EP-EDT.md#us-edt-02-подставить-данные-из-профиля) | Подставить данные из профиля | Should | Реализовано | FR-EDT-114 | — | — | `AboutEditor/hooks/useFillFromProfile.ts` | `getLocalizedProfile.test.ts` |
| [US-EDT-03](02-user-stories/EP-EDT.md#us-edt-03-контакты) | Контакты | Must | Реализовано | FR-EDT-020 – FR-EDT-025, FR-CNT-001 | — | `resume/transport/http/contacts`, `resume/domain/contacts.go` | `ResumeEditor/editors/ContactsEditor` | `contacts_test.go`, `contacts_replacement_test.go` |
| [US-EDT-04](02-user-stories/EP-EDT.md#us-edt-04-навыки) | Навыки | Should | Реализовано | FR-EDT-030, FR-EDT-031 | — | `resume/transport/http/skills`, `resume/domain/skills.go` | `ResumeEditor/editors/SkillsEditor` | `skills_test.go`, `skills_replacement_test.go` |
| [US-EDT-05](02-user-stories/EP-EDT.md#us-edt-05-языки) | Языки | Should | Реализовано | FR-EDT-040 – FR-EDT-043 | — | `resume/transport/http/languages`, `resume/domain/languages.go` | `ResumeEditor/editors/LanguagesEditor` | `languages_test.go`, `languages_replacement_test.go` |
| [US-EDT-06](02-user-stories/EP-EDT.md#us-edt-06-стек-технологий) | Стек технологий | Must | Реализовано | FR-EDT-050 – FR-EDT-052 | — | `resume/transport/http/technologies`, `resume/domain/technologies.go` | `ResumeEditor/editors/TechStackEditor` | `technologies_test.go` |
| [US-EDT-07](02-user-stories/EP-EDT.md#us-edt-07-опыт-работы) | Опыт работы | Must | Реализовано | FR-EDT-060 – FR-EDT-064 | — | `resume/transport/http/experiences`, `resume/domain/experience.go` | `ResumeEditor/editors/ExperienceEditor`, `ExperienceSection/Dates.tsx` | `experience_test.go`, `experiences_replacement_test.go`, `formatDuration.test.ts` |
| [US-EDT-08](02-user-stories/EP-EDT.md#us-edt-08-образование) | Образование | Should | Реализовано | FR-EDT-070, FR-EDT-071 | — | `resume/transport/http/educations` | `ResumeEditor/editors/EducationEditor`, `EducationSection` | `education_test.go`, `educations_replacement_test.go` |
| [US-EDT-09](02-user-stories/EP-EDT.md#us-edt-09-курсы) | Курсы | Could | Реализовано | FR-EDT-080, FR-EDT-081 | — | `resume/transport/http/courses` | `ResumeEditor/editors/CoursesEditor`, `CoursesSection` | `course_test.go`, `courses_replacement_test.go` |
| [US-EDT-10](02-user-stories/EP-EDT.md#us-edt-10-дополнительная-активность) | Дополнительная активность | Could | Реализовано | FR-EDT-090, FR-EDT-091 | — | `resume/transport/http/activities` | `ResumeEditor/editors/ActivitiesEditor`, `ActivitiesSection` | `activity_test.go`, `activities_replacement_test.go` |
| [US-EDT-11](02-user-stories/EP-EDT.md#us-edt-11-проекты) | Проекты | Should | Частично | FR-EDT-100 – FR-EDT-104, FR-EDT-116 | — | `resume/transport/http/projects` | `ResumeEditor/editors/ProjectsEditor`, `ProjectsSection` | `project_test.go`, `project_replacement_test.go`, `buildSectionRequest.test.ts` |
| [US-EDT-12](02-user-stories/EP-EDT.md#us-edt-12-живой-предпросмотр-секции) | Живой предпросмотр секции | Should | Реализовано | FR-EDT-111 | — | — | `components/SectionPreview` | — |
| [US-EDT-13](02-user-stories/EP-EDT.md#us-edt-13-предупреждение-о-несохранённых-изменениях) | Предупреждение о несохранённых изменениях | Must | Реализовано | FR-EDT-112 | — | — | `screens/LkResumeEditScreen` | — |
| [US-EDT-14](02-user-stories/EP-EDT.md#us-edt-14-сбросить-изменения-вкладки) | Сбросить изменения вкладки | Should | Реализовано | FR-EDT-113 | — | — | `screens/LkResumeEditScreen` | — |
| [US-EDT-15](02-user-stories/EP-EDT.md#us-edt-15-посмотреть-резюме-глазами-посетителя) | Посмотреть резюме глазами посетителя | Must | Реализовано | FR-EDT-115, FR-ANL-003 | — | `resume/application/queries/public_resume.go` | `screens/LkResumeScreen` | `public_resume_test.go` |
| [US-EDT-16](02-user-stories/EP-EDT.md#us-edt-16-блочная-модель-резюме) | Блочная модель резюме | Could | Идея | FR-EDT-120 | — | — | — | — |

## EP-DSG. Оформление и структура страницы

| История | Название | Приоритет | Статус | FR | NFR | Бэкенд | Фронтенд | Тесты |
|---|---|---|---|---|---|---|---|---|
| [US-DSG-01](02-user-stories/EP-DSG.md#us-dsg-01-цветовая-гамма) | Цветовая гамма | Could | Реализовано | FR-DSG-001 – FR-DSG-003 | — | `resume/domain/valueobjects/theme.go` | `ResumeEditor/editors/ThemeEditor`, `ThemePreviewDialog`, `utils/createAccentTheme` | `overview_test.go`, `createAccentTheme.test.ts`, `isHexColor.test.ts` |
| [US-DSG-02](02-user-stories/EP-DSG.md#us-dsg-02-приветствие) | Приветствие | Should | Реализовано | FR-DSG-010 | — | `resume/transport/http/overview` | `ResumeEditor/editors/HeroEditor`, `HeroSection/Greeting.tsx` | `overview/request_test.go` |
| [US-DSG-03](02-user-stories/EP-DSG.md#us-dsg-03-готовность-к-предложениям) | Готовность к предложениям | Should | Реализовано | FR-DSG-011 | — | `resume/domain/valueobjects/profile.go` (`Availability`) | `HeroSection/Availability.tsx` | `overview_test.go` |
| [US-DSG-04](02-user-stories/EP-DSG.md#us-dsg-04-файл-резюме-для-скачивания) | Файл резюме для скачивания | Should | Реализовано | FR-DSG-012, FR-PUB-040 | — | `resume/transport/http/overview` | `ControlsBar/ResumeButton.tsx` | `url_test.go` |
| [US-DSG-05](02-user-stories/EP-DSG.md#us-dsg-05-порядок-видимость-и-свёрнутость-секций) | Порядок, видимость и свёрнутость секций | Should | Реализовано | FR-DSG-020 – FR-DSG-024 | — | `resume/domain/overview.go` | `ResumeEditor/editors/SectionsEditor`, `utils/withPinnedAbout`, `sections/` | `overview_test.go`, `withPinnedAbout.test.ts`, `defaultSections.test.ts` |
| [US-DSG-06](02-user-stories/EP-DSG.md#us-dsg-06-персональное-приветствие-для-ссылки) | Персональное приветствие для ссылки | Could | Идея | FR-DSG-030 | — | — | — | — |

## EP-LNK. Именованные ссылки

| История | Название | Приоритет | Статус | FR | NFR | Бэкенд | Фронтенд | Тесты |
|---|---|---|---|---|---|---|---|---|
| [US-LNK-01](02-user-stories/EP-LNK.md#us-lnk-01-создать-ссылку-для-адресата) | Создать ссылку для адресата | Must | Реализовано | FR-LNK-001 – FR-LNK-006 | — | `resume_link/application/commands/create.go`, `resume_link/domain/link.go` | `LinksTab`, `LinkNameDialog` | `create_test.go`, `link_test.go`, `resume_link/transport/http/handler_test.go` |
| [US-LNK-02](02-user-stories/EP-LNK.md#us-lnk-02-первая-ссылка) | Первая ссылка | Must | Реализовано | FR-LNK-020 | — | — | `CreateLinkPrompt`, `MainTab` | — |
| [US-LNK-03](02-user-stories/EP-LNK.md#us-lnk-03-скопировать-и-открыть-ссылку) | Скопировать и открыть ссылку | Must | Реализовано | FR-LNK-021 | — | — | `LinksTab/LinkRow.tsx`, `utils/getResumeLinkUrl` | `getResumeLinkUrl.test.ts` |
| [US-LNK-04](02-user-stories/EP-LNK.md#us-lnk-04-переименовать-ссылку) | Переименовать ссылку | Should | Частично | FR-LNK-010 | — | — (нет ручки) | `LinksTab/hooks/useLinksTab.ts` | — |
| [US-LNK-05](02-user-stories/EP-LNK.md#us-lnk-05-выключить-и-включить-ссылку) | Выключить и включить ссылку | Must | Реализовано | FR-LNK-011, FR-LNK-012 | — | `resume_link/application/commands/lifecycle.go` | `LinksTab/LinkActionsMenu.tsx`, `hooks/useResumeLinks.ts` | `lifecycle_test.go` (commands и transport) |
| [US-LNK-06](02-user-stories/EP-LNK.md#us-lnk-06-удалить-ссылку) | Удалить ссылку | Should | Реализовано | FR-LNK-013 | — | `resume_link/application/commands/delete.go` | `LinksTab/LinkActionsMenu.tsx` | `delete_test.go` |
| [US-LNK-07](02-user-stories/EP-LNK.md#us-lnk-07-счётчики-у-ссылки) | Счётчики у ссылки | Should | Частично | FR-LNK-030, FR-ANL-010 | — | `resume_link/application/queries/statistics.go` | — | `statistics_test.go` (queries и transport) |

## EP-PUB. Публичная страница резюме

| История | Название | Приоритет | Статус | FR | NFR | Бэкенд | Фронтенд | Тесты |
|---|---|---|---|---|---|---|---|---|
| [US-PUB-01](02-user-stories/EP-PUB.md#us-pub-01-открыть-резюме-по-ссылке) | Открыть резюме по ссылке | Must | Реализовано | FR-PUB-001 – FR-PUB-006 | — | `resume_link/application/queries/public.go` | `screens/CurriculumVitaeScreen`, `HeroSection` | `resume_link/.../public_test.go`, `public/handler_test.go` |
| [US-PUB-02](02-user-stories/EP-PUB.md#us-pub-02-навигация-по-секциям) | Навигация по секциям | Should | Реализовано | FR-PUB-010 – FR-PUB-012 | — | — | `SectionNav`, `ScrollProgress`, `utils/scrollToSection` | `scrollToSection.test.ts`, `getSectionElementId.test.ts` |
| [US-PUB-03](02-user-stories/EP-PUB.md#us-pub-03-где-кандидат-использовал-технологию) | Где кандидат использовал технологию | Should | Реализовано | FR-PUB-020 – FR-PUB-024 | — | — | `TechnologyFilterProvider`, `FilterBar`, `utils/hasTechnology` | `hasTechnology.test.ts`, `isSameTechnology.test.ts` |
| [US-PUB-04](02-user-stories/EP-PUB.md#us-pub-04-длинный-список-опыта) | Длинный список опыта | Could | Реализовано | FR-PUB-013 | — | — | `ExperienceSection` | `afterExpanded.test.ts` |
| [US-PUB-05](02-user-stories/EP-PUB.md#us-pub-05-свернуть-и-развернуть-секцию) | Свернуть и развернуть секцию | Could | Реализовано | FR-PUB-014, FR-DSG-023 | — | — | `Section`, `SectionCollapseProvider` | — |
| [US-PUB-06](02-user-stories/EP-PUB.md#us-pub-06-светлая-и-тёмная-тема) | Светлая и тёмная тема | Should | Реализовано | FR-PUB-030, FR-PUB-031 | — | — | `app/providers/ColorModeProvider`, `shared/components/ThemeToggle` | — |
| [US-PUB-07](02-user-stories/EP-PUB.md#us-pub-07-язык-интерфейса) | Язык интерфейса | Should | Реализовано | FR-PUB-032, FR-PUB-033 | — | — | `shared/components/LanguageSwitch`, `shared/utils/getInitialLanguage` | `getInitialLanguage.test.ts`, `mergeLocaleDicts.test.ts` |
| [US-PUB-08](02-user-stories/EP-PUB.md#us-pub-08-распечатать-или-сохранить-в-pdf) | Распечатать или сохранить в PDF | Should | Реализовано | FR-PUB-041 | — | — | `shared/hooks/usePrintLightTheme.ts` | — |
| [US-PUB-09](02-user-stories/EP-PUB.md#us-pub-09-скачать-файл-резюме) | Скачать файл резюме | Should | Реализовано | FR-PUB-040 | — | — | `ControlsBar/ResumeButton.tsx` | — |
| [US-PUB-10](02-user-stories/EP-PUB.md#us-pub-10-понятная-ошибка) | Понятная ошибка | Must | Реализовано | FR-PUB-003, FR-PUB-050 | — | `resume_link/transport/http/public/handler.go` | `components/LoadFailure`, `utils/parseApiError` | `parseApiError.test.ts`, `public/handler_test.go` |
| [US-PUB-11](02-user-stories/EP-PUB.md#us-pub-11-удобно-с-телефона) | Удобно с телефона | Must | Реализовано | FR-PUB-060 | NFR-UX-01, NFR-UX-02 | — | `shared/animations/motion.ts`, `CursorGlow`, `TiltCard` | — |
| [US-PUB-12](02-user-stories/EP-PUB.md#us-pub-12-карточка-ссылки-в-мессенджере) | Карточка ссылки в мессенджере | Should | Идея | FR-PUB-070 | — | — | — | — |
| [US-PUB-13](02-user-stories/EP-PUB.md#us-pub-13-посетитель-с-ограниченными-возможностями) | Посетитель с ограниченными возможностями | Should | Идея | — | NFR-UX-04 | — | — | — |

## EP-CNT. Контакты и защита от сбора

| История | Название | Приоритет | Статус | FR | NFR | Бэкенд | Фронтенд | Тесты |
|---|---|---|---|---|---|---|---|---|
| [US-CNT-01](02-user-stories/EP-CNT.md#us-cnt-01-связаться-в-один-клик) | Связаться в один клик | Must | Реализовано | FR-CNT-020, FR-CNT-021 | — | — | `ContactsSection`, `utils/getContactLinks`, `shared/hooks/useCopyToClipboard.ts` | `getContactLinks.test.ts` |
| [US-CNT-02](02-user-stories/EP-CNT.md#us-cnt-02-скрыть-контакты-от-парсеров) | Скрыть контакты от парсеров | Must | Реализовано | FR-CNT-001 – FR-CNT-004 | — | `resume/application/queries/public_resume.go` (`redactHiddenContacts`) | `ContactsSection` | `public_resume_test.go` |
| [US-CNT-03](02-user-stories/EP-CNT.md#us-cnt-03-раскрыть-скрытый-контакт) | Раскрыть скрытый контакт | Must | Реализовано | FR-CNT-010 – FR-CNT-014, FR-ANL-001 | — | `resume_link/application/queries/contact_view.go`, `resume/.../readers/resume.go` | `ContactsSection/hooks` | `contact_view_test.go` (оба модуля), `public/contact_view_test.go` |
| [US-CNT-04](02-user-stories/EP-CNT.md#us-cnt-04-проверка-на-бота-при-раскрытии) | Проверка на бота при раскрытии | Should | План | FR-CNT-015, FR-CNT-016 | NFR-SEC-06 | — | — | — |
| [US-CNT-05](02-user-stories/EP-CNT.md#us-cnt-05-свой-текст-блока-контактов) | Свой текст блока контактов | Could | Реализовано | FR-EDT-024 | — | `resume/domain/valueobjects/profile.go` (`ContactsText`) | `ContactsSection/Headline.tsx` | `contacts_test.go` |

## EP-ANL. Аналитика просмотров

| История | Название | Приоритет | Статус | FR | NFR | Бэкенд | Фронтенд | Тесты |
|---|---|---|---|---|---|---|---|---|
| [US-ANL-01](02-user-stories/EP-ANL.md#us-anl-01-сколько-раз-открыли-и-прочитали) | Сколько раз открыли и прочитали | Must | Частично | FR-ANL-001, FR-ANL-010, FR-ANL-011 | — | `resume_link/infrastructure/analytics/recorder.go`, `queries/events.go` | `components/AnalyticsTab` (заглушка) | `recorder_test.go`, `events_test.go`, `event_test.go` |
| [US-ANL-02](02-user-stories/EP-ANL.md#us-anl-02-уникальные-и-вернувшиеся-посетители) | Уникальные и вернувшиеся посетители | Must | План | FR-ANL-002, FR-ANL-011 | — | — | — | — |
| [US-ANL-03](02-user-stories/EP-ANL.md#us-anl-03-лента-визитов) | Лента визитов | Must | План | FR-ANL-012 | — | — | — | — |
| [US-ANL-04](02-user-stories/EP-ANL.md#us-anl-04-график-по-дням) | График по дням | Should | План | FR-ANL-013 | — | — | — | — |
| [US-ANL-05](02-user-stories/EP-ANL.md#us-anl-05-не-считать-себя-и-ботов) | Не считать себя и ботов | Must | План | FR-ANL-003, FR-ANL-004 | — | — | — | — |
| [US-ANL-06](02-user-stories/EP-ANL.md#us-anl-06-дочитывание-и-точка-ухода) | Дочитывание и точка ухода | Should | План | FR-ANL-020, FR-ANL-021 | — | — | — | — |
| [US-ANL-07](02-user-stories/EP-ANL.md#us-anl-07-какие-технологии-интересовали) | Какие технологии интересовали | Could | План | FR-ANL-022 | — | — | — | — |
| [US-ANL-08](02-user-stories/EP-ANL.md#us-anl-08-скачивания-и-раскрытия-контактов) | Скачивания и раскрытия контактов | Should | Частично | FR-ANL-001, FR-ANL-023 | — | `resume_link/application/queries/contact_view.go` | — | `contact_view_test.go` |
| [US-ANL-09](02-user-stories/EP-ANL.md#us-anl-09-сводка-по-резюме) | Сводка по резюме | Should | План | FR-ANL-014 | — | — | — | — |
| [US-ANL-10](02-user-stories/EP-ANL.md#us-anl-10-карта-кликов) | Карта кликов | Could | План | FR-ANL-030 | — | — | — | — |
| [US-ANL-11](02-user-stories/EP-ANL.md#us-anl-11-из-какой-компании-и-переслали-ли) | Из какой компании и переслали ли | Could | Идея | FR-ANL-040, FR-ANL-041 | — | — | — | — |
| [US-ANL-12](02-user-stories/EP-ANL.md#us-anl-12-посетитель-знает-о-статистике) | Посетитель знает о статистике | Must | План | FR-ANL-050 | NFR-PRV-02 | — | — | — |

## EP-NTF. Уведомления

| История | Название | Приоритет | Статус | FR | NFR | Бэкенд | Фронтенд | Тесты |
|---|---|---|---|---|---|---|---|---|
| [US-NTF-01](02-user-stories/EP-NTF.md#us-ntf-01-подключить-telegram) | Подключить Telegram | Must | План | FR-NTF-001 – FR-NTF-004 | — | — | — | — |
| [US-NTF-02](02-user-stories/EP-NTF.md#us-ntf-02-уведомление-о-прочтении) | Уведомление о прочтении | Must | План | FR-NTF-010, FR-NTF-011 | — | — | — | — |
| [US-NTF-03](02-user-stories/EP-NTF.md#us-ntf-03-уведомление-о-раскрытии-контакта) | Уведомление о раскрытии контакта | Should | План | FR-NTF-012 | — | — | — | — |
| [US-NTF-04](02-user-stories/EP-NTF.md#us-ntf-04-выбрать-о-чём-уведомлять) | Выбрать, о чём уведомлять | Should | План | FR-NTF-005 | — | — | — | — |
| [US-NTF-05](02-user-stories/EP-NTF.md#us-ntf-05-без-спама) | Без спама | Must | План | FR-NTF-020 | — | — | — | — |
| [US-NTF-06](02-user-stories/EP-NTF.md#us-ntf-06-отклик-интересно) | Отклик «Интересно» | Could | Идея | FR-NTF-030 | — | — | — | — |

## EP-OPS. Эксплуатация

| История | Название | Приоритет | Статус | FR | NFR | Бэкенд | Фронтенд | Тесты |
|---|---|---|---|---|---|---|---|---|
| [US-OPS-01](02-user-stories/EP-OPS.md#us-ops-01-автоматический-деплой) | Автоматический деплой | Must | Реализовано | FR-OPS-001 – FR-OPS-003 | NFR-MNT-03 | `.github/workflows/ci.yml` | `.github/workflows/ci.yml` | CI |
| [US-OPS-02](02-user-stories/EP-OPS.md#us-ops-02-миграции-при-деплое) | Миграции при деплое | Must | Реализовано | FR-OPS-004 | NFR-MNT-05 | `docker-compose.prod.yml`, `backend/db/migrations` | — | `compose config -q` в CI |
| [US-OPS-03](02-user-stories/EP-OPS.md#us-ops-03-проверка-живости-и-готовности) | Проверка живости и готовности | Must | Реализовано | FR-OPS-010, FR-OPS-011 | — | `platform/http/router.go`, `platform/health`, `platform/postgres` | — | `health_test.go`, `router_test.go` |
| [US-OPS-04](02-user-stories/EP-OPS.md#us-ops-04-метрики-и-трассировка) | Метрики и трассировка | Should | Реализовано | FR-OPS-012, FR-OPS-013 | NFR-OBS-02, NFR-OBS-03 | `platform/telemetry` | — | `telemetry/metrics_test.go`, `telemetry/http_test.go` |
| [US-OPS-05](02-user-stories/EP-OPS.md#us-ops-05-найти-запрос-в-логах) | Найти запрос в логах | Must | Реализовано | FR-OPS-014 | NFR-OBS-01 | `platform/http/middleware/logger.go`, `middleware/request_id.go` | `utils/parseApiError` (показ ID запроса) | `logger_test.go`, `request_id_test.go` |
| [US-OPS-06](02-user-stories/EP-OPS.md#us-ops-06-резервные-копии-и-восстановление) | Резервные копии и восстановление | Must | План | FR-OPS-020 | NFR-REL-04 | — | — | — |
| [US-OPS-07](02-user-stories/EP-OPS.md#us-ops-07-алерты) | Алерты | Should | План | FR-OPS-021 | NFR-OBS-04 | — | — | — |
| [US-OPS-08](02-user-stories/EP-OPS.md#us-ops-08-защита-от-перегрузки) | Защита от перегрузки | Must | Частично | FR-API-020, FR-API-021 | NFR-SEC-06 | `platform/http/middleware/rate_limit.go`, `platform/http/router.go` (`BodyLimit`) | — | `rate_limit_test.go` |
