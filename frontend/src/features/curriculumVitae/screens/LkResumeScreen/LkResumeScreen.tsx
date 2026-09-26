import {
  BackToResumes,
  CurriculumVitae,
  LoadFailure,
  LoadPlaceholder,
} from '@features/curriculumVitae/components'
import { usePublicResume } from '@features/curriculumVitae/hooks'
import { useParams } from 'react-router'

/**
 * Предпросмотр резюме внутри ЛК: рендер тот же, что у публичной `/cv/:cvId`. Редактирование —
 * на `/cv/lk/:userId/:cvId`.
 */
const LkResumeScreen = () => {
  const { userId, cvId } = useParams()
  const { isLoading, error } = usePublicResume(cvId)

  if (isLoading) return <LoadPlaceholder />
  if (error) return <LoadFailure error={error} />

  return (
    <>
      {userId && <BackToResumes userId={userId} cvId={cvId} />}

      <CurriculumVitae />
    </>
  )
}

export default LkResumeScreen
