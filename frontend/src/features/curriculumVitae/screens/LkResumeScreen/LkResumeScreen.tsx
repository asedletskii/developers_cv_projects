import {
  BackToResumes,
  CurriculumVitae,
  LoadFailure,
  LoadPlaceholder,
} from '@features/curriculumVitae/components'
import { usePublicResume } from '@features/curriculumVitae/hooks'
import { useParams } from 'react-router'

/**
 * Предпросмотр резюме внутри ЛК: рендер тот же, что у публичной `/cv/:cvId`, но данные читаются
 * владельческой ручкой (`userId` в `usePublicResume`) — иначе черновик или резюме, снятое с
 * публикации, не открылось бы даже собственному владельцу. Редактирование — на `/cv/lk/:userId/:cvId`.
 */
const LkResumeScreen = () => {
  const { userId, cvId } = useParams()
  const { isLoading, error } = usePublicResume(cvId, userId)

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
