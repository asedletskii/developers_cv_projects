import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { useAnalyticsRequest } from '@features/curriculumVitae/hooks'
import { glassCard } from '@features/curriculumVitae/styles'
import { AnalyticsPeriod } from '@features/curriculumVitae/types'
import { getAnalyticsRange } from '@features/curriculumVitae/utils'
import { Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { BarList } from './BarList'

type Props = {
  linkId: string
  period: AnalyticsPeriod
}

export const TechnologyInterest = ({ linkId, period }: Props) => {
  const translate = useTranslate('Lk.Analytics.Tech')
  const { data } = useAnalyticsRequest(
    (signal) => curriculumVitaeApi.getLinkTechnologies(linkId, getAnalyticsRange(period), signal),
    [linkId, period],
  )

  // Блок необязательный: пока ручка не отвечает, просто не показываем его
  if (!data) return null

  return (
    <Stack spacing={2} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
      <Stack spacing={0.25}>
        <Typography variant="h6">{translate('title')}</Typography>
        <Typography variant="caption" sx={{ color: 'text.secondary' }}>
          {translate('hint')}
        </Typography>
      </Stack>

      {data.length === 0 ? (
        <Typography sx={{ color: 'text.secondary' }}>{translate('empty')}</Typography>
      ) : (
        <BarList
          rows={data.map(({ name, selections, visitors }) => ({
            key: name,
            label: name,
            value: selections,
            caption: translate('visitors', { count: visitors }),
          }))}
        />
      )}
    </Stack>
  )
}
