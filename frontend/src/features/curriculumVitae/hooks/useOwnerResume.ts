import { useCallback, useEffect, useState } from 'react'
import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { CurriculumVitaeError, Resume } from '@features/curriculumVitae/types'
import { parseApiError } from '@features/curriculumVitae/utils'

export const useOwnerResume = (userId?: string, resumeId?: string) => {
  const [loaded, setLoaded] = useState<{ resumeId: string; view: Resume } | null>(null)
  const [error, setError] = useState<CurriculumVitaeError | null>(null)

  const reload = useCallback(
    async (signal?: AbortSignal) => {
      if (!userId || !resumeId) return

      try {
        const view = await curriculumVitaeApi.getResume(userId, resumeId, signal)
        setLoaded({ resumeId, view })
        setError(null)
      } catch (e) {
        if (!signal?.aborted) setError(parseApiError(e))
      }
    },
    [userId, resumeId],
  )

  useEffect(() => {
    const controller = new AbortController()
    // eslint-disable-next-line react-hooks/set-state-in-effect
    reload(controller.signal)

    return () => controller.abort()
  }, [reload])

  const view = loaded && loaded.resumeId === resumeId ? loaded.view : null

  return { view, error, setError, reload }
}
