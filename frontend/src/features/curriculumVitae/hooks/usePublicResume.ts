import { useEffect } from 'react'
import {
  clearCurriculumVitaeState,
  getPersonData,
  selectIsPersonDataLoading,
  selectPersonDataError,
} from '@features/curriculumVitae/store/curriculumVitae'
import { useAppDispatch, useAppSelector } from '@shared/store'

/**
 * Загружает публичное резюме (`GET /cv/:id`) в хранилище и убирает его при уходе с экрана.
 * Смена id отменяет предыдущий запрос: иначе при быстром переключении более медленный старый
 * ответ может прийти позже нового и перезаписать экран устаревшими данными.
 */
export const usePublicResume = (cvId?: string) => {
  const dispatch = useAppDispatch()
  const isLoading = useAppSelector(selectIsPersonDataLoading)
  const error = useAppSelector(selectPersonDataError)

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

  return { isLoading, error }
}
