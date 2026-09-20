# developers-cv.ru

React + TypeScript + Vite, MUI, Redux Toolkit, i18next.

## Запуск

```bash
yarn install
yarn start
```

## Скрипты

| Команда           | Что делает                      |
| ----------------- | ------------------------------- |
| `yarn start`      | dev-сервер                      |
| `yarn build`      | продакшн-сборка в `build/`      |
| `yarn lint`       | ESLint (без автоисправления)    |
| `yarn lint:fix`   | ESLint с автоисправлением       |
| `yarn typecheck`  | проверка типов (`tsc --noEmit`) |
| `yarn test`       | тесты Jest                      |
| `yarn test:watch` | тесты Jest в режиме наблюдения  |
| `yarn format`     | Prettier по всему проекту       |

CI (`.github/workflows/ci.yml`) запускает lint, typecheck, тесты и build на каждый push в `main` и PR.

## Структура и правила

Каждый компонент лежит в своей папке и отдаётся через `index.ts`. Подкомпоненты лежат **плоскими файлами
рядом с родителем**, без вложенных папок `components/`:

```
AboutSection/
  AboutSection.tsx   ← родитель
  Portrait.tsx       ← подкомпоненты
  Summary.tsx
  Facts.tsx
  index.ts           ← экспортирует только родителя
```

Подкомпонент подкомпонента (например, `ContactButton` внутри `Actions`) тоже кладётся на этот же уровень.
Если подкомпонент нужен нескольким родителям, он поднимается в свою папку на уровень выше.

- **Сложные компоненты раскладываются на подкомпоненты.** Логика, которая разрастается, уезжает в хук рядом (`useExpandableJobs.ts`).
- **Имя ребёнка не повторяет имя родителя.** Родитель `ExperienceSection` → ребёнок `Job`, а не `ExperienceJob`; родитель `SectionNav` → `Item`, а не `SectionNavItem`.
- **Имена не должны совпадать с компонентами библиотек** (MUI, react-router и т.д.). Поэтому `Card` → `Project`, `Skeleton` → `LoadPlaceholder`, `Error` → `LoadFailure`, `RouterProvider` → `RoutingProvider`, `ThemeProvider` → `ColorModeProvider`.
- **Компонент, нужный нескольким родителям, поднимается на уровень выше** (`TechnologyList`, `IconTile`, `SparkleField` лежат в `components/` фичи и импортируются по алиасу).
- **Конечные наборы значений — только `enum`**, а не строковые union'ы и не литералы в коде (`WorkFormat.Remote`, `SectionId.About`). Значение enum совпадает с кодом из API. Доменные enum'ы лежат в `features/curriculumVitae/types/enums.ts`, общие (`Languages`, `ColorMode`, `ToastType`) в `shared/types`.
- **Утилиты не живут в папках компонентов.** Чистые функции лежат в `utils`: общие в `src/shared/utils`, относящиеся к резюме в `src/features/curriculumVitae/utils`. Каждая утилита в своей папке вместе с тестом Jest на том же уровне (`formatAge/formatAge.ts`, `formatAge/formatAge.test.ts`, `formatAge/index.ts`); наружу всё отдаёт `utils/index.ts`. Внутри слоёв с риском циклов (стор, словари) утилиты импортируются по прямому пути, например `@shared/utils/createLazySliceStateSelector`.
- Не компоненты (реестр секций, конфиги маршрутов, помощники) лежат отдельными модулями: `sections/`, `routes.tsx`, `getContactLinks.tsx`.

## Документация

- [API резюме для бэкенда](docs/backend-api.md): что фронт ждёт от сервера (эндпоинт, формат данных, правила).
