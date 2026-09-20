import { PaletteOptions } from '@mui/material/styles'
import { brand, status } from './paletteBase'

/**
 * Смысл оттенков тот же, что в светлой теме, поэтому gray инвертирован:
 * gray.dark — самый контрастный (заголовки), gray.light — приглушённый.
 */
export const paletteDark: PaletteOptions = {
  mode: 'dark',

  ...brand,
  ...status,

  surface: {
    light: '#26262A',
    main: '#151517',
    dark: '#0C0C0D',
  },

  gray: {
    light: '#8F8F96',
    main: '#C6C6CC',
    dark: '#F4F4F6',
  },

  link: {
    light: '#8DB8FB',
    main: '#5C9AF8',
    dark: '#3D82EE',
  },

  background: {
    default: '#151517',
    paper: '#26262A',
  },

  text: {
    primary: '#F4F4F6',
    secondary: '#A1A1A8',
  },

  divider: '#34343A',
}
