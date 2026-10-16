import { Stack, Typography } from '@mui/material'

type Props = {
  items: string[]
}

export const Achievements = ({ items }: Props) => (
  <Stack component="ul" spacing={0.5} sx={{ m: 0, pl: 2.5 }}>
    {items.map((item) => (
      <Typography key={item} component="li" sx={{ '&::marker': { color: 'primary.main' } }}>
        {item}
      </Typography>
    ))}
  </Stack>
)
