import { ButtonBase, Stack } from '@mui/material'
import { useLanguage, useTranslate } from '@shared/hooks'
import { LANGUAGES } from './constants'

/** Переключатель языка интерфейса: сегментированная «таблетка» RU / EN. */
export const LanguageSwitch = () => {
  const translate = useTranslate('LanguageSwitch')
  const { language, changeLanguage } = useLanguage()

  return (
    <Stack
      direction="row"
      role="group"
      aria-label={translate('label')}
      sx={{ p: '2px', borderRadius: 999, bgcolor: 'action.hover' }}
    >
      {LANGUAGES.map((code) => {
        const isActive = code === language

        return (
          <ButtonBase
            key={code}
            aria-pressed={isActive}
            onClick={() => changeLanguage(code)}
            sx={{
              minWidth: 36,
              height: 30,
              borderRadius: 999,
              fontSize: '0.75rem',
              fontWeight: 600,
              letterSpacing: '0.04em',
              color: isActive ? 'primary.contrastText' : 'gray.light',
              bgcolor: isActive ? 'primary.main' : 'transparent',
              transition: 'background-color 200ms, color 200ms',
              '&:hover': { color: isActive ? 'primary.contrastText' : 'primary.main' },
            }}
          >
            {code.toUpperCase()}
          </ButtonBase>
        )
      })}
    </Stack>
  )
}
