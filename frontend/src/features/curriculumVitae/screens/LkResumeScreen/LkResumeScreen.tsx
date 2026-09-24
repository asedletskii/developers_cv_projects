import { useEffect } from 'react'
import {
  BackToResumes,
  CurriculumVitae,
  LoadFailure,
  LoadPlaceholder,
} from '@features/curriculumVitae/components'
import {
  clearCurriculumVitaeState,
  getPersonData,
  selectIsPersonDataLoading,
  selectPersonDataError,
} from '@features/curriculumVitae/store/curriculumVitae'
import { useAppDispatch, useAppSelector } from '@shared/store'
import { useParams } from 'react-router'

/**
 * Предпросмотр резюме внутри ЛК: рендер тот же, что у публичной `/cv/:cvId`, но данные берутся
 * из кабинета владельца. Редактирование — на `/cv/lk/:userId/:cvId`.
 */
export default function LkResumeScreen() {
  const dispatch = useAppDispatch()
  const { userId, cvId } = useParams()

  const isLoading = useAppSelector(selectIsPersonDataLoading)
  const error = useAppSelector(selectPersonDataError)

  useEffect(() => {
    if (!cvId || !userId) return

    const request = dispatch(getPersonData({ cvId, userId }))

    return () => {
      request.abort()
    }
  }, [dispatch, cvId, userId])

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
