import { useMemo } from 'react'
import { ThemeSettings } from '@features/curriculumVitae/types'
import { createAccentTheme } from '@features/curriculumVitae/utils'
import { useTheme } from '@mui/material/styles'

export const useAccentTheme = (colors?: ThemeSettings) => {
  const baseTheme = useTheme()

  return useMemo(() => createAccentTheme(baseTheme, colors), [baseTheme, colors])
}
