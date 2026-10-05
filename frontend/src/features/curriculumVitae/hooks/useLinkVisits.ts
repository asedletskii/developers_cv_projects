import { useCallback, useEffect, useState } from 'react'
import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { LinkVisit } from '@features/curriculumVitae/types'

export const useLinkVisits = (linkId: string, includeFiltered = false) => {
  const [visits, setVisits] = useState<LinkVisit[]>([])
  const [cursor, setCursor] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [hasError, setHasError] = useState(false)
  const [version, setVersion] = useState(0)

  const load = useCallback(
    async (from: string | null, signal?: AbortSignal) => {
      setIsLoading(true)
      setHasError(false)

      try {
        const page = await curriculumVitaeApi.getLinkVisits(linkId, from, includeFiltered, signal)
        setVisits((current) => (from ? [...current, ...page.items] : page.items))
        setCursor(page.nextCursor)
      } catch {
        if (!signal?.aborted) setHasError(true)
      } finally {
        if (!signal?.aborted) setIsLoading(false)
      }
    },
    [linkId, includeFiltered],
  )

  useEffect(() => {
    const controller = new AbortController()
    setVisits([])
    setCursor(null)
    load(null, controller.signal)

    return () => controller.abort()
  }, [load, version])

  return {
    visits,
    isLoading,
    hasError,
    hasMore: !!cursor,
    loadMore: () => load(cursor),
    reload: () => setVersion((current) => current + 1),
  }
}
