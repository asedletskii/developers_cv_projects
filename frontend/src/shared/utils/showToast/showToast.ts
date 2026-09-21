import { DEFAULT_TOAST_DURATION } from '@shared/constants'
import { store } from '@shared/store'
import { addToast } from '@shared/store/toast'
import { ToastOptions } from '@shared/types'
import { closeToast } from '../closeToast'
import { generateId } from '../generateId'

export function showToast(toast: Omit<ToastOptions, 'id'>) {
  const id = generateId()
  const duration = toast.duration ?? DEFAULT_TOAST_DURATION[toast.type]
  store.dispatch(addToast({ id, ...toast, duration }))

  if (duration > 0) {
    setTimeout(() => {
      closeToast(id)
    }, duration)
  }
}
