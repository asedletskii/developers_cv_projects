import { ToastSliceState } from '@shared/types'

export const getToastInitialState = (): ToastSliceState => {
  return {
    toasts: [],
  }
}
