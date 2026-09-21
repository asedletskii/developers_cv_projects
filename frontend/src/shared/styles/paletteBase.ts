import { PaletteOptions } from '@mui/material/styles'

/** Цвета, которые одинаковы в обеих темах. Оттенки: light — светлее, dark — темнее. */
export const brand = {
  primary: {
    light: '#E88E74',
    main: '#DE6B48',
    dark: '#C45532',
    contrastText: '#FFFFFF',
  },

  secondary: {
    light: '#6BA3FA',
    main: '#428BF9',
    dark: '#1F6FE0',
    contrastText: '#FFFFFF',
  },
} satisfies Pick<PaletteOptions, 'primary' | 'secondary'>

/** Статусные цвета: light — фон/подложка, main — акцент, dark — текст на подложке. */
export const status = {
  success: { light: '#D3F7DF', main: '#00A63E', dark: '#016630' },
  info: { light: '#D6E6FF', main: '#155DFC', dark: '#193CB8' },
  alert: { light: '#FFF3B8', main: '#FFAB00', dark: '#894B00' },
  error: { light: '#FFD9D9', main: '#E7000B', dark: '#9F0712' },
} satisfies Pick<PaletteOptions, 'success' | 'info' | 'alert' | 'error'>
