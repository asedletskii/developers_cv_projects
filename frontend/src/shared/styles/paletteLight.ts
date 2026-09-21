import { PaletteOptions } from '@mui/material/styles'
import { brand, status } from './paletteBase'

/**
 * surface: light — приподнятые элементы (карточки), main — фон страницы, dark — подвал/акцентные блоки.
 * gray: light — приглушённый текст, main — обычный, dark — заголовки (самый контрастный).
 */
export const paletteLight: PaletteOptions = {
  mode: 'light',

  ...brand,
  ...status,

  surface: {
    light: '#FFFFFF',
    main: '#F7F6F4',
    dark: '#ECEAE6',
  },

  gray: {
    light: '#85857F',
    main: '#4A4A46',
    dark: '#1F1F1D',
  },

  link: {
    light: '#5C9AF8',
    main: '#2F79E8',
    dark: '#1B5CC4',
  },

  background: {
    default: '#F7F6F4',
    paper: '#FFFFFF',
  },

  text: {
    primary: '#1F1F1D',
    secondary: '#6B6B66',
  },

  divider: '#E2E0DB',
}
