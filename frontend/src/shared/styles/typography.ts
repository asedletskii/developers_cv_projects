import { TypographyVariantsOptions } from '@mui/material/styles'

const font =
  '"SF Pro Text", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif'

/**
 * Трекинг и интерлиньяж зависят от размера: чем крупнее текст, тем плотнее буквы и строки;
 * мелкий текст, наоборот, чуть разрежен для читаемости.
 */
export const typography: TypographyVariantsOptions = {
  fontFamily: font,
  fontWeightRegular: 400,

  h1: { fontWeight: 900, fontSize: '3.125rem', lineHeight: 1.05, letterSpacing: '-0.03em' },
  h2: { fontWeight: 700, fontSize: '2.75rem', lineHeight: 1.1, letterSpacing: '-0.025em' },
  h3: { fontWeight: 700, fontSize: '2rem', lineHeight: 1.15, letterSpacing: '-0.02em' },
  h4: { fontWeight: 700, fontSize: '1.75rem', lineHeight: 1.2, letterSpacing: '-0.015em' },
  h5: { fontWeight: 700, fontSize: '1.5rem', lineHeight: 1.25, letterSpacing: '-0.01em' },

  h6: { fontWeight: 400, fontSize: '1.25rem', lineHeight: 1.35, letterSpacing: '-0.005em' },

  body1: { fontWeight: 400, fontSize: '1rem', lineHeight: 1.5, letterSpacing: 0 },

  caption: { fontWeight: 400, fontSize: '0.75rem', lineHeight: 1.4, letterSpacing: '0.01em' },
}
