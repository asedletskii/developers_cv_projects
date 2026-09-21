import { ToastType } from '@shared/types'

export const DEFAULT_TOAST_DURATION: Record<ToastType, number> = {
  [ToastType.Success]: 3000,
  [ToastType.Info]: 5000,
  [ToastType.Warning]: 5000,
  [ToastType.Error]: 7000,
}
