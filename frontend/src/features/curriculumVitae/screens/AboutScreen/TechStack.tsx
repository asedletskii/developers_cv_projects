import { Box, Chip, Container, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { MARQUEE, TECH_ROWS } from './constants'

const Row = ({ names, reverse }: { names: string[]; reverse: boolean }) => (
  <Box
    sx={{
      overflow: 'hidden',
      maskImage: 'linear-gradient(90deg, transparent, #000 12%, #000 88%, transparent)',
      WebkitMaskImage: 'linear-gradient(90deg, transparent, #000 12%, #000 88%, transparent)',
      '&:hover > div': { animationPlayState: 'paused' },
    }}
  >
    <Stack
      direction="row"
      spacing={1.25}
      sx={{
        width: 'max-content',
        animation: `${MARQUEE} 26s linear infinite ${reverse ? 'reverse' : ''}`,
        '@media (prefers-reduced-motion: reduce)': {
          animation: 'none',
          width: 'auto',
          flexWrap: 'wrap',
          justifyContent: 'center',
          rowGap: 1.25,
          '& > [aria-hidden]': { display: 'none' },
        },
      }}
    >
      {[...names, ...names, ...names, ...names].map((name, index) => (
        <Chip
          key={index}
          label={name}
          aria-hidden={index >= names.length ? true : undefined}
          sx={{ fontSize: '0.95rem', px: 1, py: 2.5 }}
        />
      ))}
    </Stack>
  </Box>
)

export const TechStack = () => {
  const translate = useTranslate('AboutPage.Stack')

  return (
    <Container maxWidth="lg" sx={{ py: { xs: 4, md: 6 }, textAlign: 'center' }}>
      <Typography variant="h5" className="reveal" sx={{ mb: 3 }}>
        {translate('title')}
      </Typography>

      <Stack className="reveal" spacing={1.5}>
        {TECH_ROWS.map((names, index) => (
          <Row key={index} names={names} reverse={index % 2 === 1} />
        ))}
      </Stack>
    </Container>
  )
}
