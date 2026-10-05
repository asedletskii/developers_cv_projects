import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { useAnalyticsRequest } from '@features/curriculumVitae/hooks'
import { getVisibleSections } from '@features/curriculumVitae/sections'
import { glassCard } from '@features/curriculumVitae/styles'
import { AnalyticsPeriod, Resume, SectionId } from '@features/curriculumVitae/types'
import { getAnalyticsRange, orderSectionStats, toPercent } from '@features/curriculumVitae/utils'
import { Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { BarList } from './BarList'

type Props = {
  resume: Resume
  linkId: string
  period: AnalyticsPeriod
}

export const SectionFunnel = ({ resume, linkId, period }: Props) => {
  const translate = useTranslate('Lk.Analytics.Funnel')
  const translateSection = useTranslate('Lk.Edit.Sections')
  const { data } = useAnalyticsRequest(
    (signal) => curriculumVitaeApi.getLinkSections(linkId, getAnalyticsRange(period), signal),
    [linkId, period],
  )

  if (!data || data.reads === 0) return null

  return (
    <Stack spacing={2} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
      <Stack spacing={0.25}>
        <Typography variant="h6">{translate('title')}</Typography>
        <Typography variant="caption" sx={{ color: 'text.secondary' }}>
          {translate('hint')}
        </Typography>
      </Stack>

      <BarList
        max={data.reads}
        rows={orderSectionStats(
          data.sections,
          getVisibleSections(resume).map(({ id }) => id),
        ).map(({ section, reached, exits }) => ({
          key: section,
          label: translateSection(section as SectionId),
          value: reached,
          caption: `${toPercent(reached, data.reads) ?? 0}% · ${translate('exits', { count: exits })}`,
        }))}
      />
    </Stack>
  )
}
