import { PropsWithChildren, useMemo } from 'react'
import { ThemeSettings } from '@features/curriculumVitae/types'
import { ThemeProvider, useTheme } from '@mui/material/styles'
import { createAccentTheme } from '@features/curriculumVitae/utils'

type Props = PropsWithChildren<{
  colors?: ThemeSettings
}>

/** Красит резюме в цвета из настроек: вложенная тема поверх текущей, светлая/тёмная схема сохраняется. */
export const AccentTheme = ({ colors, children }: Props) => {
  const baseTheme = useTheme()
  const theme = useMemo(() => createAccentTheme(baseTheme, colors), [baseTheme, colors])

  if (!theme) return <>{children}</>

  return <ThemeProvider theme={theme}>{children}</ThemeProvider>
}
