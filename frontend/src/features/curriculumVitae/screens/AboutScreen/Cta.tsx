import { SparkleField } from '@features/curriculumVitae/components/cv/SparkleField'
import { ArrowForwardRounded } from '@mui/icons-material'
import { Button, Container, Stack, Typography } from '@mui/material'
import { ROUTES } from '@shared/constants'
import { useTranslate } from '@shared/hooks'
import { Link as RouterLink } from 'react-router'
import { DRIFT } from './constants'

export const Cta = () => {
  const translate = useTranslate('AboutPage.Cta')

  return (
    <Container maxWidth="lg" sx={{ pb: { xs: 8, md: 12 } }}>
      <Stack
        className="reveal"
        spacing={3}
        alignItems="center"
        textAlign="center"
        sx={{
          position: 'relative',
          overflow: 'hidden',
          p: { xs: 4, md: 8 },
          borderRadius: 4,
          color: 'common.black',
          animation: `${DRIFT} 12s ease-in-out infinite`,
          '@media (prefers-reduced-motion: reduce)': { animation: 'none' },
          background: ({ palette }) =>
            `linear-gradient(135deg, ${palette.primary.light}, ${palette.primary.main} 55%, ${palette.primary.dark})`,
          backgroundSize: '200% 200%',
        }}
      >
        <SparkleField count={20} tone="white" seed={5} />

        <Typography variant="h3" sx={{ color: 'inherit', fontSize: 'clamp(1.75rem, 4vw, 2.5rem)' }}>
          {translate('title')}
        </Typography>
        <Typography sx={{ color: 'inherit', maxWidth: 480, opacity: 0.8 }}>
          {translate('text')}
        </Typography>

        <Button
          component={RouterLink}
          to={ROUTES.lk}
          size="large"
          endIcon={<ArrowForwardRounded />}
          sx={{
            borderRadius: 999,
            px: 4,
            py: 1.5,
            textTransform: 'none',
            fontSize: '1rem',
            bgcolor: 'common.white',
            color: 'common.black',
            '&:hover': { bgcolor: 'common.white', transform: 'translateY(-2px)' },
            transition: 'transform .2s',
          }}
        >
          {translate('button')}
        </Button>
      </Stack>
    </Container>
  )
}
