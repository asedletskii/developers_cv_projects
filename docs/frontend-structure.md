# Структура фронтенда

Правила, по которым лежит код в `frontend/src`. Пример — фича `features/curriculumVitae`.

## Куда что класть

| Что | Где | Пример |
|---|---|---|
| Экран (страница роутера) | `features/<фича>/screens/<Экран>/<Экран>.tsx` — только сборка из компонентов и вызов хуков | `LkResumeEditScreen` |
| Компонент | своя папка `components/<Имя>/<Имя>.tsx` с `index.ts` (`export * from './<Имя>'`) | `LinksTab` |
| Подкомпонент (нужен только родителю) | файл рядом с родителем в его папке | `LinksTab/LinkRow.tsx` |
| Вложенные компоненты одного рода | подпапка родителя | `ResumeEditor/editors/<Имя>/` |
| Хук одного компонента | `components/<Имя>/hooks/` (+ `index.ts`) | `LinksTab/hooks/useLinksTab.ts` |
| Хук одного экрана | `screens/<Экран>/hooks/` | `useResumeSections` |
| Хук, общий для фичи | `features/<фича>/hooks/` | `useOwnerResume`, `useFormList` |
| Хук, общий для проекта | `shared/hooks/` | `useCopyToClipboard`, `useScrollReveal` |
| Чистая функция без React | `features/<фича>/utils/<имя>/<имя>.ts` (+ `index.ts`, `.test.ts`) или `shared/utils/` | `getResumeChecklist`, `syncList` |
| Константы | `constants/` фичи или `constants.ts` рядом с компонентом | `TECHNOLOGY_SUGGESTIONS` |
| Типы | `features/<фича>/types/` | `ResumeEditorTab` |
| Состояние и запросы | `store/`, `api/` | `curriculumVitaeApi` |

## Правила

1. **Логика — в хуки.** Компонент отвечает за разметку: получает данные из хука и вызывает его обработчики. Состояние, эффекты, обращения к API и хранилищу, производные значения и подтверждения (`window.confirm`) живут в хуках.
2. **Чистая логика — в `utils`.** Если функции не нужен React (расчёт чек-листа, сопоставление списков, сборка адреса), это утилита с тестом, а хук лишь вызывает её.
3. **Импорты — абсолютные** (`@features/curriculumVitae/...`, `@shared/...`), через `index.ts` папки; глубокий путь допустим внутри одной папки компонента и там, где барель дал бы цикл.
4. **Один `index.ts` на папку** компонента или хуков — публичный вход; всё, что в нём не экспортировано, считается внутренним.
5. **Хук общего назначения выше по дереву.** Как только хук нужен второму компоненту, он поднимается на уровень выше (компонент → фича → `shared`).

## Дерево редактора резюме (как пример)

```
components/
  ResumeEditor/            вкладки секций: конфиг tabs.ts, hooks/useResumeEditor
    editors/<Имя>/         редактор одной секции (+ hooks/ с вариантами и логикой)
  SectionForm/             форма вкладки: hooks/useSectionForm, SectionFormActions
  SectionPreview/          предпросмотр секции: hooks/useSectionPreview
  HeroPreview/, ThemePreviewDialog/
  MainTab/, LinksTab/, AnalyticsTab/   разделы страницы резюме в кабинете
  ResumeHeader/            шапка: hooks/useResumeHeader
  form/                    поля форм: SelectField, DateField, ColorField, TechnologyListField, …
hooks/                     useOwnerResume, useLkProfile, useLkResumes, usePublicResume, useFormList, …
utils/                     viewToPerson, saveResumeSection, syncList, getResumeChecklist, blankItems, …
screens/LkResumeEditScreen/  только экран + hooks/useResumeSections + constants.ts
```
