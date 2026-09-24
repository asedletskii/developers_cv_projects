import { useEffect } from 'react'
import { CurriculumVitae, LoadFailure, LoadPlaceholder } from '@features/curriculumVitae/components'
import {
  clearCurriculumVitaeState,
  getPersonData,
  selectIsPersonDataLoading,
  selectPersonDataError,
} from '@features/curriculumVitae/store/curriculumVitae'
import { useAppDispatch, useAppSelector } from '@shared/store'
import { useParams } from 'react-router'

export default function CurriculumVitaeScreen() {
  const dispatch = useAppDispatch()
  const { cvId } = useParams()

  const isLoading = useAppSelector(selectIsPersonDataLoading)
  const error = useAppSelector(selectPersonDataError)

  // Смена id отменяет предыдущий запрос: иначе при быстром переключении более медленный
  // старый ответ может прийти позже нового и перезаписать экран устаревшими данными.
  useEffect(() => {
    if (!cvId) return

    const request = dispatch(getPersonData({ cvId }))

    return () => {
      request.abort()
    }
  }, [dispatch, cvId])

  useEffect(() => {
    return () => {
      dispatch(clearCurriculumVitaeState())
    }
  }, [dispatch])

  if (isLoading) return <LoadPlaceholder />

  if (error) return <LoadFailure error={error} />

  return <CurriculumVitae />
}
