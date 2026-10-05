import { useState } from 'react'
import { CreateLinkPrompt } from '@features/curriculumVitae/components/CreateLinkPrompt'
import { LinksLoadError } from '@features/curriculumVitae/components/LinksTab/LinksLoadError'
import { useRefreshOnFocus, useResumeLinks } from '@features/curriculumVitae/hooks'
import { AnalyticsPeriod, Resume } from '@features/curriculumVitae/types'
import { LinearProgress, Stack } from '@mui/material'
import { AllLinksReport } from './AllLinksReport'
import { Controls } from './Controls'
import { LinkReport } from './LinkReport'

type Props = {
  resume: Resume
  resumeId: string
  linksState: ReturnType<typeof useResumeLinks>
}

const ALL_LINKS = 'all'

export const AnalyticsTab = ({ resume, resumeId, linksState }: Props) => {
  const [scope, setScope] = useState(ALL_LINKS)
  const [period, setPeriod] = useState<AnalyticsPeriod>(30)
  const [refreshKey, setRefreshKey] = useState(0)
  const refresh = () => setRefreshKey((key) => key + 1)

  useRefreshOnFocus(refresh)

  if (linksState.isLoading) return <LinearProgress />

  if (linksState.hasLoadError) return <LinksLoadError onRetry={linksState.reload} />

  if (linksState.links.length === 0) return <CreateLinkPrompt onCreate={linksState.create} />

  const reportKey = `${scope}-${period}-${refreshKey}`

  return (
    <Stack spacing={2.5}>
      <Controls
        links={linksState.links}
        scope={scope}
        allValue={ALL_LINKS}
        period={period}
        onScopeChange={setScope}
        onPeriodChange={setPeriod}
        onRefresh={refresh}
      />

      {scope === ALL_LINKS ? (
        <AllLinksReport key={reportKey} resumeId={resumeId} period={period} onOpenLink={setScope} />
      ) : (
        <LinkReport
          key={reportKey}
          resume={resume}
          resumeId={resumeId}
          linkId={scope}
          period={period}
        />
      )}
    </Stack>
  )
}
