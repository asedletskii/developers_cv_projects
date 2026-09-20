import { formatDuration } from './formatDuration'

describe('formatDuration', () => {
  beforeEach(() => {
    jest.useFakeTimers().setSystemTime(new Date('2026-09-21T12:00:00Z'))
  })

  afterEach(() => {
    jest.useRealTimers()
  })

  it('показывает годы и месяцы, считая оба крайних месяца', () => {
    expect(formatDuration('2023-03', '2026-09')).toBe('3 years 7 months')
  })

  it('показывает только месяцы, если года не набралось', () => {
    expect(formatDuration('2026-07', '2026-09')).toBe('3 months')
  })

  it('показывает только годы, если месяцев не осталось', () => {
    expect(formatDuration('2024-01', '2025-12')).toBe('2 years')
  })

  it('минимум один месяц, даже если начало и конец в одном месяце', () => {
    expect(formatDuration('2026-09', '2026-09')).toBe('1 month')
  })

  it('без даты окончания считает до сегодняшнего дня', () => {
    expect(formatDuration('2026-03')).toBe('7 months')
  })
})
