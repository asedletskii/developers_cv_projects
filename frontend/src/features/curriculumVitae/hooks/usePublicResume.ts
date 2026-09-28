import { useEffect } from 'react'
import {
  clearCurriculumVitaeState,
  getPublicResume,
  selectIsResumeLoading,
  selectResumeError,
} from '@features/curriculumVitae/store/curriculumVitae'
import { useAppDispatch, useAppSelector } from '@shared/store'

/**
 * Загружает резюме в хранилище и убирает его при уходе с экрана. Смена id отменяет предыдущий
 * запрос: иначе при быстром переключении более медленный старый ответ может прийти позже нового
 * и перезаписать экран устаревшими данными.
 *
 * `userId` передаёт только предпросмотр в ЛК — тогда резюме читается владельческой ручкой
 * (`GET /lk/resume/:id`), а не публичной: иначе черновик или снятое с публикации резюме
 * не открылись бы даже своему владельцу (публичная ручка отвечает `404`).
 */
export const usePublicResume = (cvId?: string, userId?: string) => {
  const dispatch = useAppDispatch()
  const isLoading = useAppSelector(selectIsResumeLoading)
  const error = useAppSelector(selectResumeError)

  useEffect(() => {
    if (!cvId) return

    const request = dispatch(getPublicResume({ cvId, userId }))
    return () => {
      request.abort()
    }
  }, [dispatch, cvId, userId])

  useEffect(() => {
    return () => {
      dispatch(clearCurriculumVitaeState())
    }
  }, [dispatch])

  return { isLoading, error }
}
