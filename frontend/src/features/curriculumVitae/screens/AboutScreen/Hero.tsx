import { useGreetingAnimation } from '@features/curriculumVitae/components/cv/HeroSection/hooks/useGreetingAnimation'
import { SparkleField } from '@features/curriculumVitae/components/cv/SparkleField'
import { ArrowForwardRounded } from '@mui/icons-material'
import { Box, Button, Chip, Container, Stack, Typography } from '@mui/material'
import { splitText } from '@shared/animations'
import { ROUTES } from '@shared/constants'
import { useTranslate } from '@shared/hooks'
import { Link as RouterLink } from 'react-router'

export const Hero = () => {
  const translate = useTranslate('AboutPage.Hero')
  const rootRef = useGreetingAnimation()

  return (
    <Box
      ref={rootRef}
      sx={{
        position: 'relative',
        display: 'flex',
        alignItems: 'center',
        minHeight: { xs: 'auto', md: '80svh' },
        py: { xs: 10, md: 0 },
        overflow: 'hidden',
        textAlign: 'center',
      }}
    >
      <Box
        aria-hidden
        sx={{
          position: 'absolute',
          inset: 0,
          background: ({ palette }) =>
            `radial-gradient(640px circle at 50% 0%, ${palette.primary.main}22, transparent 70%)`,
        }}
      />
      <SparkleField count={32} seed={11} />

      <Container maxWidth="md" sx={{ position: 'relative', zIndex: 1 }}>
        <Stack spacing={3} alignItems="center">
          <Chip className="name" label={translate('kicker')} color="primary" variant="outlined" />

          <Typography
            variant="h1"
            sx={{ fontSize: 'clamp(2.25rem, 6vw, 4.25rem)', lineHeight: 1.08 }}
          >
            {splitText(translate('title'), 'letter')}
          </Typography>

          <Typography
            className="subtitle"
            sx={{ color: 'text.secondary', fontSize: 'clamp(1rem, 2vw, 1.25rem)' }}
          >
            {translate('subtitle')}
          </Typography>

          <Button
            className="hero-action"
            component={RouterLink}
            to={ROUTES.lk}
            variant="contained"
            size="large"
            endIcon={<ArrowForwardRounded />}
            sx={{ borderRadius: 999, px: 4, py: 1.5, textTransform: 'none', fontSize: '1rem' }}
          >
            {translate('cta')}
          </Button>
        </Stack>
      </Container>
    </Box>
  )
}
