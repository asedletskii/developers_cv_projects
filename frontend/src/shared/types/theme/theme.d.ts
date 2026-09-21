import { PaletteOptions as MuiPaletteOptions } from '@mui/material/styles'

declare module '@mui/material/styles' {
  interface Palette {
    gray: Palette['primary']
    surface: Palette['primary']
    link: Palette['primary']
    info: Palette['primary']
    text: Palette['primary']
    alert: Palette['primary']
  }

  interface PaletteOptions {
    gray: MuiPaletteOptions['primary']
    surface: MuiPaletteOptions['primary']
    link: MuiPaletteOptions['primary']
    info: MuiPaletteOptions['primary']
    text: MuiPaletteOptions['primary']
    alert: MuiPaletteOptions['primary']
  }

  type SvgIconPropsColorOverrides = 'gray'
}
