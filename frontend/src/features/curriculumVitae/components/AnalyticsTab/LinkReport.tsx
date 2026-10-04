import { useState } from 'react'
import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { useAnalyticsRequest } from '@features/curriculumVitae/hooks'
import { AnalyticsPeriod, Resume } from '@features/curriculumVitae/types'
import { getAnalyticsRange } from '@features/curriculumVitae/utils'
import { LinearProgress, Stack, ToggleButton, ToggleButtonGroup } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { Breakdown } from './Breakdown'
import { ClickMap } from './ClickMap'
import { DailyChart } from './DailyChart'
import { EmptyReport } from './EmptyReport'
import { ReportError } from './ReportError'
import { SectionFunnel } from './SectionFunnel'
import { Summary } from './Summary'
import { TechnologyInterest } from './TechnologyInterest'
import { VisitFeed } from './VisitFeed'

type Props = {
  resume: Resume
  resumeId: string
  linkId: string
  period: AnalyticsPeriod
}

type View = 'overview' | 'clickMap'

export const LinkReport = ({ resume, resumeId, linkId, period }: Props) => {
  const translate = useTranslate('Lk.Analytics.Tabs')
  const [view, setView] = useState<View>('overview')
  const { data, isLoading, hasError, reload } = useAnalyticsRequest(
    async (signal) => {
      const range = getAnalyticsRange(period)
      const [statistics, daily] = await Promise.all([
        curriculumVitaeApi.getLinkStatistics(linkId, range, signal),
        curriculumVitaeApi.getLinkDaily(linkId, range, signal),
      ])
      return { statistics, daily }
    },
    [linkId, period],
  )

  if (isLoading) return <LinearProgress />

  if (hasError || !data) return <ReportError onRetry={reload} />

  const isEmpty = data.statistics.views === 0 && !data.statistics.firstViewAt

  return (
    <Stack spacing={2.5}>
      <Summary statistics={data.statistics} />

      {isEmpty ? (
        <EmptyReport />
      ) : (
        <>
          <ToggleButtonGroup
            exclusive
            size="small"
            value={view}
            onChange={(_, value: View | null) => value && setView(value)}
            sx={{ alignSelf: 'flex-start' }}
          >
            {(['overview', 'clickMap'] as const).map((key) => (
              <ToggleButton key={key} value={key} sx={{ textTransform: 'none', px: 2 }}>
                {translate(key)}
              </ToggleButton>
            ))}
          </ToggleButtonGroup>

          {view === 'overview' ? (
            <>
              <DailyChart days={data.daily} />
              <Breakdown statistics={data.statistics} />
              <SectionFunnel resume={resume} linkId={linkId} period={period} />
              <TechnologyInterest linkId={linkId} period={period} />
              <VisitFeed linkId={linkId} />
            </>
          ) : (
            <ClickMap linkId={linkId} resumeId={resumeId} />
          )}
        </>
      )}
    </Stack>
  )
}
