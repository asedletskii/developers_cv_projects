import { useEffect } from 'react'
import {
  clearCurriculumVitaeState,
  getPublicResume,
  selectIsResumeLoading,
  selectResumeError,
} from '@features/curriculumVitae/store/curriculumVitae'
import { useAppDispatch, useAppSelector } from '@shared/store'

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
