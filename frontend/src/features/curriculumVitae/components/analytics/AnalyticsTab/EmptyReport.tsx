import { glassCard } from '@features/curriculumVitae/styles'
import { Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'

export const EmptyReport = () => {
  const translate = useTranslate('Lk.Analytics')

  return (
    <Stack spacing={0.5} alignItems="center" textAlign="center" sx={{ p: 4, ...glassCard }}>
      <Typography variant="h6">{translate('emptyTitle')}</Typography>
      <Typography sx={{ color: 'text.secondary' }}>{translate('emptyText')}</Typography>
    </Stack>
  )
}
