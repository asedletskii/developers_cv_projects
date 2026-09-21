import { formatMonth } from './formatMonth'

describe('formatMonth', () => {
  it('форматирует месяц и год', () => {
    expect(formatMonth('2023-03')).toBe('Mar 2023')
  })

  it('принимает полную дату', () => {
    expect(formatMonth('1999-12-31')).toBe('Dec 1999')
  })
})
