# Стайлгайд фронтенда

Правила кода в `frontend/`. Часть из них проверяет линтер (`yarn lint`), остальные держатся на ревью. Примеры — из фичи `features/curriculumVitae`.

Перед коммитом должно проходить:

```bash
yarn lint && yarn typecheck && yarn test && yarn build
```

`yarn lint:fix` сам исправляет форматирование и порядок импортов.

## 1. Форматирование

Настроено в `.prettierrc` и `eslint.config.mjs`, руками ничего выравнивать не нужно.

| Правило                               | Значение                                                                                   |
| ------------------------------------- | ------------------------------------------------------------------------------------------ |
| Отступ                                | 2 пробела                                                                                  |
| Ширина строки                         | 100 символов                                                                               |
| Точка с запятой                       | не ставим                                                                                  |
| Кавычки                               | одинарные                                                                                  |
| Запятая в конце многострочных списков | всегда (`trailingComma: all`)                                                              |
| Скобки у одного параметра стрелки     | всегда: `(item) => …`                                                                      |
| Порядок импортов                      | `simple-import-sort`: `react` первым, затем побочные, пакеты и алиасы, затем относительные |
| Дубли импортов                        | запрещены (`import/no-duplicates`)                                                         |
| Неиспользуемые переменные             | предупреждение; `catch (e)` без использования допустим                                     |

Комментарии и подписи в интерфейсе — на русском. Документирующий комментарий (`/** … */`) пишем над компонентом, хуком и утилитой: **что делает и почему так**, а не пересказ кода. Локальные пометки — обычные `//`.

## 2. Функции: только стрелочные выражения

Всё пишем через `() => {}`. Это проверяется правилами `func-style: expression` и `prefer-arrow-callback`, объявления `function` не проходят линтер.

```ts
// хорошо
export const formatAge = (dateOfBirth: string) => {
  return dayjs().diff(dayjs(dateOfBirth), 'year')
}

// плохо
export function formatAge(dateOfBirth: string) { … }
```

- **Компоненты, хуки, утилиты, обработчики, колбэки** — стрелки.
- **Экран (страница роутера)** — тоже стрелка, а экспортируется по умолчанию отдельной строкой (нужно для `lazy`):

  ```tsx
  const LkResumesScreen = () => {
    …
  }

  export default LkResumesScreen
  ```

- **Порядок объявлений важен.** Стрелка в `const` не «поднимается»: то, что используется на верхнем уровне модуля, объявляем раньше использования. Внутри функций порядок не критичен — тело выполняется позже.
- Фабрики, возвращающие функцию, возвращают стрелку, присвоенную константе (нужно, например, чтобы линтер хуков видел имя `use…`):

  ```ts
  const useTranslate = (prefix: string) => { … }
  return useTranslate
  ```

## 3. Куда что класть

| Что                                  | Где                                                                                        | Пример                                  |
| ------------------------------------ | ------------------------------------------------------------------------------------------ | --------------------------------------- |
| Экран (страница роутера)             | `features/<фича>/screens/<Экран>/<Экран>.tsx` — только сборка из компонентов и вызов хуков | `LkResumeEditScreen`                    |
| Компонент                            | своя папка `components/<Имя>/<Имя>.tsx` с `index.ts` (`export * from './<Имя>'`)           | `LinksTab`                              |
| Подкомпонент (нужен только родителю) | файл рядом с родителем в его папке                                                         | `LinksTab/LinkRow.tsx`                  |
| Вложенные компоненты одного рода     | подпапка родителя                                                                          | `ResumeEditor/editors/<Имя>/`           |
| Хук одного компонента                | `components/<Имя>/hooks/` (+ `index.ts`)                                                   | `LinksTab/hooks/useLinksTab.ts`         |
| Хук одного экрана                    | `screens/<Экран>/hooks/`                                                                   | `useResumeSections`                     |
| Хук, общий для фичи                  | `features/<фича>/hooks/`                                                                   | `useOwnerResume`, `useFormList`         |
| Хук, общий для проекта               | `shared/hooks/`                                                                            | `useCopyToClipboard`, `useScrollReveal` |
| Чистая функция без React             | `features/<фича>/utils/<имя>/<имя>.ts` (+ `index.ts`, `<имя>.test.ts`) или `shared/utils/` | `getResumeChecklist`, `syncList`        |
| Константы компонента                 | `constants.ts` в папке компонента                                                          | `HeroSection/constants.ts`              |
| Константы фичи                       | `features/<фича>/constants/`                                                               | `TECHNOLOGY_SUGGESTIONS`, `API_ROUTES`  |
| Константы проекта                    | `shared/constants/`                                                                        | `ROUTES`                                |
| Типы                                 | `features/<фича>/types/`                                                                   | `ResumeEditorTab`                       |
| Состояние и запросы                  | `store/`, `api/`                                                                           | `curriculumVitaeApi`                    |

