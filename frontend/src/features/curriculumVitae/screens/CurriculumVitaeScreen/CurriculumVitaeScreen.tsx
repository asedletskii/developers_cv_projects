import { useMemo } from 'react'
import { CurriculumVitae, LoadFailure, LoadPlaceholder } from '@features/curriculumVitae/components'
import { AnalyticsNotice } from '@features/curriculumVitae/components/AnalyticsNotice'
import { AnalyticsProvider } from '@features/curriculumVitae/components/AnalyticsProvider'
import { usePublicResume } from '@features/curriculumVitae/hooks'
import { getVisibleSections } from '@features/curriculumVitae/sections'
import { selectResume, selectViewId } from '@features/curriculumVitae/store/curriculumVitae'
import { useAppSelector } from '@shared/store'
import { useParams } from 'react-router'

const CurriculumVitaeScreen = () => {
  const { cvId } = useParams()
  const { isLoading, error } = usePublicResume(cvId)
  const resume = useAppSelector(selectResume)
  const viewId = useAppSelector(selectViewId)
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
