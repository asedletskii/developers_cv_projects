import { AuthShell } from '@features/auth/components/AuthShell'
import { RISE, SHAKE, STILL } from '@features/auth/components/AuthShell/constants'
import { glassCard, secondaryButton } from '@features/curriculumVitae/styles'
import { HomeRounded, RefreshRounded, ReportProblemRounded } from '@mui/icons-material'
import { Box, Button, Chip, Stack, Typography } from '@mui/material'
import { ROUTES } from '@shared/constants'
import { useTranslate } from '@shared/hooks'

type Props = {
  error: Error
}

export const ErrorFallback = ({ error }: Props) => {
  const translate = useTranslate('ErrorBoundary')

  return (
    <AuthShell>
      <Stack
        role="alert"
        spacing={2.5}
        alignItems="center"
        textAlign="center"
        sx={{
          maxWidth: 480,
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
          <ReportProblemRounded sx={{ fontSize: 40 }} />
        </Box>

        <Typography variant="h4" component="h1">
          {translate('title')}
        </Typography>
        <Typography sx={{ color: 'text.secondary' }}>{translate('text')}</Typography>

        {error.message && (
          <Chip
            variant="outlined"
            label={error.message.slice(0, 160)}
            sx={{
              maxWidth: '100%',
              height: 'auto',
              py: 0.5,
              fontFamily: 'monospace',
              '& .MuiChip-label': { whiteSpace: 'normal', px: 1.5 },
            }}
          />
        )}

        <Stack spacing={1.5} sx={{ width: '100%' }}>
          <Button
            variant="contained"
            size="large"
            startIcon={<RefreshRounded />}
            onClick={() => window.location.reload()}
            sx={{ borderRadius: 999, textTransform: 'none' }}
          >
            {translate('reload')}
          </Button>

          {/* Обычная ссылка, а не роутер: перезагрузка сбрасывает и сломанное состояние приложения */}
          <Button
            href={ROUTES.home}
            variant="outlined"
            size="large"
            startIcon={<HomeRounded />}
            sx={{ borderRadius: 999, textTransform: 'none', ...secondaryButton }}
          >
            {translate('home')}
          </Button>
        </Stack>
      </Stack>
    </AuthShell>
  )
}
