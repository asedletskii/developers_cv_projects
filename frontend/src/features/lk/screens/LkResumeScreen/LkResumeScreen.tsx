import { useEffect } from 'react'
import { CurriculumVitae, LoadFailure, LoadPlaceholder } from '@features/curriculumVitae/components'
import {
  clearCurriculumVitaeState,
  getPersonData,
  selectIsPersonDataLoading,
  selectPersonDataError,
} from '@features/curriculumVitae/store/curriculumVitae'
import { BackToResumes } from '@features/lk/components'
import { useLanguage } from '@shared/hooks'
import { useAppDispatch, useAppSelector } from '@shared/store'
import { useParams } from 'react-router'

/**
 * Предпросмотр резюме внутри ЛК: то, как оно выглядит на публичной странице. Данные и рендер —
 * те же, что у `/cv/:cvId` (переиспользуем фичу curriculumVitae целиком), редактирование — на
 * `/cv/lk/:userId/:cvId`.
 */
export default function LkResumeScreen() {
  const dispatch = useAppDispatch()
  const { userId, cvId } = useParams()
  const { language } = useLanguage()

  const isLoading = useAppSelector(selectIsPersonDataLoading)
  const error = useAppSelector(selectPersonDataError)

  useEffect(() => {
    if (!cvId) return

    const request = dispatch(getPersonData({ cvId, lang: language }))

    return () => {
      request.abort()
    }
  }, [dispatch, cvId, language])

  useEffect(() => {
    return () => {
      dispatch(clearCurriculumVitaeState())
    }
  }, [dispatch])

  if (isLoading) return <LoadPlaceholder />

  if (error) return <LoadFailure error={error} />

  return (
    <>
      {userId && <BackToResumes userId={userId} cvId={cvId} />}

      <CurriculumVitae />
    </>
  )
}
