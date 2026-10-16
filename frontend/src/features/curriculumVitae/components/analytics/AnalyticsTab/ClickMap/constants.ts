import { ClickBreakpoint } from '@features/curriculumVitae/types'

export const BREAKPOINT_WIDTH: Record<ClickBreakpoint, number> = { xs: 375, md: 900, lg: 1280 }

export const BREAKPOINTS = Object.keys(BREAKPOINT_WIDTH) as ClickBreakpoint[]

export const VIEWPORT_HEIGHT = 760
