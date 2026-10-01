import { useEffect } from 'react'
import {
  getLkResumes,
  selectLkResumes,
  selectLkResumesError,
} from '@features/curriculumVitae/store/lk'
import { useAppDispatch, useAppSelector } from '@shared/store'

export const useLkResumes = (userId?: string) => {
  const dispatch = useAppDispatch()
  const resumes = useAppSelector(selectLkResumes)
  const error = useAppSelector(selectLkResumesError)

  useEffect(() => {
    if (!userId) return

    const request = dispatch(getLkResumes(userId))
    return () => {
      request.abort()
    }
  }, [dispatch, userId])

  return { resumes, error }
}
