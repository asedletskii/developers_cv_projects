import { PropsWithChildren } from 'react'
import { ThemeSettings } from '@features/curriculumVitae/types'
import { ThemeProvider } from '@mui/material/styles'
import { useAccentTheme } from './hooks'

type Props = PropsWithChildren<{
  colors?: ThemeSettings
}>

export const AccentTheme = ({ colors, children }: Props) => {
  const theme = useAccentTheme(colors)

  if (!theme) return <>{children}</>

  return <ThemeProvider theme={theme}>{children}</ThemeProvider>
}
