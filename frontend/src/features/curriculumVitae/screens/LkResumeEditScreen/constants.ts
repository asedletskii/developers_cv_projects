/** Разделы страницы резюме в кабинете, в порядке показа. */
export const TOP_SECTIONS = ['main', 'edit', 'links', 'analytics'] as const

export type TopSection = (typeof TOP_SECTIONS)[number]
