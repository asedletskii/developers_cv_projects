import { DEFAULT_TOAST_DURATION } from '@shared/constants'
import { store } from '@shared/store'
import { selectToasts } from '@shared/store/toast'
import { ToastType } from '@shared/types'
import { showToast } from './showToast'

const toasts = () => selectToasts(store.getState())

describe('showToast', () => {
  beforeEach(() => {
    store.dispatch({ type: 'store/reset' })
    jest.useFakeTimers()
  })

  afterEach(() => {
    jest.useRealTimers()
  })

  it('кладёт тост в стор со временем жизни по умолчанию для его типа', () => {
    showToast({ messageKey: 'copied', type: ToastType.Success })

    expect(toasts()).toHaveLength(1)
    expect(toasts()[0]).toMatchObject({
      messageKey: 'copied',
      type: ToastType.Success,
      duration: DEFAULT_TOAST_DURATION[ToastType.Success],
    })
  })

  it('убирает тост из стора, когда время жизни вышло', () => {
    showToast({ messageKey: 'copied', type: ToastType.Success })

    jest.advanceTimersByTime(DEFAULT_TOAST_DURATION[ToastType.Success] - 1)
    expect(toasts()).toHaveLength(1)

    jest.advanceTimersByTime(1)
    expect(toasts()).toHaveLength(0)
  })

  it('уважает своё время жизни', () => {
    showToast({ messageKey: 'copied', type: ToastType.Success, duration: 500 })

    jest.advanceTimersByTime(500)

    expect(toasts()).toHaveLength(0)
  })

  it('тост с нулевой длительностью остаётся, пока его не закроют', () => {
    showToast({ messageKey: 'copyFailed', type: ToastType.Error, duration: 0 })

    jest.advanceTimersByTime(60_000)

    expect(toasts()).toHaveLength(1)
  })

  it('разные типы живут разное время', () => {
    expect(DEFAULT_TOAST_DURATION[ToastType.Error]).toBeGreaterThan(
      DEFAULT_TOAST_DURATION[ToastType.Success],
    )
  })

  it('каждому тосту выдаётся свой id', () => {
    showToast({ messageKey: 'copied', type: ToastType.Success })
    showToast({ messageKey: 'copyFailed', type: ToastType.Error })

    const [first, second] = toasts()

    expect(first.id).not.toBe(second.id)
  })
})
