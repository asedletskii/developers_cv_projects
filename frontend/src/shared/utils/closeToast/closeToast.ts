import { store } from '@shared/store'
import { removeToast } from '@shared/store/toast'

export function closeToast(id: number) {
  store.dispatch(removeToast(id))
}
