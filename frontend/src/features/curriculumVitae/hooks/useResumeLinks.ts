import { useCallback, useEffect, useState } from 'react'
import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { LinkStatus, ResumeLink } from '@features/curriculumVitae/types'
import { ToastType } from '@shared/types'
import { showToast } from '@shared/utils'

export const useResumeLinks = (userId?: string, resumeId?: string) => {
  const [links, setLinks] = useState<ResumeLink[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [hasLoadError, setHasLoadError] = useState(false)

  const fetchLinks = useCallback(
    (signal?: AbortSignal) => {
      if (!userId || !resumeId) return Promise.resolve(undefined)
      return curriculumVitaeApi.getResumeLinks(resumeId, signal)
    },
    [userId, resumeId],
  )

  const load = useCallback(
    async (signal?: AbortSignal) => {
      setIsLoading(true)
      setHasLoadError(false)

      try {
        const result = await fetchLinks(signal)
        if (result) setLinks(result)
      } catch {
        if (!signal?.aborted) setHasLoadError(true)
      } finally {
        if (!signal?.aborted) setIsLoading(false)
      }
    },
    [fetchLinks],
  )

  useEffect(() => {
    const controller = new AbortController()
    load(controller.signal)

    return () => controller.abort()
  }, [load])

  const refresh = async () => {
    const result = await fetchLinks()
    if (result) setLinks(result)
  }

  const create = async (name: string) => {
    if (!userId || !resumeId) return

    try {
      await curriculumVitaeApi.createResumeLink(resumeId, name)
      await refresh()
    } catch {
      showToast({ messageKey: 'linkActionFailed', type: ToastType.Error })
    }
  }

  const rename = async (id: string, version: number, name: string) => {
    if (!userId) return

    try {
      await curriculumVitaeApi.renameResumeLink(id, version, name)
      await refresh()
    } catch {
      showToast({ messageKey: 'linkActionFailed', type: ToastType.Error })
    }
  }

  const remove = async (id: string) => {
    if (!userId) return

    try {
      await curriculumVitaeApi.deleteResumeLink(id)
      setLinks((current) => current.filter((link) => link.id !== id))
    } catch {
      showToast({ messageKey: 'linkActionFailed', type: ToastType.Error })
    }
  }

  const changeStatus = async (id: string, version: number, status: LinkStatus) => {
    if (!userId) return

    try {
      await curriculumVitaeApi.changeLinkStatus(id, version, status)
      await refresh()
    } catch {
      showToast({ messageKey: 'linkActionFailed', type: ToastType.Error })
    }
  }

  return {
    links,
    isLoading,
    hasLoadError,
    reload: () => load(),
    create,
    rename,
    remove,
    changeStatus,
  }
}
