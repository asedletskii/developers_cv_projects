import { glassCard } from '@features/curriculumVitae/styles'
import { Box, Container, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { STATS } from './constants'
import { useCountUp } from './hooks/useCountUp'

const Stat = ({ value, label }: { value: number; label: string }) => {
  const countRef = useCountUp(value)

  return (
    <Stack
      className="reveal-item"
      spacing={0.5}
      alignItems="center"
      textAlign="center"
      sx={{ p: { xs: 2, md: 3 }, ...glassCard }}
    >
      <Typography
        component="span"
        sx={{
          fontSize: 'clamp(2.5rem, 6vw, 4rem)',
          fontWeight: 800,
          lineHeight: 1,
          color: 'primary.main',
          fontVariantNumeric: 'tabular-nums',
        }}
      >
        <span ref={countRef}>{value}</span>
      </Typography>

      <Typography sx={{ color: 'text.secondary', fontSize: '0.95rem' }}>{label}</Typography>
    </Stack>
  )
}

export const Stats = () => {
  const translate = useTranslate('AboutPage.Stats')

  return (
    <Container maxWidth="lg" sx={{ py: { xs: 4, md: 6 } }}>
      <Box
        className="reveal-group"
        sx={{
          display: 'grid',
          gridTemplateColumns: { xs: 'repeat(2, 1fr)', md: 'repeat(4, 1fr)' },
          gap: 2.5,
        }}
      >
        {STATS.map(({ key, value }) => (
          <Stat key={key} value={value} label={translate(key)} />
        ))}
      </Box>
    </Container>
  )
}
