import { glassCard } from '@features/curriculumVitae/styles'
import { Box, Container, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { FEATURES } from './constants'

export const Features = () => {
  const translate = useTranslate('AboutPage.Features')

  return (
    <Container maxWidth="lg" sx={{ py: { xs: 4, md: 10 } }}>
      <Stack
        spacing={1}
        alignItems="center"
        textAlign="center"
        className="reveal"
        sx={{ mb: { xs: 3, md: 6 } }}
      >
        <Typography variant="h3" sx={{ fontSize: 'clamp(1.75rem, 4vw, 2.5rem)' }}>
          {translate('title')}
        </Typography>
        <Typography sx={{ color: 'text.secondary' }}>{translate('subtitle')}</Typography>
      </Stack>

      <Box
        className="reveal-group"
        sx={{
          display: 'grid',
          gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)' },
          gap: 2.5,
        }}
      >
        {FEATURES.map(({ key, Icon }) => (
          <Stack
            key={key}
            className="reveal-item"
            spacing={1.5}
            sx={{
              p: 3,
              ...glassCard,
              transition: 'transform .25s ease, border-color .25s ease',
              '&:hover': { transform: 'translateY(-4px)', borderColor: 'primary.main' },
              '@media (prefers-reduced-motion: reduce)': { '&:hover': { transform: 'none' } },
            }}
          >
            <Stack
              alignItems="center"
              justifyContent="center"
              sx={{
                width: 48,
                height: 48,
                borderRadius: 2,
                bgcolor: 'action.hover',
                color: 'primary.main',
              }}
            >
              <Icon />
            </Stack>

            <Typography variant="h6">{translate(`${key}.title`)}</Typography>
            <Typography sx={{ color: 'text.secondary' }}>{translate(`${key}.text`)}</Typography>
          </Stack>
        ))}
      </Box>
    </Container>
  )
}
