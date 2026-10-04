import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { useAnalyticsRequest } from '@features/curriculumVitae/hooks'
import { AnalyticsPeriod } from '@features/curriculumVitae/types'
import { getAnalyticsRange } from '@features/curriculumVitae/utils'
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
    (signal) => curriculumVitaeApi.getResumeStatistics(resumeId, getAnalyticsRange(period), signal),
    [resumeId, period],
  )

  if (isLoading) return <LinearProgress />

  if (hasError || !data) return <ReportError onRetry={reload} />

  return (
    <Stack spacing={2.5}>
      <Summary statistics={data} />
      <LinksBreakdown links={data.links} onOpenLink={onOpenLink} />
    </Stack>
  )
}
