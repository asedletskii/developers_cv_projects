import { Stack, Typography } from '@mui/material'

type Props = {
  title: string
  hasBackButton: boolean
}

export const Title = ({ title, hasBackButton }: Props) => (
  <Stack
    direction={'row'}
    spacing={2}
    sx={{
      maxWidth: 'max-content',
      mt: hasBackButton ? 1 : 0,
      mb: 3,
      alignItems: 'center',
    }}
  >
    <Typography noWrap title={title} sx={{ typography: 'h3' }}>
      {title}
    </Typography>
  </Stack>
)
