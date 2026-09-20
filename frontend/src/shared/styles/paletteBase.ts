import { PaletteOptions } from '@mui/material/styles'

/** Цвета, которые одинаковы в обеих темах. Оттенки: light — светлее, dark — темнее. */
export const brand = {
  primary: {
    light: '#FFDD7A',
    main: '#FFCC3C',
    dark: '#C99A1E',
    contrastText: '#1F1F1D',
  },

  secondary: {
    light: '#F2A65A',
    main: '#D97A3D',
    dark: '#A85A26',
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
