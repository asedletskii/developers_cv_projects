import { useEffect } from 'react'
import {
  fetchSession,
  selectSessionStatus,
  selectSessionUserId,
} from '@features/auth/store/session'
import { useAppDispatch, useAppSelector } from '@shared/store'

export const useSession = () => {
  const dispatch = useAppDispatch()
  const status = useAppSelector(selectSessionStatus)
  const userId = useAppSelector(selectSessionUserId)

  useEffect(() => {
    if (status === 'idle') dispatch(fetchSession())
  }, [dispatch, status])

  return {
    userId,
    status,
    isLoading: status === 'idle' || status === 'loading',
    isAuthenticated: status === 'authenticated',
  }
}
