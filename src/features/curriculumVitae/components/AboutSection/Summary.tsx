import { Stack, Typography } from '@mui/material'

type Props = {
  fio: string
  position: string
}

export const Summary = ({ fio, position }: Props) => (
  <Stack spacing={0.5} alignItems={{ xs: 'center', sm: 'flex-start' }}>
    <Typography variant="h4" sx={{ fontSize: { xs: '1.5rem', sm: '1.75rem' } }}>
      {fio}
    </Typography>

    <Typography variant="h6" sx={{ color: 'primary.main' }}>
      {position}
    </Typography>
  </Stack>
)
