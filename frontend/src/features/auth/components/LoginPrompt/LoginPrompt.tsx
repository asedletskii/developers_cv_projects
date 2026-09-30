import { authApi } from '@features/auth/api'
import { glassCard } from '@features/curriculumVitae/styles'
import { Telegram as TelegramIcon } from '@mui/icons-material'
import { Button, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'

export const LoginPrompt = () => {
  const translate = useTranslate('Auth')

  return (
    <Stack
      spacing={2}
      alignItems="center"
      justifyContent="center"
      textAlign="center"
      sx={{ flex: 1, minHeight: '70vh', px: 2 }}
    >
      <Stack spacing={2} alignItems="center" sx={{ p: 4, maxWidth: 420, ...glassCard }}>
        <Typography variant="h5">{translate('loginTitle')}</Typography>
        <Typography sx={{ color: 'text.secondary' }}>{translate('loginText')}</Typography>

        <Button
          variant="contained"
          startIcon={<TelegramIcon />}
          href={authApi.getTelegramLoginUrl()}
          sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}
        >
          {translate('login')}
        </Button>
      </Stack>
    </Stack>
  )
}
