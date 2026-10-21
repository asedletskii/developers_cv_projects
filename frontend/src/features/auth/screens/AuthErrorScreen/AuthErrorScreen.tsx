import { authApi } from '@features/auth/api'
import { AuthShell } from '@features/auth/components/AuthShell'
import { RISE, SHAKE, STILL } from '@features/auth/components/AuthShell/constants'
import { glassCard, secondaryButton } from '@features/curriculumVitae/styles'
import { ErrorOutlineRounded } from '@mui/icons-material'
import { Box, Button, Stack, Typography } from '@mui/material'
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
    <AuthShell>
      <Stack
        role="alert"
        spacing={2.5}
        alignItems="center"
        textAlign="center"
        sx={{
          maxWidth: 440,
          mx: 'auto',
          p: { xs: 3, md: 5 },
          ...glassCard,
          animation: `${RISE} 0.5s ease-out both`,
          ...STILL,
        }}
      >
        <Box
          sx={{
            display: 'grid',
            placeItems: 'center',
            width: 72,
            height: 72,
            borderRadius: '50%',
            bgcolor: 'action.hover',
            color: 'error.main',
            animation: `${SHAKE} 0.6s ease-in-out 0.3s both`,
            ...STILL,
          }}
        >
          <ErrorOutlineRounded sx={{ fontSize: 40 }} />
        </Box>

        <Typography variant="h4">{translate('errorTitle')}</Typography>
        <Typography sx={{ color: 'text.secondary' }}>{translate(textKey)}</Typography>

        <Stack direction={{ xs: 'column', sm: 'row' }} spacing={1.5} sx={{ width: '100%' }}>
          <Button
            href={authApi.getTelegramLoginUrl()}
            variant="contained"
            size="large"
            sx={{ flex: 1, borderRadius: 999, textTransform: 'none' }}
          >
            {translate('retry')}
          </Button>

          <Button
            component={RouterLink}
            to={ROUTES.home}
            variant="outlined"
            size="large"
            sx={{ flex: 1, borderRadius: 999, textTransform: 'none', ...secondaryButton }}
          >
            {translate('backHome')}
          </Button>
        </Stack>
      </Stack>
    </AuthShell>
  )
}

export default AuthErrorScreen
