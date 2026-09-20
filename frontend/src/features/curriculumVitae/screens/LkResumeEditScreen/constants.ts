export const TOP_SECTIONS = ['main', 'edit', 'links', 'analytics'] as const

export type TopSection = (typeof TOP_SECTIONS)[number]

export const TAB_SX = { textTransform: 'none', minWidth: 0, px: { xs: 1.5, sm: 2 } } as const
