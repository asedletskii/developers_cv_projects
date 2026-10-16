import { LoadFailure } from '@features/curriculumVitae/components'
import { AnalyticsTab } from '@features/curriculumVitae/components/analytics/AnalyticsTab'
import { ResumeEditor } from '@features/curriculumVitae/components/editor/ResumeEditor'
import { LinksTab } from '@features/curriculumVitae/components/lk/LinksTab'
import { MainTab } from '@features/curriculumVitae/components/lk/MainTab'
import { ResumeHeader } from '@features/curriculumVitae/components/lk/ResumeHeader'
import {
  useLkOutletContext,
  useOwnerResume,
  useRefreshOnFocus,
  useResumeLinks,
  useResumeStatistics,
} from '@features/curriculumVitae/hooks'
import { LinearProgress, Stack, Tab, Tabs } from '@mui/material'
import { LocalizationProvider } from '@mui/x-date-pickers'
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs'
import { useLanguage, useTranslate } from '@shared/hooks'
import { useParams } from 'react-router'
import { TAB_SX, TOP_SECTIONS } from './constants'
import { useResumeSections } from './hooks'

const LkResumeEditScreen = () => {
  const { userId } = useLkOutletContext()
  const { cvId } = useParams()
  const { language } = useLanguage()
  const translateSections = useTranslate('Lk.Sections')
  const { view, error, setError, reload } = useOwnerResume(userId, cvId)
  const { section, changeSection, openEditTab, isDirtyRef } = useResumeSections()
  const linksState = useResumeLinks(userId, view?.id)
  const { data: statistics, reload: reloadStatistics } = useResumeStatistics(view?.id)

  useRefreshOnFocus(reloadStatistics)

  if (error) return <LoadFailure error={error} />
  if (!view || !userId) return <LinearProgress />

  return (
    <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale={language}>
      <Stack spacing={3} sx={{ width: '100%', minWidth: 0 }}>
        <ResumeHeader view={view} onReload={reload} onError={setError} />

        <Tabs
          value={section}
          onChange={(_, id) => changeSection(id)}
          variant="scrollable"
          scrollButtons="auto"
        >
          {TOP_SECTIONS.map((id) => (
            <Tab key={id} value={id} label={translateSections(id)} sx={TAB_SX} />
          ))}
        </Tabs>

        {section === 'main' && (
          <MainTab
            view={view}
            linksState={linksState}
            statistics={statistics}
            onOpenTab={openEditTab}
            onOpenAnalytics={() => changeSection('analytics')}
          />
        )}

        {section === 'edit' && (
          <ResumeEditor view={view} onSaved={reload} isDirtyRef={isDirtyRef} />
        )}

        {section === 'links' && <LinksTab linksState={linksState} statistics={statistics} />}

        {section === 'analytics' && (
          <AnalyticsTab resume={view} resumeId={view.id} linksState={linksState} />
        )}
      </Stack>
    </LocalizationProvider>
  )
}

export default LkResumeEditScreen
