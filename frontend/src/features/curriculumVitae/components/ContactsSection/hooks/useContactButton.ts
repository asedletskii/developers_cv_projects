import { useState } from 'react'
import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { ContactLink } from '@features/curriculumVitae/utils'
import { ToastType } from '@shared/types'
import { showToast } from '@shared/utils'

export const useContactButton = (
  resumeId: string,
  link: ContactLink,
  getToken: () => Promise<string | undefined>,
) => {
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
      setValue(await curriculumVitaeApi.revealContact(resumeId, link.key, token))
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
      showToast({ messageKey: 'copied', type: ToastType.Success })
    } catch {
      showToast({ messageKey: 'copyFailed', type: ToastType.Error })
    }
  }

  return { value, isRevealing, reveal, copy }
}
