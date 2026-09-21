import { ThemeSettings } from '@features/curriculumVitae/types'
import { createTheme, Theme } from '@mui/material/styles'
import { isHexColor } from '@shared/utils'

const ACCENTS = ['primary', 'secondary'] as const

/**
 * Накладывает акцентные цвета с бэка на текущую тему. Светлый и тёмный оттенки и цвет текста
 * считает MUI (augmentColor), так что контраст на кнопках и плашках сохраняется при любом цвете.
 * Невалидные значения игнорируются; если применять нечего, возвращает null.
 */
export const createAccentTheme = (base: Theme, colors?: ThemeSettings): Theme | null => {
  const palette = Object.fromEntries(
    ACCENTS.flatMap((name) => {
      const main = colors?.[name]

      return isHexColor(main) ? [[name, base.palette.augmentColor({ color: { main }, name })]] : []
    }),
  )

  return Object.keys(palette).length ? createTheme(base, { palette }) : null
}
