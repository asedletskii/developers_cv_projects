import { ClickBreakpoint } from '@features/curriculumVitae/types'

/** Ширина предпросмотра на каждый брейкпоинт: `xs` < 600, `md` < 1200, `lg` — остальное */
export const BREAKPOINT_WIDTH: Record<ClickBreakpoint, number> = { xs: 375, md: 900, lg: 1280 }

export const BREAKPOINTS = Object.keys(BREAKPOINT_WIDTH) as ClickBreakpoint[]

/** Высота окна предпросмотра в натуральных пикселях; страница прокручивается внутри него, как на устройстве */
export const VIEWPORT_HEIGHT = 760
