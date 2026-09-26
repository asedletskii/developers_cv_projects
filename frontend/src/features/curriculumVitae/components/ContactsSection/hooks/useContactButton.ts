import { useState } from 'react'
import { ToastType } from '@shared/types'
import { showToast } from '@shared/utils'

export const useContactButton = (value: string, isHidden: boolean) => {
  const [isRevealed, setIsRevealed] = useState(!isHidden)

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(value)
      showToast({ messageKey: 'copied', type: ToastType.Success })
    } catch {
      showToast({ messageKey: 'copyFailed', type: ToastType.Error })
    }
  }

  return { isRevealed, reveal: () => setIsRevealed(true), copy }
}
