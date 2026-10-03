import {
  DashboardCustomizeRounded,
  LinkRounded,
  PaletteRounded,
  PrintRounded,
  Telegram,
  VisibilityRounded,
} from '@mui/icons-material'
import { keyframes } from '@mui/material/styles'

/** фон закрывающего баннера плавно переливается — та же идея, что у блока контактов резюме. */
export const DRIFT = keyframes`
  0%, 100% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
`

/** заполнение полоски прогресса активного шага — как в историях. */
export const STEP_FILL = keyframes`
  from { width: 0%; }
  to { width: 100%; }
`

export const STEP_DURATION_MS = 5000

/** извилистая «дорога» для десктопной версии шагов: SVG-путь + координаты узлов на нём —
 *  числа подобраны руками под viewBox ниже, но держат равный шаг между узлами. */
export const ROAD_VIEWBOX = '0 0 300 900'
export const ROAD_PATH_D =
  'M 150 30 C 280 110, 280 230, 150 310 C 20 390, 20 510, 150 590 C 280 670, 280 790, 150 870'
export const ROAD_NODES = [
  { x: 150, y: 30 },
  { x: 150, y: 310 },
  { x: 150, y: 590 },
  { x: 150, y: 870 },
]

export const STEPS = [
  { key: 'login', Icon: Telegram },
  { key: 'build', Icon: DashboardCustomizeRounded },
  { key: 'preview', Icon: VisibilityRounded },
  { key: 'share', Icon: LinkRounded },
] as const

export const FEATURES = [
  { key: 'auth', Icon: Telegram },
  { key: 'sections', Icon: DashboardCustomizeRounded },
  { key: 'livePreview', Icon: VisibilityRounded },
  { key: 'links', Icon: LinkRounded },
  { key: 'theme', Icon: PaletteRounded },
  { key: 'print', Icon: PrintRounded },
] as const

export const TECH_STACK = [
  'React',
  'TypeScript',
  'MUI',
  'Redux Toolkit',
  'React Hook Form',
  'GSAP',
  'Go',
  'PostgreSQL',
  'Docker',
  'Caddy',
]
