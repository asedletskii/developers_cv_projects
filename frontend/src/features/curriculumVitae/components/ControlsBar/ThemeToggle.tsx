import { DarkModeOutlined, LightModeOutlined } from '@mui/icons-material'
import { IconButton } from '@mui/material'
import { useColorMode, useTranslate } from '@shared/hooks'
import { ColorMode } from '@shared/types'

export const ThemeToggle = () => {
  const translate = useTranslate('CurriculumVitae.Controls')
  const { colorMode, changeColorMode } = useColorMode()

  const isDark = colorMode === ColorMode.Dark

  const toggle = () => {
    const root = document.documentElement

    root.classList.add('theme-transition')
    window.setTimeout(() => root.classList.remove('theme-transition'), 500)

    changeColorMode(isDark ? ColorMode.Light : ColorMode.Dark)
  }

  return (
    <IconButton
      aria-label={translate('theme')}
      title={translate('theme')}
      onClick={toggle}
      sx={{ width: 40, height: 40 }}
    >
      {isDark ? <LightModeOutlined /> : <DarkModeOutlined />}
    </IconButton>
  )
}
