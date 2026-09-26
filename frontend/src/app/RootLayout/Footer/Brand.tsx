import { Box, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { getBrandGradient } from '@shared/styles/gradients'

type Props = {
  appName: string
}

/** Название проекта крупным словом: последнее слово — акцентным градиентом, ниже — описание. */
export const Brand = ({ appName }: Props) => {
  const translate = useTranslate('Footer')

  const words = appName.split(' ')
  const accent = words.pop()

  return (
    <Stack spacing={1.5} sx={{ maxWidth: 400, flexShrink: 0 }}>
      <Typography variant="h2" sx={{ fontWeight: 800, letterSpacing: '-0.02em', lineHeight: 1.05 }}>
        {words.join(' ')}{' '}
        <Box
          component="span"
          sx={(theme) => ({
            background: getBrandGradient(theme),
            WebkitBackgroundClip: 'text',
            backgroundClip: 'text',
            WebkitTextFillColor: 'transparent',
          })}
        >
          {accent}
        </Box>
      </Typography>

      <Typography sx={{ color: 'gray.light', whiteSpace: 'pre-line' }}>
        {translate('info')}
      </Typography>
    </Stack>
  )
}
