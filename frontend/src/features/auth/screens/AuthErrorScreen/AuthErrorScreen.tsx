import { authApi } from '@features/auth/api'
import { Button, Stack, Typography } from '@mui/material'
import { ROUTES } from '@shared/constants'
import { useTranslate } from '@shared/hooks'
import { Link as RouterLink, useSearchParams } from 'react-router'

const REASON_TEXT_KEY = {
  provider_denied: 'reasonProviderDenied',
  authentication_failed: 'reasonAuthenticationFailed',
} as const

const AuthErrorScreen = () => {
  const translate = useTranslate('Auth')
  const [searchParams] = useSearchParams()
  const reason = searchParams.get('reason') ?? ''
  const textKey = REASON_TEXT_KEY[reason as keyof typeof REASON_TEXT_KEY] ?? 'errorText'

  return (
    <Stack
      spacing={2}
      alignItems="center"
      justifyContent="center"
      textAlign="center"
      sx={{ flex: 1, minHeight: '70vh', px: 2 }}
    >
      <Typography variant="h3">{translate('errorTitle')}</Typography>
      <Typography sx={{ color: 'text.secondary' }}>{translate(textKey)}</Typography>

      <Stack direction="row" spacing={2}>
        <Button
          href={authApi.getTelegramLoginUrl()}
          variant="contained"
          sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}
        >
          {translate('retry')}
        </Button>

        <Button
          component={RouterLink}
          to={ROUTES.home}
          variant="outlined"
          sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}
        >
          {translate('backHome')}
        </Button>
      </Stack>
    </Stack>
  )
}

export default AuthErrorScreen
