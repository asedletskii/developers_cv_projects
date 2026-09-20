import { useCallback, useState } from 'react'
import { logout } from '@features/auth/store/session'
import { useAppDispatch } from '@shared/store'

export const useLogout = () => {
  const dispatch = useAppDispatch()
  const [isLoggingOut, setIsLoggingOut] = useState(false)

  const signOut = useCallback(async () => {
    setIsLoggingOut(true)
    await dispatch(logout())
    setIsLoggingOut(false)
  }, [dispatch])

  return { signOut, isLoggingOut }
}
