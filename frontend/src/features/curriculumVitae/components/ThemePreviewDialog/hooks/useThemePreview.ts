import { useState } from 'react'
import { Resume } from '@features/curriculumVitae/types'
import { createAccentTheme } from '@features/curriculumVitae/utils'
import { darkTheme } from '@shared/styles/darkTheme'
import { lightTheme } from '@shared/styles/lightTheme'
import { ColorMode } from '@shared/types'
import { useFormContext, useWatch } from 'react-hook-form'

/** Тема предпросмотра: цвета из формы поверх светлой или тёмной базовой темы. */
export const useThemePreview = () => {
  const { control } = useFormContext<Resume>()
  const primary = useWatch({ control, name: 'overview.theme.primary' })
  const secondary = useWatch({ control, name: 'overview.theme.secondary' })

  const [mode, setMode] = useState(ColorMode.Dark)

  const baseTheme = mode === ColorMode.Light ? lightTheme : darkTheme
  const theme = createAccentTheme(baseTheme, { primary, secondary }) ?? baseTheme

  return { mode, setMode, theme }
}
