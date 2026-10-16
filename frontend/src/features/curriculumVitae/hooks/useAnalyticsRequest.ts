import { useCallback, useEffect, useRef, useState } from 'react'

export const useAnalyticsRequest = <T>(
  fetcher: (signal: AbortSignal) => Promise<T>,
  deps: unknown[],
) => {
  const [data, setData] = useState<T>()
  const [isLoading, setIsLoading] = useState(true)
  const [hasError, setHasError] = useState(false)
  const [version, setVersion] = useState(0)
  const fetcherRef = useRef(fetcher)
  fetcherRef.current = fetcher

  useEffect(() => {
    const controller = new AbortController()
    setIsLoading(true)
    setHasError(false)

    fetcherRef
      .current(controller.signal)
      .then(setData)
      .catch(() => !controller.signal.aborted && setHasError(true))
      .finally(() => !controller.signal.aborted && setIsLoading(false))

    return () => controller.abort()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [...deps, version])

  const reload = useCallback(() => setVersion((current) => current + 1), [])

  return { data, isLoading, hasError, reload }
}
