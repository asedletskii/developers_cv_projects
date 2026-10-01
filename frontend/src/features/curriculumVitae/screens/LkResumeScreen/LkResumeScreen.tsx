import { useSession } from '@features/auth/hooks'
import {
  BackToResumes,
  CurriculumVitae,
  LoadFailure,
  LoadPlaceholder,
} from '@features/curriculumVitae/components'
import { usePublicResume } from '@features/curriculumVitae/hooks'
import { useParams } from 'react-router'

const LkResumeScreen = () => {
  const { cvId } = useParams()
  const { userId, isLoading: isSessionLoading } = useSession()
  const { isLoading: isResumeLoading, error } = usePublicResume(cvId, userId ?? undefined)

  if (isSessionLoading || isResumeLoading) return <LoadPlaceholder />
  if (error) return <LoadFailure error={error} />

  return (
    <>
      <BackToResumes cvId={cvId} />

      <CurriculumVitae />
    </>
  )
}

export default LkResumeScreen
