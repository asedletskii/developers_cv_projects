import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { useAnalyticsRequest } from '@features/curriculumVitae/hooks'
import { AnalyticsPeriod } from '@features/curriculumVitae/types'
import { getAnalyticsRange, getPreviousRange } from '@features/curriculumVitae/utils'
import { LinearProgress, Stack } from '@mui/material'
import { LinksBreakdown } from './LinksBreakdown'
import { ReportError } from './ReportError'
import { Summary } from './Summary'

type Props = {
  resumeId: string
  period: AnalyticsPeriod
  onOpenLink: (linkId: string) => void
}

export const AllLinksReport = ({ resumeId, period, onOpenLink }: Props) => {
  const { data, isLoading, hasError, reload } = useAnalyticsRequest(
    async (signal) => {
      const [statistics, previous] = await Promise.all([
        curriculumVitaeApi.getResumeStatistics(resumeId, getAnalyticsRange(period), signal),
        curriculumVitaeApi
          .getResumeStatistics(resumeId, getPreviousRange(period), signal)
          .catch(() => undefined),
      ])
      return { statistics, previous }
    },
    [resumeId, period],
  )

  if (isLoading) return <LinearProgress />

  if (hasError || !data) return <ReportError onRetry={reload} />

  return (
    <Stack spacing={2.5}>
      <Summary statistics={data.statistics} previous={data.previous} period={period} />
      <LinksBreakdown links={data.statistics.links} onOpenLink={onOpenLink} />
    </Stack>
  )
}
