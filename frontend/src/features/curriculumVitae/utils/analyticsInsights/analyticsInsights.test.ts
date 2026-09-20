import { getDelta, getPeakDay, groupByWeekday, toPercent } from './analyticsInsights'

const day = (date: string, views: number, reads = 0) => ({
  date,
  views,
  reads,
  uniqueVisitors: 0,
})

describe('toPercent', () => {
  it('округляет долю', () => {
    expect(toPercent(1, 3)).toBe(33)
    expect(toPercent(2, 3)).toBe(67)
  })

  it('без базы — null', () => {
    expect(toPercent(5, 0)).toBeNull()
  })
})

describe('getDelta', () => {
  it('рост и падение в процентах', () => {
    expect(getDelta(13, 10)).toBe(30)
    expect(getDelta(5, 10)).toBe(-50)
    expect(getDelta(10, 10)).toBe(0)
  })

  it('если прошлый период пуст — null, а не бесконечность', () => {
    expect(getDelta(7, 0)).toBeNull()
  })
})

describe('groupByWeekday', () => {
  it('понедельник — первый, сумма по одноимённым дням', () => {
    // 2026-09-28 — понедельник, 2026-10-04 — воскресенье
    const result = groupByWeekday([
      day('2026-09-28', 2, 1),
      day('2026-10-05', 3, 2),
      day('2026-10-04', 4),
    ])

    expect(result[0]).toEqual({ views: 5, reads: 3 })
    expect(result[6]).toEqual({ views: 4, reads: 0 })
    expect(result[2]).toEqual({ views: 0, reads: 0 })
  })
})

describe('getPeakDay', () => {
  it('находит день с максимумом открытий', () => {
    expect(getPeakDay([day('a', 1), day('b', 5), day('c', 3)])?.date).toBe('b')
  })

  it('без открытий — undefined', () => {
    expect(getPeakDay([day('a', 0), day('b', 0)])).toBeUndefined()
    expect(getPeakDay([])).toBeUndefined()
  })
})
