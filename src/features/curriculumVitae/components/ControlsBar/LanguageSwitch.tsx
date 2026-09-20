import { ButtonBase, Stack } from '@mui/material'
import { useLanguage, useTranslate } from '@shared/hooks'
import { Languages } from '@shared/types'

const LANGUAGES: Languages[] = [Languages.Ru, Languages.En]

export const LanguageSwitch = () => {
  const translate = useTranslate('CurriculumVitae.Controls')
  const { language, changeLanguage } = useLanguage()

  return (
    <Stack
      direction="row"
      role="group"
      aria-label={translate('language')}
      sx={{ p: '2px', borderRadius: 999, bgcolor: 'action.hover' }}
    >
      {LANGUAGES.map((lang) => {
        const isActive = lang === language

        return (
          <ButtonBase
            key={lang}
            aria-pressed={isActive}
            onClick={() => changeLanguage(lang)}
            sx={{
              minWidth: 36,
              height: 34,
              borderRadius: 999,
              fontSize: '0.75rem',
              fontWeight: 600,
              letterSpacing: '0.04em',
              color: isActive ? 'primary.contrastText' : 'text.primary',
              bgcolor: isActive ? 'primary.main' : 'transparent',
              transition: 'background-color 200ms, color 200ms, scale 120ms ease-out',
            }}
          >
            {lang.toUpperCase()}
          </ButtonBase>
        )
      })}
    </Stack>
  )
}
