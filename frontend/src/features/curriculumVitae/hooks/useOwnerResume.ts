import { useCallback, useEffect, useState } from 'react'
import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { CurriculumVitaeError, ResumeView } from '@features/curriculumVitae/types'
import { parseApiError } from '@features/curriculumVitae/utils'

/**
 * Резюме владельца целиком (`GET /lk/resume/:id`). `reload` перечитывает его без сброса экрана в
 * загрузку — так после сохранения секции значения в форме обновляются на месте.
 */
export const useOwnerResume = (userId?: string, resumeId?: string) => {
  const [view, setView] = useState<ResumeView | null>(null)
  const [error, setError] = useState<CurriculumVitaeError | null>(null)

  const reload = useCallback(
    async (signal?: AbortSignal) => {
      if (!userId || !resumeId) return

      try {
        setView(await curriculumVitaeApi.getResume(userId, resumeId, signal))
        setError(null)
      } catch (e) {
        if (!signal?.aborted) setError(parseApiError(e))
      }
    },
    [userId, resumeId],
  )

  useEffect(() => {
    setView(null)

    const controller = new AbortController()
    reload(controller.signal)

    return () => controller.abort()
  }, [reload])

  return { view, error, setError, reload }
}
