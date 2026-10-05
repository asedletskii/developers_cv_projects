import { useState } from 'react'
import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { useAnalytics } from '@features/curriculumVitae/hooks/useAnalytics'
import { AnalyticsEventType } from '@features/curriculumVitae/types'
import { ContactLink } from '@features/curriculumVitae/utils'
import { ToastType } from '@shared/types'
import { showToast } from '@shared/utils'

export const useContactButton = (
  resumeId: string,
  link: ContactLink,
  getToken: () => Promise<string | undefined>,
) => {
    const { viewId, track } = useAnalytics()
  const [value, setValue] = useState(link.isHidden ? undefined : link.value)
  const [isRevealing, setIsRevealing] = useState(false)

  const reveal = async () => {
    if (value || isRevealing) return

    if (link.value) {
      setValue(link.value)
      return
    }

    setIsRevealing(true)
    try {
      const token = await getToken()
      setValue(await curriculumVitaeApi.revealContact(resumeId, link.key, viewId, token))
    } catch {
      showToast({ messageKey: 'revealFailed', type: ToastType.Error })
    } finally {
      setIsRevealing(false)
    }
  }

  const copy = async () => {
    if (!value) return

    try {
      await navigator.clipboard.writeText(value)
      track({ type: AnalyticsEventType.ContactCopy, contact: link.key })
      showToast({ messageKey: 'copied', type: ToastType.Success })
    } catch {
      showToast({ messageKey: 'copyFailed', type: ToastType.Error })
    }
  }

  const trackOpen = () => track({ type: AnalyticsEventType.ContactClick, contact: link.key })

  return { value, isRevealing, reveal, copy, trackOpen }
}
