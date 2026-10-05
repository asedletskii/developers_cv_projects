import { useState } from 'react'
import { ANALYTICS_NOTICE_STORAGE_KEY } from '@features/curriculumVitae/constants'

const readDismissed = () => {
  try {
    return localStorage.getItem(ANALYTICS_NOTICE_STORAGE_KEY) === '1'
  } catch {
    return false
  }
}

export const useAnalyticsNotice = () => {
  const [isDismissed, setIsDismissed] = useState(readDismissed)

  const dismiss = () => {
    setIsDismissed(true)

    try {
      localStorage.setItem(ANALYTICS_NOTICE_STORAGE_KEY, '1')
    } catch {
      // не страшно: покажем снова при следующем заходе
    }
  }

  return { isDismissed, dismiss }
}
