import { SparkleField } from '@features/curriculumVitae/components/cv/SparkleField'
import { secondaryButton } from '@features/curriculumVitae/styles'
import { ArrowBackRounded, HomeRounded, LinkOffRounded } from '@mui/icons-material'
import { Box, Button, Chip, Stack, Typography } from '@mui/material'
import { ROUTES } from '@shared/constants'
import { useDocumentTitle, useTranslate } from '@shared/hooks'
import { Link as RouterLink, useLocation, useNavigate } from 'react-router'
import { BREATHE, ORBIT, POP, RISE, STILL, WOBBLE } from './constants'

const DIGIT_SX = {
  fontSize: 'clamp(6rem, 22vw, 13rem)',
  fontWeight: 800,
  lineHeight: 1,
  color: 'primary.main',
  letterSpacing: '-0.04em',
} as const

const NotFoundScreen = () => {
  const translate = useTranslate('NotFound')
  const { pathname } = useLocation()
  const navigate = useNavigate()
  useDocumentTitle(translate('title'))

  const canGoBack = Boolean(window.history.state?.idx)

  return (
    <Box
      sx={{
        position: 'relative',
        flex: 1,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '70vh',
        px: 2,
        py: { xs: 6, md: 8 },
        overflow: 'hidden',
      }}
    >
      <Box
        aria-hidden
        sx={{
          position: 'absolute',
          inset: 0,
          maskImage: 'linear-gradient(to bottom, transparent, #000 30%, #000 70%, transparent)',
          WebkitMaskImage:
            'linear-gradient(to bottom, transparent, #000 30%, #000 70%, transparent)',
          background: ({ palette }) =>
            `radial-gradient(760px circle at 50% 42%, ${palette.primary.main}26, transparent 70%)`,
        }}
      />
      <SparkleField count={28} seed={404} />

      <Stack spacing={3} alignItems="center" textAlign="center" sx={{ position: 'relative' }}>
        <Stack direction="row" alignItems="center" role="img" aria-label="404">
          <Typography
            aria-hidden
            component="span"
            sx={{ ...DIGIT_SX, animation: `${POP} 0.8s ease-out both`, ...STILL }}
          >
            4
          </Typography>

          <Box
            aria-hidden
            sx={{
              position: 'relative',
              display: 'grid',
              placeItems: 'center',
              width: 'clamp(5rem, 18vw, 10.5rem)',
              height: 'clamp(5rem, 18vw, 10.5rem)',
              mx: { xs: 0.5, md: 1.5 },
              animation: `${POP} 0.8s ease-out 0.15s both`,
              ...STILL,
            }}
          >
            <Box
              sx={{
                position: 'absolute',
                inset: 0,
                borderRadius: '50%',
                border: 'clamp(10px, 2.4vw, 22px) solid',
                borderColor: 'primary.main',
                animation: `${BREATHE} 3.2s ease-in-out infinite`,
                ...STILL,
              }}
            />
            <Box
              sx={{
                position: 'absolute',
                inset: 0,
                animation: `${ORBIT} 5s linear infinite`,
                ...STILL,
              }}
            >
              <Box
                sx={{
                  position: 'absolute',
                  top: -6,
                  left: '50%',
                  ml: '-6px',
                  width: 12,
                  height: 12,
                  borderRadius: '50%',
                  bgcolor: 'primary.light',
                  boxShadow: ({ palette }) => `0 0 14px ${palette.primary.main}`,
                }}
              />
            </Box>

            <LinkOffRounded
              sx={{
                fontSize: 'clamp(2rem, 7vw, 4rem)',
                color: 'text.secondary',
                animation: `${WOBBLE} 3s ease-in-out infinite`,
                ...STILL,
              }}
            />
          </Box>

          <Typography
            aria-hidden
            component="span"
            sx={{ ...DIGIT_SX, animation: `${POP} 0.8s ease-out 0.3s both`, ...STILL }}
          >
            4
          </Typography>
        </Stack>

        <Stack
          spacing={1.5}
          alignItems="center"
          sx={{ animation: `${RISE} 0.6s ease-out 0.7s both`, ...STILL }}
        >
          <Typography variant="h4" component="h1">
            {translate('title')}
          </Typography>

          <Typography sx={{ color: 'text.secondary', maxWidth: 460 }}>
            {translate('text')}
          </Typography>

          <Chip
            label={pathname}
            variant="outlined"
            sx={{ maxWidth: '100%', fontFamily: 'monospace', '& .MuiChip-label': { px: 1.5 } }}
          />
        </Stack>

        <Stack
          direction={{ xs: 'column', sm: 'row' }}
          spacing={1.5}
          sx={{
            width: { xs: '100%', sm: 'auto' },
            animation: `${RISE} 0.6s ease-out 0.9s both`,
            ...STILL,
          }}
        >
          <Button
            component={RouterLink}
            to={ROUTES.home}
            variant="contained"
            size="large"
            startIcon={<HomeRounded />}
            sx={{ borderRadius: 999, px: 4, textTransform: 'none' }}
          >
            {translate('home')}
          </Button>

          {canGoBack && (
            <Button
              variant="outlined"
              size="large"
              startIcon={<ArrowBackRounded />}
              onClick={() => navigate(-1)}
              sx={{ borderRadius: 999, px: 4, textTransform: 'none', ...secondaryButton }}
            >
              {translate('back')}
            </Button>
          )}
        </Stack>
      </Stack>
    </Box>
  )
}

export default NotFoundScreen
