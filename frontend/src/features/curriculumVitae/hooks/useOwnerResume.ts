import { useCallback, useEffect, useState } from 'react'
import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { CurriculumVitaeError, ResumeView } from '@features/curriculumVitae/types'
import { parseApiError } from '@features/curriculumVitae/utils'

/**
 * Резюме владельца целиком (`GET /lk/resume/:id`). `reload` перечитывает его без сброса экрана в
 * загрузку — так после сохранения секции значения в форме обновляются на месте.
 */
export const useOwnerResume = (userId?: string, resumeId?: string) => {
  // Резюме хранится вместе с id, к которому относится: при смене id старое не показываем, и
  // сбрасывать состояние в эффекте не нужно
  const [loaded, setLoaded] = useState<{ resumeId: string; view: ResumeView } | null>(null)
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
    // Загрузка с бэка: setState вызывается после ответа, а не синхронно в теле эффекта
    // eslint-disable-next-line react-hooks/set-state-in-effect
    reload(controller.signal)

    return () => controller.abort()
  }, [reload])

  const view = loaded && loaded.resumeId === resumeId ? loaded.view : null

  return { view, error, setError, reload }
}
