import {
  DashboardCustomizeRounded,
  DownloadRounded,
  InsightsRounded,
  LinkRounded,
  PaletteRounded,
  PrintRounded,
  ShieldRounded,
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
  { key: 'analytics', Icon: InsightsRounded },
  { key: 'protection', Icon: ShieldRounded },
  { key: 'download', Icon: DownloadRounded },
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

/** две бегущие строки стека едут в разные стороны: первая половина и вторая */
export const TECH_ROWS = [TECH_STACK.slice(0, 5), TECH_STACK.slice(5)]

export const STATS = [
  { key: 'sections', value: 9 },
  { key: 'languages', value: 2 },
  { key: 'contacts', value: 3 },
  { key: 'read', value: 5 },
] as const

export const SHOWCASE = [
  { key: 'builder', Icon: DashboardCustomizeRounded },
  { key: 'analytics', Icon: InsightsRounded },
  { key: 'protection', Icon: ShieldRounded },
  { key: 'download', Icon: DownloadRounded },
] as const

export const FAQ = ['data', 'dnt', 'hidden', 'languages', 'links'] as const

/** парящий элемент: плавно вверх-вниз, у каждого свой сдвиг по времени */
export const FLOAT_Y = keyframes`
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-12px); }
`

/** бегущая строка: сдвиг на половину ширины, потому что содержимое продублировано */
export const MARQUEE = keyframes`
  from { transform: translateX(0); }
  to { transform: translateX(-50%); }
`

/** строки макета резюме «дорисовываются» слева направо */
export const DRAW_LINE = keyframes`
  from { transform: scaleX(0); }
  to { transform: scaleX(1); }
`

/** столбик диаграммы вырастает от основания */
export const GROW_BAR = keyframes`
  from { transform: scaleY(0); }
  to { transform: scaleY(1); }
`

/** расходящееся кольцо вокруг щита и индикатора доступности */
export const RING = keyframes`
  0% { transform: scale(0.6); opacity: 0.55; }
  100% { transform: scale(1.9); opacity: 0; }
`

/** новый контент вкладки выезжает снизу */
export const SLIDE_UP = keyframes`
  from { opacity: 0; transform: translateY(18px); }
  to { opacity: 1; transform: none; }
`

/** карточки секций в «конструкторе» по очереди съезжаются в стопку */
export const STACK_IN = keyframes`
  0% { opacity: 0; transform: translateX(-24px); }
  20%, 85% { opacity: 1; transform: none; }
  100% { opacity: 0; transform: translateX(24px); }
`

/** листок падает в «загрузки» и всё начинается сначала */
export const DROP_FILE = keyframes`
  0% { transform: translateY(-28px); opacity: 0; }
  30% { opacity: 1; }
  70%, 100% { transform: translateY(14px); opacity: 1; }
`

/** тепловое пятно на карте кликов пульсирует */
export const HEAT_PULSE = keyframes`
  0%, 100% { transform: scale(0.85); opacity: 0.55; }
  50% { transform: scale(1.15); opacity: 1; }
`

/** вместе с `prefers-reduced-motion` движение отключается целиком */
export const STILL = { '@media (prefers-reduced-motion: reduce)': { animation: 'none' } } as const
