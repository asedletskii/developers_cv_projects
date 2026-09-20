import { store } from '@shared/store'
import { addToast, selectToasts } from '@shared/store/toast'
import { ToastType } from '@shared/types'
import { closeToast } from './closeToast'

const add = (id: number) =>
  store.dispatch(addToast({ id, messageKey: 'copied', type: ToastType.Success }))

const ids = () => selectToasts(store.getState()).map(({ id }) => id)

describe('closeToast', () => {
  beforeEach(() => {
    store.dispatch({ type: 'store/reset' })
  })

  it('убирает тост из стора по id', () => {
    add(1)
    add(2)

    closeToast(1)

    expect(ids()).toEqual([2])
  })

  it('неизвестный id ничего не меняет', () => {
    add(1)

    closeToast(999)

    expect(ids()).toEqual([1])
  })
})
