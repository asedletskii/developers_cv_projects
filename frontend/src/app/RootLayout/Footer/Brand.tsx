import { Box, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  appName: string
}

export const Brand = ({ appName }: Props) => {
  const translate = useTranslate('Footer')

  const words = appName.split(' ')
  const accent = words.pop()

  return (
    <Stack spacing={1.5} sx={{ maxWidth: 400, flexShrink: 0 }}>
      <Typography variant="h2" sx={{ fontWeight: 800, letterSpacing: '-0.02em', lineHeight: 1.05 }}>
        {words.join(' ')}{' '}
        <Box component="span" sx={{ color: 'primary.main' }}>
          {accent}
        </Box>
      </Typography>

      <Typography sx={{ color: 'gray.light', whiteSpace: 'pre-line' }}>
        {translate('info')}
      </Typography>
    </Stack>
  )
}
