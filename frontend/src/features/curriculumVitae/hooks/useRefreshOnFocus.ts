import { useEffect, useRef } from 'react'

const DEFAULT_MIN_HIDDEN_MS = 60_000

export const useRefreshOnFocus = (
  onRefresh: () => void,
  minHiddenMs: number = DEFAULT_MIN_HIDDEN_MS,
) => {
  const onRefreshRef = useRef(onRefresh)

  useEffect(() => {
    onRefreshRef.current = onRefresh
  }, [onRefresh])

  useEffect(() => {
    let hiddenAt: number | null = null

    const onVisibilityChange = () => {
      if (document.visibilityState === 'hidden') {
        hiddenAt = Date.now()
        return
      }

      if (hiddenAt !== null && Date.now() - hiddenAt >= minHiddenMs) onRefreshRef.current()
      hiddenAt = null
    }

    document.addEventListener('visibilitychange', onVisibilityChange)
    return () => document.removeEventListener('visibilitychange', onVisibilityChange)
  }, [minHiddenMs])
}
