import { glassCard } from '@features/curriculumVitae/styles'
import { Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'

/**TODO: пока заглушка — на бэке нет ручек со статистикой просмотров. */
export const AnalyticsTab = () => {
  const translate = useTranslate('Lk.Analytics')

  return (
    <Stack
      spacing={1}
      alignItems="center"
      sx={{ p: { xs: 4, sm: 6 }, textAlign: 'center', ...glassCard }}
    >
      <Typography variant="h6">{translate('title')}</Typography>

      <Typography sx={{ color: 'text.secondary' }}>{translate('text')}</Typography>
    </Stack>
  )
}
