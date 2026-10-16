import { glassCard } from '@features/curriculumVitae/styles'
import { ErrorOutlineRounded, RefreshRounded } from '@mui/icons-material'
import { Button, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  onRetry: () => void
}

export const LinksLoadError = ({ onRetry }: Props) => {
  const translate = useTranslate('Lk.Links')

  return (
    <Stack
      spacing={1.5}
      alignItems="center"
      textAlign="center"
      sx={{ p: { xs: 3, sm: 4 }, ...glassCard }}
    >
      <ErrorOutlineRounded color="error" fontSize="large" />

      <Typography sx={{ fontWeight: 600 }}>{translate('loadErrorTitle')}</Typography>

      <Typography variant="body2" sx={{ color: 'text.secondary' }}>
        {translate('loadError')}
      </Typography>

      <Button
        startIcon={<RefreshRounded />}
        onClick={onRetry}
        variant="outlined"
        sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}
      >
        {translate('retry')}
      </Button>
    </Stack>
  )
}
