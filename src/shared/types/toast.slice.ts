import { AppLocaleDictionary, NestedKeyOf } from './locale'
import { PayloadAction, SliceCaseReducers } from './reduxToolkit'

/** Ключ текста тоста из словаря `Toasts` */
export type ToastKey = NestedKeyOf<AppLocaleDictionary['Toasts']>

/** Вид тоста (совпадает с вариантами notistack) */
export enum ToastType {
  Success = 'success',
  Error = 'error',
  Warning = 'warning',
  Info = 'info',
}

export type ToastOptions = {
  id: number
  messageKey: ToastKey
  titleKey?: ToastKey
  type: ToastType
  duration?: number
}

export type ToastSliceState = {
  toasts: ToastOptions[]
}

export type ToastSliceReducers = SliceCaseReducers<ToastSliceState> & {
  addToast(state: ToastSliceState, action: PayloadAction<ToastOptions>): void

  removeToast(state: ToastSliceState, action: PayloadAction<number>): void
}
