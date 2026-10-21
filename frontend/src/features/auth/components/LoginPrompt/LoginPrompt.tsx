import { authApi } from '@features/auth/api'
import { AuthShell } from '@features/auth/components/AuthShell'
import { RISE, STILL } from '@features/auth/components/AuthShell/constants'
import { badge, iconTile } from '@features/curriculumVitae/styles'
import {
  ArrowForwardRounded,
  DashboardCustomizeRounded,
  InsightsRounded,
  Telegram as TelegramIcon,
} from '@mui/icons-material'
import { Box, Button, Chip, Container, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { LoginFlow } from './LoginFlow'

const BENEFITS = [
  { key: 'benefitNoPassword', Icon: TelegramIcon },
  { key: 'benefitSections', Icon: DashboardCustomizeRounded },
  { key: 'benefitAnalytics', Icon: InsightsRounded },
] as const

export const LoginPrompt = () => {
  const translate = useTranslate('Auth')

  return (
    <AuthShell>
      <Container maxWidth="md">
        <Box
          sx={{
            display: 'grid',
            gridTemplateColumns: { xs: '1fr', md: '1.15fr 1fr' },
            gap: { xs: 4, md: 6 },
            alignItems: 'center',
            animation: `${RISE} 0.6s ease-out both`,
            ...STILL,
          }}
        >
          <Stack spacing={3} alignItems={{ xs: 'center', md: 'flex-start' }}>
            <Chip label={translate('loginTitle')} variant="outlined" sx={badge} />

            <Typography
              variant="h2"
              sx={{
                fontSize: 'clamp(2rem, 5vw, 3.25rem)',
                lineHeight: 1.1,
                textAlign: { xs: 'center', md: 'left' },
              }}
            >
              {translate('loginHeadline')}
            </Typography>

            <Typography sx={{ color: 'text.secondary', textAlign: { xs: 'center', md: 'left' } }}>
              {translate('loginText')}
            </Typography>

            <Stack component="ul" spacing={1.25} sx={{ m: 0, p: 0, listStyle: 'none' }}>
              {BENEFITS.map(({ key, Icon }) => (
                <Stack key={key} component="li" direction="row" spacing={1.5} alignItems="center">
                  <Stack
                    alignItems="center"
                    justifyContent="center"
                    sx={{
                      width: 36,
                      height: 36,
                      borderRadius: 2,
                      ...iconTile,
                      flexShrink: 0,
                    }}
                  >
                    <Icon fontSize="small" />
                  </Stack>
                  <Typography>{translate(key)}</Typography>
                </Stack>
              ))}
            </Stack>

            <Button
              variant="contained"
              size="large"
              startIcon={<TelegramIcon />}
              endIcon={<ArrowForwardRounded />}
              href={authApi.getTelegramLoginUrl()}
              sx={{ borderRadius: 999, px: 4, py: 1.5, textTransform: 'none', fontSize: '1rem' }}
            >
              {translate('login')}
            </Button>
          </Stack>

          <LoginFlow />
        </Box>
      </Container>
    </AuthShell>
  )
}
