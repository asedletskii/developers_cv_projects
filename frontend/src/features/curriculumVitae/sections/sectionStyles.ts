export const SECTION_STYLES = ['first', 'second', 'third'] as const

export type SectionStyle = (typeof SECTION_STYLES)[number]

export const DEFAULT_SECTION_STYLE: SectionStyle = 'first'

export const resolveSectionStyle = (requested?: string): SectionStyle =>
  SECTION_STYLES.find((style) => style === requested) ?? DEFAULT_SECTION_STYLE

export const byStyle = <T>(style: SectionStyle, variants: Record<SectionStyle, T>): T =>
  variants[style]
