import { DarkModeOutlined, LightModeOutlined } from '@mui/icons-material'
import { IconButton, SxProps, Theme } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useThemeToggle } from './hooks'

type Props = {
  sx?: SxProps<Theme>
}

export const ThemeToggle = ({ sx }: Props) => {
  const translate = useTranslate('ThemeToggle')
  const { isDark, toggle } = useThemeToggle()

  return (
    <IconButton
      aria-label={translate('label')}
      title={translate('label')}
      onClick={toggle}
      sx={[{ width: 40, height: 40 }, ...(Array.isArray(sx) ? sx : [sx])]}
    >
      {isDark ? <LightModeOutlined /> : <DarkModeOutlined />}
    </IconButton>
  )
}
