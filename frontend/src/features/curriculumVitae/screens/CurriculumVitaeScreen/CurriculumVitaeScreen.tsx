import { useMemo } from 'react'
import { CurriculumVitae, LoadFailure, LoadPlaceholder } from '@features/curriculumVitae/components'
import { AnalyticsNotice } from '@features/curriculumVitae/components/analytics/AnalyticsNotice'
import { AnalyticsProvider } from '@features/curriculumVitae/components/analytics/AnalyticsProvider'
import { usePublicResume } from '@features/curriculumVitae/hooks'
import { getVisibleSections } from '@features/curriculumVitae/sections'
import { selectResume, selectViewId } from '@features/curriculumVitae/store/curriculumVitae'
import { useDocumentTitle, useNoIndex } from '@shared/hooks'
import { useAppSelector } from '@shared/store'
import { useParams } from 'react-router'

const CurriculumVitaeScreen = () => {
  const { cvId } = useParams()
  const { isLoading, error } = usePublicResume(cvId)
  const resume = useAppSelector(selectResume)
  const viewId = useAppSelector(selectViewId)
  const profile = resume.profile
  useNoIndex()
  useDocumentTitle(
    profile?.name
      ? `${profile.name} ${profile.surname}${profile.position ? ` — ${profile.position}` : ''}`
      : undefined,
  )
  const sectionOrder = useMemo(() => getVisibleSections(resume).map(({ id }) => id), [resume])

  if (isLoading) return <LoadPlaceholder />
  if (error) return <LoadFailure error={error} />

  return (
    <AnalyticsProvider linkId={cvId} viewId={viewId} sectionOrder={sectionOrder}>
      <CurriculumVitae />
      <AnalyticsNotice />
    </AnalyticsProvider>
  )
}

export default CurriculumVitaeScreen
