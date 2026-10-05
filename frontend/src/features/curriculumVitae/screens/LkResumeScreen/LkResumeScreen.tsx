import { useSession } from '@features/auth/hooks'
import {
  BackToResumes,
  CurriculumVitae,
  LoadFailure,
  LoadPlaceholder,
} from '@features/curriculumVitae/components'
import { usePublicResume } from '@features/curriculumVitae/hooks'
import { useParams, useSearchParams } from 'react-router'

const LkResumeScreen = () => {
  const { cvId } = useParams()
  // `?embed` — предпросмотр внутри iframe на вкладке аналитики: без плашки «назад»
  const isEmbedded = useSearchParams()[0].has('embed')
  const { userId, isLoading: isSessionLoading } = useSession()
  const { isLoading: isResumeLoading, error } = usePublicResume(cvId, userId ?? undefined)

  if (isSessionLoading || isResumeLoading) return <LoadPlaceholder />
  if (error) return <LoadFailure error={error} />

  return (
    <>
      {!isEmbedded && <BackToResumes cvId={cvId} />}

      <CurriculumVitae />
    </>
  )
}

export default LkResumeScreen
