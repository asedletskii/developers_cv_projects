import { useEffect } from 'react'
import {
  clearLkState,
  getLkProfile,
  selectIsLkProfileLoading,
  selectLkProfile,
  selectLkProfileError,
} from '@features/curriculumVitae/store/lk'
import { useAppDispatch, useAppSelector } from '@shared/store'

/** Общие данные владельца кабинета: загрузка при входе и очистка состояния кабинета при выходе. */
export const useLkProfile = (userId?: string) => {
  const dispatch = useAppDispatch()
  const isLoading = useAppSelector(selectIsLkProfileLoading)
  const error = useAppSelector(selectLkProfileError)
  const profile = useAppSelector(selectLkProfile)

  useEffect(() => {
    if (!userId) return

    const request = dispatch(getLkProfile(userId))
    return () => {
      request.abort()
    }
  }, [dispatch, userId])

  useEffect(() => {
    return () => {
      dispatch(clearLkState())
    }
  }, [dispatch])

  return { profile, isLoading, error }
}
