import { Stack, Typography } from '@mui/material'

type Props = {
  title: string
  subtitle: string
}

export const Headline = ({ title, subtitle }: Props) => {
  if (!title && !subtitle) return null

  return (
    <Stack spacing={1} sx={{ position: 'relative' }}>
      {title && (
        <Typography variant="h3" sx={{ color: 'inherit', fontSize: { xs: '1.75rem', md: '2rem' } }}>
          {title}
        </Typography>
      )}

      {subtitle && <Typography sx={{ color: 'inherit', opacity: 0.9 }}>{subtitle}</Typography>}
    </Stack>
  )
}
