import { useState } from 'react'
import { glassCard } from '@features/curriculumVitae/styles'
import { Person } from '@features/curriculumVitae/types'
import { createAccentTheme } from '@features/curriculumVitae/utils'
import { Brightness4Rounded, Brightness7Rounded } from '@mui/icons-material'
import {
  Button,
  ButtonBase,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Stack,
  Typography,
} from '@mui/material'
import { ThemeProvider } from '@mui/material/styles'
import { useTranslate } from '@shared/hooks'
import { darkTheme } from '@shared/styles/darkTheme'
import { lightTheme } from '@shared/styles/lightTheme'
import { ColoredTypography } from '@shared/components'
import { ColorMode } from '@shared/types'
import { useFormContext, useWatch } from 'react-hook-form'

type Props = {
  onClose: () => void
}

const MODES = [ColorMode.Light, ColorMode.Dark]

/**
 * Предпросмотр цветов резюме: переключатель светлая/тёмная — только для этой модалки, на
 * остальное приложение не влияет. Тема считается из текущих (ещё не сохранённых) значений формы.
 * Тумблер оставляем в оформлении модалки (внешняя тема), а вложенная тема красит только карточку
 * превью — так же, как LanguageSwitch остаётся своим цветом вне резюме, которое переключает.
 */
export const ThemePreviewDialog = ({ onClose }: Props) => {
  const translate = useTranslate('Lk.Edit.Theme')
  const { control } = useFormContext<Person>()
  const primary = useWatch({ control, name: 'theme.primary' })
  const secondary = useWatch({ control, name: 'theme.secondary' })

  const [mode, setMode] = useState(ColorMode.Dark)

  const baseTheme = mode === ColorMode.Light ? lightTheme : darkTheme
  const theme = createAccentTheme(baseTheme, { primary, secondary }) ?? baseTheme

  return (
    <Dialog
      open
      onClose={onClose}
      fullWidth
      maxWidth="xs"
      slotProps={{ paper: { sx: { ...glassCard, backgroundImage: 'none' } } }}
    >
      <DialogTitle
        sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 2 }}
      >
        {translate('previewTitle')}

        <Stack
          direction="row"
          role="group"
          aria-label={translate('title')}
          sx={{ p: '2px', borderRadius: 999, bgcolor: 'action.hover' }}
        >
          {MODES.map((option) => {
            const isActive = option === mode
            const Icon = option === ColorMode.Light ? Brightness7Rounded : Brightness4Rounded

            return (
              <ButtonBase
                key={option}
                aria-pressed={isActive}
                aria-label={option === ColorMode.Light ? translate('lightMode') : translate('darkMode')}
                onClick={() => setMode(option)}
                sx={{
                  width: 34,
                  height: 34,
                  borderRadius: 999,
                  color: isActive ? 'primary.contrastText' : 'text.primary',
                  bgcolor: isActive ? 'primary.main' : 'transparent',
                  transition: 'background-color 200ms, color 200ms, scale 120ms ease-out',
                }}
              >
                <Icon fontSize="small" />
              </ButtonBase>
            )
          })}
        </Stack>
      </DialogTitle>

      <ThemeProvider theme={theme}>
        <DialogContent sx={{ bgcolor: 'background.default' }}>
          <Stack spacing={2} alignItems="center" sx={{ p: 3, ...glassCard }}>
            <ColoredTypography
              variant="h5"
              firstPart="Иван"
              coloredPart="Иванов"
              color="primary.main"
            />

            <Typography sx={{ color: 'text.secondary' }}>Frontend-разработчик</Typography>

            <Stack direction="row" spacing={1} flexWrap="wrap" justifyContent="center">
              <Chip label="React" color="primary" />
              <Chip
                label="TypeScript"
                variant="outlined"
                sx={{ borderColor: 'secondary.main', color: 'secondary.main' }}
              />
            </Stack>

            <Button variant="contained" sx={{ borderRadius: 999, textTransform: 'none' }}>
              {translate('previewCta')}
            </Button>
          </Stack>
        </DialogContent>
      </ThemeProvider>

      <DialogActions sx={{ px: 3, pb: 3 }}>
        <Button onClick={onClose} sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}>
          {translate('close')}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