### Правила распределения

1. **Логика — в хуки.** Компонент отвечает за разметку: получает данные из хука и вызывает его обработчики. Состояние, эффекты, обращения к API и хранилищу, производные значения и подтверждения (`window.confirm`) живут в хуках.
2. **Чистая логика — в `utils`.** Если функции не нужен React (расчёт чек-листа, сопоставление списков, сборка адреса), это утилита с тестом, а хук лишь вызывает её.
3. **Константы — в `constants.ts` рядом с компонентом.** Модульные константы (списки, карты значений, `keyframes`, повторяющиеся `sx`-объекты) не объявляем в файле компонента или хука, а выносим и импортируем. Имена констант — `UPPER_SNAKE_CASE` (`PULSE`, `ORBS`, `STATUS_COLORS`).
4. **Один `index.ts` на папку** компонента или хуков — публичный вход; всё, что в нём не экспортировано, считается внутренним.
5. **Хук общего назначения выше по дереву.** Как только хук нужен второму компоненту, он поднимается на уровень выше (компонент → фича → `shared`).
6. **Экран тонкий.** В файле экрана только вызов хуков и композиция компонентов; разметки и логики в нём нет.

### Дерево редактора резюме (пример)

```
components/
  ResumeEditor/            вкладки секций: конфиг в constants.ts, hooks/useResumeEditor
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

## 4. Импорты

- Только абсолютные: `@features/curriculumVitae/...`, `@shared/...`, `@app/...` (алиасы в `tsconfig.paths.json` и `jest.config.cjs`).
- Импортируем через `index.ts` папки: `@features/curriculumVitae/components/LinksTab`, а не `…/LinksTab/LinksTab`.
- Глубокий путь допустим внутри одной папки компонента и там, где импорт через барель дал бы циклическую зависимость (например, `@features/curriculumVitae/hooks/useSectionCollapse` из провайдера).
- Относительные импорты — только между соседними файлами одной папки (`./hooks`, `./constants`).

## 5. Роутинг: только константы и `generatePath`

Адреса не пишем строками и не склеиваем через шаблонные строки и `encodeURIComponent`.

- **Страницы** — шаблоны в `ROUTES` (`shared/constants/routes.constants.ts`). Ссылки, `navigate` и `<Navigate>` собираются через `generatePath(ROUTES.x, params)`; она же экранирует параметры.

  ```tsx
  <Link to={generatePath(ROUTES.lkResume, { userId, cvId: resume.id })} />
  navigate(generatePath(ROUTES.lkProfile, { userId }))
  <Navigate to={ROUTES.home} replace />   // без параметров генератор не нужен
  ```

- **Определения маршрутов** (`routes.tsx`) используют те же `ROUTES` абсолютными путями: `<Route path={ROUTES.lkResume} … />`.
- **Ручки бэкенда** — шаблоны в `API_ROUTES` (`features/<фича>/constants/apiRoutes.ts`), запросы собираем так же: `api.get(generatePath(API_ROUTES.lkResume, { resumeId }))`.
- Проверка «находимся ли мы в разделе» — через `matchPath({ path: ROUTES.cv, end: false }, pathname)`, а не `startsWith`.
- Новая страница = новая запись в `ROUTES` + маршрут в `routes.tsx`; строк с путями в компонентах быть не должно.

## 6. Компоненты и стили

- Один компонент — один файл; `Props` объявляем над компонентом (`type Props = { … }`), без `React.FC`.
- Стили — через `sx` и тему MUI; общие куски (`glassCard`) лежат в `styles/`. Повторяющиеся `sx`-объекты выносим в `constants.ts`.
- Тексты интерфейса — только через `useTranslate('Раздел.Подраздел')` и словари `locale/`; литералов на русском в разметке нет (исключение — демо-данные предпросмотра).
- Формы — `react-hook-form`; путь поля типизирован моделью формы. Списки записей в форме ведёт `useFormList`.

## 7. Тесты

- Jest, файл теста рядом с кодом: `<имя>.test.ts`.
- Чистые утилиты покрываем тестом обязательно: минимум один случай на ветку логики.
- Компоненты и хуки отдельно не тестируем, пока в них нет логики, которой нет в утилитах; логика уходит в утилиты именно поэтому.
