import { getAnalyticsRange, getPreviousRange } from './getAnalyticsRange'

describe('getAnalyticsRange', () => {
  it('период включает сегодняшний день', () => {
    const range = getAnalyticsRange(7, new Date(2026, 8, 30, 12))

    expect(range.from).toBe('2026-09-24')
    expect(range.to).toBe('2026-09-30')
  })

  it('переходит через границу месяца', () => {
    expect(getAnalyticsRange(30, new Date(2026, 9, 5, 12)).from).toBe('2026-09-06')
  })

  it('предыдущий период — той же длины и вплотную перед текущим', () => {
    const now = new Date(2026, 8, 30, 12)

    expect(getPreviousRange(7, now)).toMatchObject({ from: '2026-09-17', to: '2026-09-23' })
    expect(getAnalyticsRange(7, now).from).toBe('2026-09-24')
  })
})
